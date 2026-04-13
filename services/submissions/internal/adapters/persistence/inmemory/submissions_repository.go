package inmemory

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/submissions/internal/core/domain"
)

type InMemorySubmissionsRepository struct {
	logger *log.Logger
}

func NewInMemorySubmissionsRepository(logger *log.Logger) *InMemorySubmissionsRepository {
	return &InMemorySubmissionsRepository{
		logger: logger,
	}
}

func (r *InMemorySubmissionsRepository) Find(context.Context) ([]*domain.Submission, error) {
	return nil, nil
}

func (r *InMemorySubmissionsRepository) FindById(context.Context, domain.SubmissionID) (*domain.Submission, error) {
	return nil, nil
}

func (r *InMemorySubmissionsRepository) FindByReferenceId(context.Context, domain.ReferenceID) (*domain.Submission, error) {
	return nil, nil
}
