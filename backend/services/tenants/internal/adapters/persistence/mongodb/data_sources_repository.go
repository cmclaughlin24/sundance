package mongodb

import (
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoDBDataSourcesRepository struct {
	mongodbBaseRepository
}

func NewMongoDBDataSourcesRepository(db *mongo.Database, logger *log.Logger) *MongoDBDataSourcesRepository {
	return &MongoDBDataSourcesRepository{
		mongodbBaseRepository{
			collection: db.Collection("data_sources"),
			logger:     logger,
		},
	}
}
