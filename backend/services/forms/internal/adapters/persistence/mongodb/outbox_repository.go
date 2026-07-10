package mongodb

import (
	"context"
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

func (r *mongoDBOutboxRepository) Find(ctx context.Context, filters ports.FindEventsFilter) ([]*domain.Event, error) {
	opts := options.Find()
	if filters.Take > 0 {
		opts.SetLimit(int64(filters.Take))
	}

	f := bson.M{
		"attempts":   bson.M{"$lt": filters.RetryLimit},
		"created_at": bson.M{"$gte": filters.CreatedAfter},
	}

	if len(filters.Statuses) > 0 {
		f["status"] = bson.M{"$in": filters.Statuses}
	}

	docs, err := r.base.Find(ctx, f, opts)
	if err != nil {
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
	update := bson.M{"$set": doc}
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
