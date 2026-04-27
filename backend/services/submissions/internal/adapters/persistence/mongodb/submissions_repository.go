package mongodb

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoDBSubmissionsRepository struct {
	base *database.MongoDBRepository[any]
}

func newMongoDBSubmissionsRepository(db *mongo.Database, logger *log.Logger) ports.SubmissionsRepository {
	base := database.NewMongoDBRepository[any](
		db.Collection("submissions"),
		logger,
	)

	return &mongoDBSubmissionsRepository{base}
}

func (r *mongoDBSubmissionsRepository) Find(ctx context.Context) ([]*domain.Submission, error) {
	return nil, nil
}

func (r *mongoDBSubmissionsRepository) FindByID(context.Context, domain.SubmissionID) (*domain.Submission, error) {
	return nil, nil
}

func (r *mongoDBSubmissionsRepository) FindByReferenceID(context.Context, domain.ReferenceID) (*domain.Submission, error) {
	return nil, nil
}
