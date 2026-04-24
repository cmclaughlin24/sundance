package mongodb

import (
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Bootstrap(client *mongo.Client, logger *log.Logger) *ports.Repository {
	db := client.Database("tenants")

	return &ports.Repository{
		Database:    database.NewMongoDBDatabase(client, db),
		DataSources: newMongoDBDataSourcesRepository(db, logger),
		Tenants:     newMongoDBTenantsRepository(db, logger),
	}
}
