package mongodb

import (
	"log/slog"
	"time"

	"sundance/backend/pkg/database"
	"sundance/backend/services/forms/internal/core/ports"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var Now = time.Now

func Bootstrap(client *mongo.Client, logger *slog.Logger, databaseName string) (*ports.Repository, error) {
	db := client.Database(databaseName)

	outbox, err := newMongoDBOutboxRepository(db, logger)
	if err != nil {
		return nil, err
	}

	forms, err := newMongoDBFormsRepository(db, logger)
	if err != nil {
		return nil, err
	}

	formVersions, err := newMongoDBFormVersionsRepository(db, outbox.(*mongoDBOutboxRepository), logger)
	if err != nil {
		return nil, err
	}

	submissions, err := newMongoDBSubmissionsRepository(db, outbox.(*mongoDBOutboxRepository), logger)
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
		Outbox:       outbox,
		Forms:        forms,
		FormVersions: formVersions,
		Submissions:  submissions,
		Tags:         tags,
		TagVersions:  tagVersions,
	}, nil
}
