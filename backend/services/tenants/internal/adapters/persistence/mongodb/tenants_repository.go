package mongodb

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoDBTenantsRepository struct {
	base *database.MongoDBRepository[tenantDocument]
}

func newMongoDBTenantsRepository(db *mongo.Database, logger *log.Logger) ports.TenantsRepository {
	repository := database.NewMongoDBRepository[tenantDocument](
		db.Collection("tenants"),
		logger,
	)

	return &mongoDBTenantsRepository{
		base: repository,
	}
}

func (r *mongoDBTenantsRepository) Find(ctx context.Context) ([]*domain.Tenant, error) {
	documents, err := r.base.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	tenants := make([]*domain.Tenant, 0, len(documents))

	for _, document := range documents {
		tenants = append(tenants, fromTenantDocument(&document))
	}

	return tenants, nil
}

func (r *mongoDBTenantsRepository) FindByID(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
	result, err := r.base.FindOne(ctx, bson.M{"_id": id})

	if err != nil {
		return nil, err
	}

	return fromTenantDocument(&result), nil
}

func (r *mongoDBTenantsRepository) Exists(ctx context.Context, id domain.TenantID) (bool, error) {
	return r.base.Exists(ctx, bson.M{"_id": id})
}

func (r *mongoDBTenantsRepository) Upsert(ctx context.Context, t *domain.Tenant) (*domain.Tenant, error) {
	doc := toTenantDocument(t)
	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result tenantDocument
	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.base.Collection().FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		return nil, err
	}

	return fromTenantDocument(&result), nil
}

func (r *mongoDBTenantsRepository) Remove(ctx context.Context, id domain.TenantID) error {
	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		result, err := r.base.Collection().DeleteOne(sctx, bson.M{"_id": id})

		if err != nil {
			return err
		}

		if result.DeletedCount == 0 {
			return common.ErrNotFound
		}

		return nil
	})

	return err
}
