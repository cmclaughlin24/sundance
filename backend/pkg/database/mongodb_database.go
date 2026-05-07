package database

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoDBDatabase struct {
	client *mongo.Client
	logger *slog.Logger
}

func NewMongoDBDatabase(client *mongo.Client, _ *mongo.Database, logger *slog.Logger) Database {
	return &MongoDBDatabase{
		client: client,
		logger: logger,
	}
}

func (db *MongoDBDatabase) Close() error {
	return db.client.Disconnect(context.Background())
}

func (db *MongoDBDatabase) BeginTx(ctx context.Context) (context.Context, error) {
	db.logger.DebugContext(ctx, "starting transaction")

	session, err := db.client.StartSession()

	if err != nil {
		db.logger.ErrorContext(ctx, "failed to start session", "error", err)
		return ctx, err
	}

	if err := session.StartTransaction(); err != nil {
		db.logger.ErrorContext(ctx, "failed to start transaction", "error", err)
		return ctx, err
	}

	return mongo.NewSessionContext(ctx, session), nil
}

func (db *MongoDBDatabase) CommitTx(ctx context.Context) error {
	session := mongo.SessionFromContext(ctx)

	if session == nil {
		return nil
	}

	db.logger.DebugContext(ctx, "committing transaction")

	defer session.EndSession(ctx)
	if err := session.CommitTransaction(ctx); err != nil {
		db.logger.ErrorContext(ctx, "failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (db *MongoDBDatabase) RollbackTx(ctx context.Context) error {
	session := mongo.SessionFromContext(ctx)

	if session == nil {
		return nil
	}

	db.logger.DebugContext(ctx, "rolling back transaction")

	defer session.EndSession(ctx)
	if err := session.AbortTransaction(ctx); err != nil {
		db.logger.ErrorContext(ctx, "failed to rollback transaction", "error", err)
		return err
	}

	return nil
}
