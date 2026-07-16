package mongodb

import (
	"context"
	"log/slog"

	"sundance/backend/pkg/database"
	"sundance/backend/services/forms/internal/adapters/persistence/mongodb/documents"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	formVersionIndexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "form_id", Value: 1},
				{Key: "version", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
)

type mongoDBFormVersionsRepository struct {
	base   *database.MongoDBRepository[documents.FormVersionDocument]
	outbox *mongoDBOutboxRepository
}

func newMongoDBFormVersionsRepository(db *mongo.Database, outbox *mongoDBOutboxRepository, logger *slog.Logger) (ports.FormVersionRepository, error) {
	base := database.NewMongoDBRepository[documents.FormVersionDocument](
		db.Collection("form_versions"),
		logger,
	)
	repository := &mongoDBFormVersionsRepository{base, outbox}

	if err := repository.migrate(context.Background()); err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *mongoDBFormVersionsRepository) migrate(ctx context.Context) error {
	_, err := r.base.Collection().Indexes().CreateMany(ctx, formVersionIndexes)
	return err
}

func (r *mongoDBFormVersionsRepository) Find(ctx context.Context, formID domain.FormID) ([]*domain.FormVersion, error) {
	docs, err := r.base.Find(ctx, bson.M{"form_id": formID})

	if err != nil {
		return nil, err
	}

	versions := make([]*domain.FormVersion, 0, len(docs))
	for _, document := range docs {
		v, err := documents.FromFormVersionDocument(&document)

		if err != nil {
			return nil, err
		}

		versions = append(versions, v)
	}

	return versions, nil
}

func (r *mongoDBFormVersionsRepository) FindByID(ctx context.Context, versionID domain.FormVersionID) (*domain.FormVersion, error) {
	document, err := r.base.FindOne(ctx, bson.M{"_id": versionID})

	if err != nil {
		return nil, err
	}

	return documents.FromFormVersionDocument(&document)
}

func (r *mongoDBFormVersionsRepository) FindNextVersionNumber(ctx context.Context, formID domain.FormID) (int, error) {
	r.base.Logger().DebugContext(ctx, "finding next version number", "form_id", formID)

	filter := bson.M{"form_id": formID}
	opts := options.Find().SetSort(bson.M{"version": -1}).SetLimit(1).SetProjection(bson.M{"version": 1})

	var lv int
	err := r.base.WithSession(ctx, func(sctx context.Context) error {
		cursor, err := r.base.Collection().Find(sctx, filter, opts)

		if err != nil {
			return err
		}

		defer cursor.Close(sctx)

		var documents []documents.FormVersionDocument
		if err := cursor.All(sctx, &documents); err != nil {
			return err
		}

		if len(documents) == 0 {
			lv = 1
			return nil
		}

		lv = documents[0].Version + 1
		return nil
	})

	if err != nil {
		r.base.Logger().ErrorContext(ctx, "failed to find next version number", "form_id", formID, "error", err)
		return 0, err
	}

	return lv, nil
}

func (r *mongoDBFormVersionsRepository) Upsert(ctx context.Context, v *domain.FormVersion) (*domain.FormVersion, error) {
	r.base.Logger().DebugContext(ctx, "upsert version", "version_id", v.ID, "form_id", v.FormID)

	doc, err := documents.ToFormVersionDocument(v)

	if err != nil {
		r.base.Logger().ErrorContext(ctx, "failed to convert version to document", "version_id", v.ID, "form_id", v.FormID, "error", err)
		return nil, err
	}

	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result documents.FormVersionDocument
	err = r.base.WithSession(ctx, func(sctx context.Context) error {
		if err := r.base.Collection().FindOneAndUpdate(sctx, filter, update, opts).Decode(&result); err != nil {
			return err
		}

		return r.WriteEvents(sctx, v)
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, domain.ErrDuplicateVersion
		}

		r.base.Logger().ErrorContext(ctx, "mongo upsert failed", "version_id", v.ID, "form_id", v.FormID, "error", err)
		return nil, err
	}

	v.DrainEvents()

	return documents.FromFormVersionDocument(&result)
}

func (r *mongoDBFormVersionsRepository) WriteEvents(ctx context.Context, e domain.HasEvents) error {
	for event := range e.PeekEvents() {
		if _, err := r.outbox.Upsert(ctx, &event); err != nil {
			return err
		}
	}

	return nil
}
