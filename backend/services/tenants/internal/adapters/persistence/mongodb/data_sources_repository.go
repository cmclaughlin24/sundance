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

type MongoDBDataSourcesRepository struct {
	mongodbBaseRepository
}

func NewMongoDBDataSourcesRepository(db *mongo.Database, logger *log.Logger) ports.DataSourcesRepository {
	return &MongoDBDataSourcesRepository{
		mongodbBaseRepository{
			collection: db.Collection("data_sources"),
			logger:     logger,
		},
	}
}

func (r *MongoDBDataSourcesRepository) Find(ctx context.Context, tenantID domain.TenantID) ([]*domain.DataSource, error) {
	var documents []dataSourceDocument

	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		cursor, err := r.collection.Find(sctx, bson.M{"tenant_id": tenantID})

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

	dataSources := make([]*domain.DataSource, 0, len(documents))

	for _, document := range documents {
		ds, err := fromDataSourceDocument(&document)

		if err != nil {
			return nil, err
		}

		dataSources = append(dataSources, ds)
	}

	return dataSources, nil
}

func (r *MongoDBDataSourcesRepository) FindById(ctx context.Context, tenantID domain.TenantID, sourceID domain.DataSourceID) (*domain.DataSource, error) {
	var result dataSourceDocument

	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.collection.FindOne(sctx, bson.M{"_id": sourceID, "tenant_id": tenantID}).Decode(&result)
	})

	if err != nil {
		return nil, err
	}

	return fromDataSourceDocument(&result)
}

func (r *MongoDBDataSourcesRepository) Exists(ctx context.Context, tenantID domain.TenantID, sourceID domain.DataSourceID) (bool, error) {
	filter := bson.M{"_id": sourceID, "tenant_id": tenantID}
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

func (r *MongoDBDataSourcesRepository) Upsert(ctx context.Context, ds *domain.DataSource) (*domain.DataSource, error) {
	now := time.Now()

	// TODO: Move to the domain layer of the application.
	if ds.ID == "" {
		ds.ID = domain.DataSourceID(uuid.New().String())
		ds.CreatedAt = now
	}
	ds.UpdatedAt = now

	doc, err := toDataSourceDocument(ds)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result dataSourceDocument
	err = mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.collection.FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		return nil, err
	}

	return fromDataSourceDocument(&result)
}

func (r *MongoDBDataSourcesRepository) Remove(ctx context.Context, tenantID domain.TenantID, sourceID domain.DataSourceID) error {
	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		result, err := r.collection.DeleteOne(sctx, bson.M{"_id": sourceID, "tenant_id": tenantID})

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
