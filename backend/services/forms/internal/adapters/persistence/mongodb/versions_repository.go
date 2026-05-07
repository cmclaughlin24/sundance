package mongodb

import (
	"context"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "form_id", Value: 1},
				{Key: "version", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
)

type mongoDBVersionsRepository struct {
	base *database.MongoDBRepository[versionDocument]
}

func newMongoDBVersionsRepository(db *mongo.Database, logger *slog.Logger) (ports.VersionRepository, error) {
	base := database.NewMongoDBRepository[versionDocument](
		db.Collection("versions"),
		logger,
	)
	repository := &mongoDBVersionsRepository{base}

	if err := repository.migrate(context.Background()); err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *mongoDBVersionsRepository) migrate(ctx context.Context) error {
	_, err := r.base.Collection().Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *mongoDBVersionsRepository) Find(ctx context.Context, formID domain.FormID) ([]*domain.Version, error) {
	documents, err := r.base.Find(ctx, bson.M{"form_id": formID})

	if err != nil {
		return nil, err
	}

	versions := make([]*domain.Version, 0, len(documents))
	for _, document := range documents {
		v, err := fromVersionDocument(&document)

		if err != nil {
			return nil, err
		}

		versions = append(versions, v)
	}

	return versions, nil
}

func (r *mongoDBVersionsRepository) FindByID(ctx context.Context, formID domain.FormID, versionID domain.VersionID) (*domain.Version, error) {
	document, err := r.base.FindOne(ctx, bson.M{"form_id": formID, "_id": versionID})

	if err != nil {
		return nil, err
	}

	return fromVersionDocument(&document)
}

func (r *mongoDBVersionsRepository) FindNextVersionNumber(ctx context.Context, formID domain.FormID) (int, error) {
	r.base.Logger().DebugContext(ctx, "finding next version number", "form_id", formID)

	filter := bson.M{"form_id": formID}
	opts := options.Find().SetSort(bson.M{"version": -1}).SetLimit(1).SetProjection(bson.M{"version": 1})

	var lv int
	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		cursor, err := r.base.Collection().Find(sctx, filter, opts)

		if err != nil {
			return err
		}

		defer cursor.Close(sctx)

		var documents []versionDocument
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

func (r *mongoDBVersionsRepository) Upsert(ctx context.Context, v *domain.Version) (*domain.Version, error) {
	r.base.Logger().DebugContext(ctx, "upsert version", "version_id", v.ID, "form_id", v.FormID)

	doc, err := toVersionDocument(v)

	if err != nil {
		r.base.Logger().ErrorContext(ctx, "failed to convert version to document", "version_id", v.ID, "form_id", v.FormID, "error", err)
		return nil, err
	}

	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result versionDocument
	err = mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.base.Collection().FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, domain.ErrDuplicateVersion
		}

		r.base.Logger().ErrorContext(ctx, "mongo upsert failed", "version_id", v.ID, "form_id", v.FormID, "error", err)
		return nil, err
	}

	return fromVersionDocument(&result)
}
