package inmemory

import (
	"context"
	"log"
	"sync"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
)

type InMemorySubmissionsRepository struct {
	mu          sync.RWMutex
	submissions map[string]*domain.Submission
	logger      *log.Logger
}

func NewInMemorySubmissionsRepository(logger *log.Logger) ports.SubmissionsRepository {
	return &InMemorySubmissionsRepository{
		submissions: make(map[string]*domain.Submission),
		logger:      logger,
	}
}

func (r *InMemorySubmissionsRepository) Find(ctx context.Context) ([]*domain.Submission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	submissions := make([]*domain.Submission, 0, len(r.submissions))

	for _, submission := range r.submissions {
		submissions = append(submissions, submission)
	}

	return submissions, nil
}

func (r *InMemorySubmissionsRepository) FindByID(ctx context.Context, id domain.SubmissionID) (*domain.Submission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	submission, ok := r.submissions[string(id)]

	if !ok {
		return nil, common.ErrNotFound
	}

	return submission, nil
}

func (r *InMemorySubmissionsRepository) FindByReferenceID(ctx context.Context, referenceID domain.ReferenceID) (*domain.Submission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, submission := range r.submissions {
		if submission.ReferenceID == referenceID {
			return submission, nil
		}
	}

	return nil, common.ErrNotFound
}
