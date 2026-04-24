package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoDBDataSourcesRepository struct {
	base *database.MongoDBRepository[dataSourceDocument]
}

func newMongoDBDataSourcesRepository(db *mongo.Database, logger *log.Logger) ports.DataSourcesRepository {
	repository := database.NewMongoDBRepository[dataSourceDocument](
		db.Collection("data_sources"),
		logger,
	)

	return &mongoDBDataSourcesRepository{
		base: repository,
	}
}

func (r *mongoDBDataSourcesRepository) Find(ctx context.Context, tenantID domain.TenantID) ([]*domain.DataSource, error) {
	documents, err := r.base.Find(ctx, bson.M{"tenant_id": tenantID})

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
func (r *mongoDBDataSourcesRepository) FindById(ctx context.Context, tenantID domain.TenantID, sourceID domain.DataSourceID) (*domain.DataSource, error) {
	result, err := r.base.FindById(ctx, bson.M{"_id": sourceID, "tenant_id": tenantID})

	if err != nil {
		return nil, err
	}

	return fromDataSourceDocument(&result)
}

func (r *mongoDBDataSourcesRepository) Exists(ctx context.Context, tenantID domain.TenantID, sourceID domain.DataSourceID) (bool, error) {
	return r.base.Exists(ctx, bson.M{"_id": sourceID, "tenant_id": tenantID})
}

func (r *mongoDBDataSourcesRepository) Upsert(ctx context.Context, ds *domain.DataSource) (*domain.DataSource, error) {
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
		return r.base.Collection().FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		return nil, err
	}

	return fromDataSourceDocument(&result)
}

func (r *mongoDBDataSourcesRepository) Remove(ctx context.Context, tenantID domain.TenantID, sourceID domain.DataSourceID) error {
	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		result, err := r.base.Collection().DeleteOne(sctx, bson.M{"_id": sourceID, "tenant_id": tenantID})

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
