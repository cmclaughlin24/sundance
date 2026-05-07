package mongodb

import (
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Bootstrap(client *mongo.Client, logger *slog.Logger) (*ports.Repository, error) {
	db := client.Database("forms")

	forms, err := newMongoDBFormsRepository(db, logger)
	if err != nil {
		return nil, err
	}

	versions, err := newMongoDBVersionsRepository(db, logger)
	if err != nil {
		return nil, err
	}

	return &ports.Repository{
		Database: database.NewMongoDBDatabase(client, db, logger),
		Forms:    forms,
		Versions: versions,
	}, nil
}
