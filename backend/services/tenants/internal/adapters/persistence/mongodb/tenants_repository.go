package mongodb

import (
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoDBTenantsRepository struct {
	mongodbBaseRepository
}

func NewMongoDBTenantsRepository(db *mongo.Database, logger *log.Logger) *MongoDBTenantsRepository {
	return &MongoDBTenantsRepository{
		mongodbBaseRepository{
			collection: db.Collection("tenants"),
			logger:     logger,
		},
	}
}
