package mongodb

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/database"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoDBDatabase struct {
	client *mongo.Client
}

func NewMongoDBDatabase(client *mongo.Client, _ *mongo.Database) database.Database {
	return &MongoDBDatabase{
		client: client,
	}
}

func (db *MongoDBDatabase) Close() error {
	return db.client.Disconnect(context.Background())
}

func (db *MongoDBDatabase) BeginTx(ctx context.Context) (context.Context, error) {
	session, err := db.client.StartSession()

	if err != nil {
		return ctx, err
	}

	if err := session.StartTransaction(); err != nil {
		return ctx, err
	}

	return mongo.NewSessionContext(ctx, session), nil
}

func (db *MongoDBDatabase) CommitTx(ctx context.Context) error {
	session := mongo.SessionFromContext(ctx)

	if session == nil {
		return nil
	}

	defer session.EndSession(ctx)
	return session.CommitTransaction(ctx)
}

func (db *MongoDBDatabase) RollbackTx(ctx context.Context) error {
	session := mongo.SessionFromContext(ctx)

	if session == nil {
		return nil
	}

	defer session.EndSession(ctx)
	return session.AbortTransaction(ctx)
}
