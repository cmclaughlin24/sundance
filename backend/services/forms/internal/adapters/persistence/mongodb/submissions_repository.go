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
	submissionIndexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "idempotency_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
)

type mongoDBSubmissionsRepository struct {
	base *database.MongoDBRepository[documents.SubmissionDocument]
}

func newMongoDBSubmissionsRepository(db *mongo.Database, logger *slog.Logger) (ports.SubmissionsRepository, error) {
	base := database.NewMongoDBRepository[documents.SubmissionDocument](
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
	_, err := r.base.Collection().Indexes().CreateMany(ctx, submissionIndexes)
	return err
}

func (r *mongoDBSubmissionsRepository) Find(ctx context.Context, filter *ports.FindSubmissionsFilter) ([]*domain.Submission, error) {
	f := newSubmissionFilter(filter)
	opts := options.Find()

	if filter.Take > 0 {
		opts.SetLimit(int64(filter.Take))
	}

	docs, err := r.base.Find(ctx, f, opts)
	if err != nil {
		return nil, err
	}

	submissions := make([]*domain.Submission, 0, len(docs))
	for _, doc := range docs {
		s, err := documents.FromSubmissionDocument(&doc)
		if err != nil {
			return nil, err
		}

		submissions = append(submissions, s)
	}

	return submissions, nil
}

func (r *mongoDBSubmissionsRepository) FindJobs(ctx context.Context, filter *ports.FindSubmissionsFilter) ([]domain.SubmissionID, error) {
	f := newSubmissionFilter(filter)
	opts := options.Find()

	if filter.Take > 0 {
		opts.SetLimit(int64(filter.Take))
	}

	// TODO: Time permiting this query could be optimized by using the collection to query for only thi IDs.
	docs, err := r.base.Find(ctx, f, opts)
	if err != nil {
		return nil, err
	}

	ids := make([]domain.SubmissionID, 0, len(docs))
	for _, doc := range docs {
		ids = append(ids, domain.SubmissionID(doc.ID))
	}

	return ids, nil
}

func (r *mongoDBSubmissionsRepository) FindByID(ctx context.Context, id domain.SubmissionID) (*domain.Submission, error) {
	document, err := r.base.FindOne(ctx, bson.M{"_id": id})

	if err != nil {
		return nil, err
	}

	return documents.FromSubmissionDocument(&document)
}

func (r *mongoDBSubmissionsRepository) FindByReferenceID(ctx context.Context, id domain.ReferenceID) (*domain.Submission, error) {
	doc, err := r.base.FindOne(ctx, bson.M{"reference_id": id})

	if err != nil {
		return nil, err
	}

	return documents.FromSubmissionDocument(&doc)
}

func (r *mongoDBSubmissionsRepository) FindByIdempotencyID(ctx context.Context, id domain.IdempotencyID) (*domain.Submission, error) {
	doc, err := r.base.FindOne(ctx, bson.M{"idempotency_id": id})

	if err != nil {
		return nil, err
	}

	return documents.FromSubmissionDocument(&doc)
}

func (r *mongoDBSubmissionsRepository) Upsert(ctx context.Context, s *domain.Submission) (*domain.Submission, error) {
	r.base.Logger().DebugContext(ctx, "upsert submission", "submission_id", s.ID)

	doc, err := documents.ToSubmissionDocument(s)

	if err != nil {
		r.base.Logger().ErrorContext(ctx, "failed to convert submission to document", "submission_id", s.ID, "error", err)
		return nil, err
	}

	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result documents.SubmissionDocument
	err = r.base.WithSession(ctx, func(sctx context.Context) error {
		return r.base.Collection().FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, domain.ErrDuplicateSubmission
		}

		r.base.Logger().ErrorContext(ctx, "mongo upsert failed", "submission_id", s.ID, "error", err)
		return nil, err
	}

	return documents.FromSubmissionDocument(&result)
}

func newSubmissionFilter(filter *ports.FindSubmissionsFilter) bson.M {
	f := bson.M{}

	if filter.TenantID != "" {
		f["tenant_id"] = filter.TenantID
	}

	if len(filter.Statuses) > 0 {
		f["status"] = bson.M{"$in": filter.Statuses}
	}

	return f
}
