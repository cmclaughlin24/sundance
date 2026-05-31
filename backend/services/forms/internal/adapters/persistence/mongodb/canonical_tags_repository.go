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
	canonicalTagIndexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "tenant_id", Value: 1},
				{Key: "key", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
)

type mongoDBCanonicalTagsRespository struct {
	base *database.MongoDBRepository[documents.CanonicalTagDocument]
}

func newMongoDBCanonicalTagsRepository(db *mongo.Database, logger *slog.Logger) (ports.CanonicalTagRepository, error) {
	base := database.NewMongoDBRepository[documents.CanonicalTagDocument](
		db.Collection("canonical_tags"),
		logger,
	)
	repository := &mongoDBCanonicalTagsRespository{base}

	if err := repository.migrate(context.Background()); err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *mongoDBCanonicalTagsRespository) migrate(ctx context.Context) error {
	_, err := r.base.Collection().Indexes().CreateMany(ctx, canonicalTagIndexes)
	return err
}

func (r *mongoDBCanonicalTagsRespository) Find(ctx context.Context, ctf ports.CanonicalTagFilters) ([]*domain.CanonicalTag, error) {
	filter := bson.M{}

	if ctf.TenantID != "" {
		filter["tenant_id"] = ctf.TenantID
	}

	docs, err := r.base.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	tags := make([]*domain.CanonicalTag, 0, len(docs))
	for _, document := range docs {
		tags = append(tags, documents.FromCanonicalTagDocument(document))
	}

	return tags, nil
}

func (r *mongoDBCanonicalTagsRespository) FindByID(ctx context.Context, id domain.CanonicalTagID) (*domain.CanonicalTag, error) {
	doc, err := r.base.FindOne(ctx, bson.M{"_id": id})

	if err != nil {
		return nil, err
	}

	return documents.FromCanonicalTagDocument(doc), nil
}

func (r *mongoDBCanonicalTagsRespository) Upsert(ctx context.Context, t *domain.CanonicalTag) (*domain.CanonicalTag, error) {
	r.base.Logger().DebugContext(ctx, "upsert canonical tag", "canonical_tag_id", t.ID)

	doc := documents.ToCanonicalTagDocument(t)
	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result documents.CanonicalTagDocument
	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.base.Collection().FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		r.base.Logger().ErrorContext(ctx, "mongo upsert failed", "canonical_tag_id", t.ID, "error", err)
		return nil, err
	}

	return documents.FromCanonicalTagDocument(result), nil
}

func (r *mongoDBCanonicalTagsRespository) Delete(ctx context.Context, id domain.CanonicalTagID) error {
	return r.base.Delete(ctx, bson.M{"_id": id})
}
