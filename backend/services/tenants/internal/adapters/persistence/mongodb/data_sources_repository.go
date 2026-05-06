package mongodb

import (
	"context"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoDBDataSourcesRepository struct {
	base *database.MongoDBRepository[dataSourceDocument]
}

func newMongoDBDataSourcesRepository(db *mongo.Database, logger *slog.Logger) ports.DataSourcesRepository {
	base := database.NewMongoDBRepository[dataSourceDocument](
		db.Collection("data_sources"),
		logger,
	)

	return &mongoDBDataSourcesRepository{base}
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
func (r *mongoDBDataSourcesRepository) FindByID(ctx context.Context, tenantID domain.TenantID, sourceID domain.DataSourceID) (*domain.DataSource, error) {
	result, err := r.base.FindOne(ctx, bson.M{"_id": sourceID, "tenant_id": tenantID})

	if err != nil {
		return nil, err
	}

	return fromDataSourceDocument(&result)
}

func (r *mongoDBDataSourcesRepository) Exists(ctx context.Context, tenantID domain.TenantID, sourceID domain.DataSourceID) (bool, error) {
	return r.base.Exists(ctx, bson.M{"_id": sourceID, "tenant_id": tenantID})
}

func (r *mongoDBDataSourcesRepository) Upsert(ctx context.Context, ds *domain.DataSource) (*domain.DataSource, error) {
	r.base.Logger().DebugContext(ctx, "upsert data source", "tenant_id", ds.TenantID, "data_source_id", ds.ID)

	doc, err := toDataSourceDocument(ds)
	if err != nil {
		r.base.Logger().ErrorContext(ctx, "failed to convert data source to document", "tenant_id", ds.TenantID, "data_source_id", ds.ID, "error", err)
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
		r.base.Logger().ErrorContext(ctx, "mongo upsert failed", "tenant_id", ds.TenantID, "data_source_id", ds.ID, "error", err)
		return nil, err
	}

	return fromDataSourceDocument(&result)
}

func (r *mongoDBDataSourcesRepository) Delete(ctx context.Context, tenantID domain.TenantID, sourceID domain.DataSourceID) error {
	return r.base.Delete(ctx, bson.M{"_id": sourceID, "tenant_id": tenantID})
}

func (r *mongoDBDataSourcesRepository) DeleteAll(ctx context.Context, tenantID domain.TenantID) error {
	r.base.Logger().DebugContext(ctx, "deleting all data sources for tenant", "tenant_id", tenantID)

	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		_, err := r.base.Collection().DeleteMany(sctx, bson.M{"tenant_id": tenantID})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		r.base.Logger().ErrorContext(ctx, "failed to delete all data sources", "tenant_id", tenantID, "error", err)
	}

	return err
}
