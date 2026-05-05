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

type mongoDBSubmissionsRepository struct {
	base *database.MongoDBRepository[submissionDocument]
}

func newMongoDBSubmissionsRepository(db *mongo.Database, logger *slog.Logger) ports.SubmissionsRepository {
	base := database.NewMongoDBRepository[submissionDocument](
		db.Collection("submissions"),
		logger,
	)

	return &mongoDBSubmissionsRepository{base}
}

func (r *mongoDBSubmissionsRepository) Find(ctx context.Context, filter *ports.FindSubmissionsFilter) ([]*domain.Submission, error) {
	documents, err := r.base.Find(ctx, bson.M{"tenant_id": filter.TenantID})

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

func (r *mongoDBSubmissionsRepository) Upsert(ctx context.Context, s *domain.Submission) (*domain.Submission, error) {
	doc, err := toSubmissionDocument(s)

	if err != nil {
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
		return nil, err
	}

	return fromSubmissionDocument(&result), nil
}
