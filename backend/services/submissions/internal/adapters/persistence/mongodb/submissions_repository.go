package mongodb

import (
	"context"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "idempotency_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
)

type mongoDBSubmissionsRepository struct {
	base *database.MongoDBRepository[submissionDocument]
}

func newMongoDBSubmissionsRepository(db *mongo.Database, logger *slog.Logger) (ports.SubmissionsRepository, error) {
	base := database.NewMongoDBRepository[submissionDocument](
		db.Collection("submissions"),
		logger,
	)
	repository := &mongoDBSubmissionsRepository{base}
	
	if err := repository.migrate(context.Background()); err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *mongoDBSubmissionsRepository) migrate(ctx context.Context) error {
	_, err := r.base.Collection().Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *mongoDBSubmissionsRepository) Find(ctx context.Context, filter *ports.FindSubmissionsFilter) ([]*domain.Submission, error) {
	f := bson.M{}

	if filter.TenantID != "" {
		f["tenant_id"] = filter.TenantID
	}

	if len(filter.Statuses) > 0 {
		f["status"] = bson.M{"$in": filter.Statuses}
	}

	documents, err := r.base.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	submissions := make([]*domain.Submission, 0, len(documents))
	for _, doc := range documents {
		submissions = append(submissions, fromSubmissionDocument(&doc))
	}

	return submissions, nil
}

func (r *mongoDBSubmissionsRepository) FindByID(ctx context.Context, id domain.SubmissionID) (*domain.Submission, error) {
	document, err := r.base.FindOne(ctx, bson.M{"_id": id})

	if err != nil {
		return nil, err
	}

	return fromSubmissionDocument(&document), nil
}

func (r *mongoDBSubmissionsRepository) FindByReferenceID(ctx context.Context, id domain.ReferenceID) (*domain.Submission, error) {
	document, err := r.base.FindOne(ctx, bson.M{"reference_id": id})

	if err != nil {
		return nil, err
	}

	return fromSubmissionDocument(&document), nil
}

func (r *mongoDBSubmissionsRepository) FindByIdempotencyID(ctx context.Context, id domain.IdempotencyID) (*domain.Submission, error) {
	document, err := r.base.FindOne(ctx, bson.M{"idempotency_id": id})

	if err != nil {
		return nil, err
	}

	return fromSubmissionDocument(&document), nil
}

func (r *mongoDBSubmissionsRepository) Upsert(ctx context.Context, s *domain.Submission) (*domain.Submission, error) {
	r.base.Logger().DebugContext(ctx, "upsert submission", "submission_id", s.ID)

	doc, err := toSubmissionDocument(s)

	if err != nil {
		r.base.Logger().ErrorContext(ctx, "failed to convert submission to document", "submission_id", s.ID, "error", err)
		return nil, err
	}

	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result submissionDocument
	err = mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.base.Collection().FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, domain.ErrDuplicateSubmission
		}

		r.base.Logger().ErrorContext(ctx, "mongo upsert failed", "submission_id", s.ID, "error", err)
		return nil, err
	}

	return fromSubmissionDocument(&result), nil
}
