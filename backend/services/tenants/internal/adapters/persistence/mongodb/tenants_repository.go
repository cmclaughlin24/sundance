package mongodb

import (
	"log"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
