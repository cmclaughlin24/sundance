package mongodb

import (
	"log/slog"

	"sundance/backend/pkg/database"
	"sundance/backend/services/forms/internal/core/ports"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Bootstrap(client *mongo.Client, logger *slog.Logger, databaseName string) (*ports.Repository, error) {
	db := client.Database(databaseName)

	forms, err := newMongoDBFormsRepository(db, logger)
	if err != nil {
		return nil, err
	}

	formVersions, err := newMongoDBFormVersionsRepository(db, logger)
	if err != nil {
		return nil, err
	}

	submissions, err := newMongoDBSubmissionsRepository(db, logger)
	if err != nil {
		return nil, err
	}

	tags, err := newMongoDBTagsRepository(db, logger)
	if err != nil {
		return nil, err
	}

	tagVersions, err := newMongoDBTagVersionsRepository(db, logger)
	if err != nil {
		return nil, err
	}

	return &ports.Repository{
		Database:     database.NewMongoDBDatabase(client, db, logger),
		Forms:        forms,
		FormVersions: formVersions,
		Submissions:  submissions,
		Tags:         tags,
		TagVersions:  tagVersions,
	}, nil
}
