package mongodb

import (
	"log/slog"

	"sundance/backend/pkg/database"
	"sundance/backend/services/forms/internal/core/ports"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Bootstrap(client *mongo.Client, logger *slog.Logger, databaseName string) (*ports.Repository, error) {
	db := client.Database(databaseName)

	tags, err := newMongoDBCanonicalTagsRepository(db, logger)
	if err != nil {
		return nil, err
	}

	forms, err := newMongoDBFormsRepository(db, logger)
	if err != nil {
		return nil, err
	}

	versions, err := newMongoDBFormVersionsRepository(db, logger)
	if err != nil {
		return nil, err
	}

	submissions, err := newMongoDBSubmissionsRepository(db, logger)
	if err != nil {
		return nil, err
	}

	return &ports.Repository{
		Database:      database.NewMongoDBDatabase(client, db, logger),
		Tags: tags,
		Forms:         forms,
		FormVersions:  versions,
		Submissions:   submissions,
	}, nil
}
