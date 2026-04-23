package mongodb

import (
	"log"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
