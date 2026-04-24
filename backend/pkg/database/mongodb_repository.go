package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBRepository[T any] struct {
	collection *mongo.Collection
	logger     *log.Logger
}

func NewMongoDBRepository[T any](collection *mongo.Collection, logger *log.Logger) *MongoDBRepository[T] {
	return &MongoDBRepository[T]{
		collection: collection,
		logger:     logger,
	}
}

func (r *MongoDBRepository[T]) Collection() *mongo.Collection {
	return r.collection
}

func (r *MongoDBRepository[T]) Logger() *log.Logger {
	return r.logger
}

func (r *MongoDBRepository[T]) Find(ctx context.Context, filter bson.M) ([]T, error) {
	var documents []T

	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		cursor, err := r.collection.Find(sctx, filter)

		if err != nil {
			return err
		}

		if err = cursor.All(sctx, &documents); err != nil {
			return fmt.Errorf("an error occurred reading the documents: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *MongoDBRepository[T]) FindOne(ctx context.Context, filter bson.M) (T, error) {
	var result T

	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.collection.FindOne(sctx, filter).Decode(&result)
	})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, common.ErrNotFound
		}

		return result, err
	}

	return result, nil
}

func (r *MongoDBRepository[T]) Exists(ctx context.Context, filter bson.M) (bool, error) {
	opts := options.Count().SetLimit(1)
	var count int64

	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		c, err := r.collection.CountDocuments(sctx, filter, opts)

		if err != nil {

			return err
		}

		count = c
		return nil
	})

	return count > 0, err
}
