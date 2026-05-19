package mongodb

import (
	"log/slog"

	"sundance/backend/pkg/database"
	"sundance/backend/services/tenants/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Bootstrap(client *mongo.Client, logger *slog.Logger, databaseName string) *ports.Repository {
	db := client.Database(databaseName)

	return &ports.Repository{
		Database:    database.NewMongoDBDatabase(client, db, logger),
		DataSources: newMongoDBDataSourcesRepository(db, logger),
		Tenants:     newMongoDBTenantsRepository(db, logger),
	}
}
