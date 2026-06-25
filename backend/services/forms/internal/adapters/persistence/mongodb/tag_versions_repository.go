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
	tagVersionIndexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "tag_id", Value: 1},
				{Key: "version", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
)

type mongoDBTagVersionsRepository struct {
	base *database.MongoDBRepository[documents.TagVersionDocument]
}

func newMongoDBTagVersionsRepository(db *mongo.Database, logger *slog.Logger) (ports.TagVersionsRepository, error) {
	base := database.NewMongoDBRepository[documents.TagVersionDocument](
		db.Collection("tag_versions"),
		logger,
	)
	repository := &mongoDBTagVersionsRepository{base}

	if err := repository.migrate(context.Background()); err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *mongoDBTagVersionsRepository) migrate(ctx context.Context) error {
	_, err := r.base.Collection().Indexes().CreateMany(ctx, tagVersionIndexes)
	return err
}

func (r *mongoDBTagVersionsRepository) Find(ctx context.Context, filters ports.TagVersionFilters) ([]*domain.TagVersion, error) {
	f := bson.M{}

	if filters.TagID != "" {
		f["tag_id"] = filters.TagID
	}

	if len(filters.Statuses) != 0 {
		f["status"] = bson.M{"$in": filters.Statuses}
	}

	docs, err := r.base.Find(ctx, f)

	if err != nil {
		return nil, err
	}

	versions := make([]*domain.TagVersion, 0, len(docs))
	for _, document := range docs {
		versions = append(versions, documents.FromTagVersionDocument(document))
	}

	return versions, nil
}

func (r *mongoDBTagVersionsRepository) FindByIDs(ctx context.Context, ids []domain.TagVersionID) ([]*domain.TagVersion, error) {
	docs, err := r.base.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	versions := make([]*domain.TagVersion, 0, len(docs))
	for _, document := range docs {
		versions = append(versions, documents.FromTagVersionDocument(document))
	}

	return versions, nil
}

func (r *mongoDBTagVersionsRepository) FindByID(ctx context.Context, id domain.TagVersionID) (*domain.TagVersion, error) {
	doc, err := r.base.FindOne(ctx, bson.M{"_id": id})

	if err != nil {
		return nil, err
	}

	return documents.FromTagVersionDocument(doc), nil
}

func (r *mongoDBTagVersionsRepository) FindNextVersionNumber(ctx context.Context, tagID domain.TagID) (int, error) {
	r.base.Logger().DebugContext(ctx, "finding next tag version number", "tag_id", tagID)

	filter := bson.M{"tag_id": tagID}
	opts := options.Find().SetSort(bson.M{"version": -1}).SetLimit(1).SetProjection(bson.M{"version": 1})

	var lv int
	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		cursor, err := r.base.Collection().Find(sctx, filter, opts)

		if err != nil {
			return err
		}

		defer cursor.Close(sctx)

		var docs []documents.TagVersionDocument
		if err := cursor.All(sctx, &docs); err != nil {
			return err
		}

		if len(docs) == 0 {
			lv = 1
			return nil
		}

		lv = docs[0].Version + 1
		return nil
	})

	if err != nil {
		r.base.Logger().ErrorContext(ctx, "failed to find next tag version number", "tag_id", tagID, "error", err)
		return 0, err
	}

	return lv, nil
}

func (r *mongoDBTagVersionsRepository) Upsert(ctx context.Context, v *domain.TagVersion) (*domain.TagVersion, error) {
	r.base.Logger().DebugContext(ctx, "upsert tag version", "tag_version_id", v.ID, "tag_id", v.TagID)

	doc := documents.ToTagVersionDocument(v)
	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result documents.TagVersionDocument
	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.base.Collection().FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, domain.ErrDuplicateTagVersion
		}

		r.base.Logger().ErrorContext(ctx, "mongo upsert failed", "tag_version_id", v.ID, "tag_id", v.TagID, "error", err)
		return nil, err
	}

	return documents.FromTagVersionDocument(result), nil
}
