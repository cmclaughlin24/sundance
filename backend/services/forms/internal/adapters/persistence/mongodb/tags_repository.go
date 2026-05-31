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
	tagIndexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "tenant_id", Value: 1},
				{Key: "key", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
)

type mongoDBTagsRespository struct {
	base *database.MongoDBRepository[documents.TagDocument]
}

func newMongoDBTagsRepository(db *mongo.Database, logger *slog.Logger) (ports.TagsRepository, error) {
	base := database.NewMongoDBRepository[documents.TagDocument](
		db.Collection("tags"),
		logger,
	)
	repository := &mongoDBTagsRespository{base}

	if err := repository.migrate(context.Background()); err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *mongoDBTagsRespository) migrate(ctx context.Context) error {
	_, err := r.base.Collection().Indexes().CreateMany(ctx, tagIndexes)
	return err
}

func (r *mongoDBTagsRespository) Find(ctx context.Context, ctf ports.TagFilters) ([]*domain.Tag, error) {
	filter := bson.M{}

	if ctf.TenantID != "" {
		filter["tenant_id"] = ctf.TenantID
	}

	docs, err := r.base.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	tags := make([]*domain.Tag, 0, len(docs))
	for _, document := range docs {
		tags = append(tags, documents.FromTagDocument(document))
	}

	return tags, nil
}

func (r *mongoDBTagsRespository) FindByID(ctx context.Context, id domain.TagID) (*domain.Tag, error) {
	doc, err := r.base.FindOne(ctx, bson.M{"_id": id})

	if err != nil {
		return nil, err
	}

	return documents.FromTagDocument(doc), nil
}

func (r *mongoDBTagsRespository) Upsert(ctx context.Context, t *domain.Tag) (*domain.Tag, error) {
	r.base.Logger().DebugContext(ctx, "upsert canonical tag", "canonical_tag_id", t.ID)

	doc := documents.ToTagDocument(t)
	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result documents.TagDocument
	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.base.Collection().FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		r.base.Logger().ErrorContext(ctx, "mongo upsert failed", "canonical_tag_id", t.ID, "error", err)
		return nil, err
	}

	return documents.FromTagDocument(result), nil
}

func (r *mongoDBTagsRespository) Delete(ctx context.Context, id domain.TagID) error {
	return r.base.Delete(ctx, bson.M{"_id": id})
}
