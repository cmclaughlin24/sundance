package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBRepository[T any] struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewMongoDBRepository[T any](collection *mongo.Collection, logger *slog.Logger) *MongoDBRepository[T] {
	return &MongoDBRepository[T]{
		collection: collection,
		logger:     logger,
	}
}

func (r *MongoDBRepository[T]) Collection() *mongo.Collection {
	return r.collection
}

func (r *MongoDBRepository[T]) Logger() *slog.Logger {
	return r.logger
}

func (r *MongoDBRepository[T]) Find(ctx context.Context, filter bson.M) ([]T, error) {
	r.logger.DebugContext(ctx, "mongodb find", "collection", r.collection.Name())

	var documents []T

	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		cursor, err := r.collection.Find(sctx, filter)

		if err != nil {
			return err
		}

		defer cursor.Close(sctx)

		if err = cursor.All(sctx, &documents); err != nil {
			return fmt.Errorf("an error occurred reading the documents: %w", err)
		}

		return nil
	})

	if err != nil {
		r.logger.ErrorContext(ctx, "mongodb find failed", "collection", r.collection.Name(), "error", err)
		return nil, err
	}

	return documents, nil
}

func (r *MongoDBRepository[T]) FindOne(ctx context.Context, filter bson.M) (T, error) {
	r.logger.DebugContext(ctx, "mongodb find one", "collection", r.collection.Name())

	var result T

	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.collection.FindOne(sctx, filter).Decode(&result)
	})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, common.ErrNotFound
		}

		r.logger.ErrorContext(ctx, "mongodb find one failed", "collection", r.collection.Name(), "error", err)
		return result, err
	}

	return result, nil
}

func (r *MongoDBRepository[T]) Exists(ctx context.Context, filter bson.M) (bool, error) {
	r.logger.DebugContext(ctx, "mongodb exists", "collection", r.collection.Name())

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

	if err != nil {
		r.logger.ErrorContext(ctx, "mongodb exists failed", "collection", r.collection.Name(), "error", err)
		return false, err
	}

	return count > 0, nil
}

func (r *MongoDBRepository[T]) Delete(ctx context.Context, filter bson.M) error {
	r.logger.DebugContext(ctx, "mongodb delete", "collection", r.collection.Name())

	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		result, err := r.collection.DeleteOne(sctx, filter)

		if err != nil {
			return err
		}

		if result.DeletedCount == 0 {
			return common.ErrNotFound
		}

		return nil
	})

	if err != nil && err != common.ErrNotFound {
		r.logger.ErrorContext(ctx, "mongodb delete failed", "collection", r.collection.Name(), "error", err)
	}

	return err
}
