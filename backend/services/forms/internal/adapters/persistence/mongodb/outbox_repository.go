package mongodb

import (
	"context"
	"errors"
	"log/slog"
	"sundance/backend/pkg/database"
	"sundance/backend/services/forms/internal/adapters/persistence/mongodb/documents"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	outboxIndexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "attempts", Value: 1},
				{Key: "created_at", Value: 1},
			},
			Options: options.Index(),
		},
	}
)

type mongoDBOutboxRepository struct {
	base *database.MongoDBRepository[documents.EventDocument]
}

func newMongoDBOutboxRepository(db *mongo.Database, logger *slog.Logger) (ports.OutboxRepository, error) {
	base := database.NewMongoDBRepository[documents.EventDocument](
		db.Collection("outbox"),
		logger,
	)
	repository := &mongoDBOutboxRepository{base}

	if err := repository.migrate(context.Background()); err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *mongoDBOutboxRepository) migrate(ctx context.Context) error {
	_, err := r.base.Collection().Indexes().CreateMany(ctx, outboxIndexes)
	return err
}

func (r *mongoDBOutboxRepository) Claim(ctx context.Context, o ports.ClaimEventsOptions) ([]*domain.Event, error) {
	now := Now()
	filters := bson.M{
		"attempts":   bson.M{"$lt": o.RetryLimit},
		"created_at": bson.M{"$gte": o.CreatedAfter},
		"$or": bson.A{
			bson.M{"status": bson.M{"$in": []domain.EventStatus{domain.EventStatusPending, domain.EventStatusError}}},
			bson.M{"status": domain.EventStatusProcessing, "locked_until": bson.M{"$lt": now}},
		},
	}
	update := bson.M{
		"$set": bson.M{
			"status":       domain.EventStatusProcessing,
			"locked_until": now.Add(o.LeaseDuration),
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	docs := make([]documents.EventDocument, 0)

	err := r.base.WithSession(ctx, func(sctx context.Context) error {
		for range o.BatchSize {
			var doc documents.EventDocument

			err := r.base.Collection().FindOneAndUpdate(sctx, filters, update, opts).Decode(&doc)
			if errors.Is(err, mongo.ErrNoDocuments) {
				break
			} else if err != nil {
				return err
			}

			docs = append(docs, doc)
		}

		return nil
	})

	if err != nil {
		r.base.Logger().ErrorContext(ctx, "mongo find and update failed", "error", err)
		return nil, err
	}

	events := make([]*domain.Event, 0, len(docs))
	for _, doc := range docs {
		events = append(events, documents.FromEventDocument(&doc))
	}

	return events, nil
}

func (r *mongoDBOutboxRepository) Upsert(ctx context.Context, e *domain.Event) (*domain.Event, error) {
	r.base.Logger().DebugContext(ctx, "upsert event", "event", e.ID)

	doc := documents.ToEventDocument(*e)
	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc, "$unset": bson.M{"locked_until": ""}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result documents.EventDocument
	err := r.base.WithSession(ctx, func(sctx context.Context) error {
		return r.base.Collection().FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		r.base.Logger().ErrorContext(ctx, "mongo upsert failed", "event_id", e.ID, "error", err)
		return nil, err
	}

	return documents.FromEventDocument(&result), nil
}
