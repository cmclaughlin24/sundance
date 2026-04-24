package mongodb

import (
	"context"
	"errors"
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongodbBaseRepository[T any] struct {
	collection *mongo.Collection
	logger     *log.Logger
}

func (r *mongodbBaseRepository[T]) findById(ctx context.Context, filter bson.M) (T, error) {
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

func (r *mongodbBaseRepository[T]) exists(ctx context.Context, filter bson.M) (bool, error) {
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
