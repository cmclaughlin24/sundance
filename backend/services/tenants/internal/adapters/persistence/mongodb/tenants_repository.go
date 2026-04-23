package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBTenantsRepository struct {
	mongodbBaseRepository
}

func NewMongoDBTenantsRepository(db *mongo.Database, logger *log.Logger) ports.TenantsRepository {
	return &MongoDBTenantsRepository{
		mongodbBaseRepository{
			collection: db.Collection("tenants"),
			logger:     logger,
		},
	}
}

func (r *MongoDBTenantsRepository) Find(ctx context.Context) ([]*domain.Tenant, error) {
	var documents []tenantDocument

	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		cursor, err := r.collection.Find(sctx, bson.M{})

		if err != nil {
			return err
		}

		if err = cursor.All(sctx, &documents); err != nil {
			return fmt.Errorf("an error occurred reading the documents: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	tenants := make([]*domain.Tenant, 0, len(documents))

	for _, document := range documents {
		tenants = append(tenants, fromTenantDocument(&document))
	}

	return tenants, nil
}

func (r *MongoDBTenantsRepository) FindById(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
	var result tenantDocument

	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.collection.FindOne(sctx, bson.M{"_id": id}).Decode(&result)
	})

	if err != nil {
		return nil, err
	}

	return fromTenantDocument(&result), nil
}

func (r *MongoDBTenantsRepository) Exists(ctx context.Context, id domain.TenantID) (bool, error) {
	filter := bson.M{"_id": id}
	opts := options.Count().SetLimit(1)
	var count int64

	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		c, err := r.collection.CountDocuments(sctx, filter, opts)

		if err != nil {
			return err
		}

		count = c
		return nil
	})

	return count > 0, err
}

func (r *MongoDBTenantsRepository) Upsert(ctx context.Context, t *domain.Tenant) (*domain.Tenant, error) {
	now := time.Now()

	// TODO: Move to the domain layer of the application.
	if t.ID == "" {
		t.ID = domain.TenantID(uuid.New().String())
		t.CreatedAt = now
	}
	t.UpdatedAt = now

	doc := toTenantDocument(t)
	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result tenantDocument
	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.collection.FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		return nil, err
	}

	return fromTenantDocument(&result), nil
}

func (r *MongoDBTenantsRepository) Remove(ctx context.Context, id domain.TenantID) error {
	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		result, err := r.collection.DeleteOne(sctx, bson.M{"_id": id})

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
