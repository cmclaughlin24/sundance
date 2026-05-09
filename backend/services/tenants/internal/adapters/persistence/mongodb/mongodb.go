package mongodb

import (
	"log/slog"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Package declaration for the current time function. Allows for easier testing by enabling the injection of a
// mock time function.
var Now = time.Now

func Bootstrap(client *mongo.Client, logger *slog.Logger) *ports.Repository {
	db := client.Database("tenants")

	return &ports.Repository{
		Database:    database.NewMongoDBDatabase(client, db, logger),
		DataSources: newMongoDBDataSourcesRepository(db, logger),
		Tenants:     newMongoDBTenantsRepository(db, logger),
	}
}
