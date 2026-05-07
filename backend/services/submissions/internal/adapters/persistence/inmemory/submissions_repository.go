package inmemory

import (
	"context"
	"log/slog"
	"sync"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
)

type InMemorySubmissionsRepository struct {
	mu          sync.RWMutex
	submissions map[string]*domain.Submission
	logger      *slog.Logger
}

func NewInMemorySubmissionsRepository(logger *slog.Logger) ports.SubmissionsRepository {
	return &InMemorySubmissionsRepository{
		submissions: make(map[string]*domain.Submission),
		logger:      logger,
	}
}

func (r *InMemorySubmissionsRepository) Find(ctx context.Context, filter *ports.FindSubmissionsFilter) ([]*domain.Submission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	submissions := make([]*domain.Submission, 0, len(r.submissions))

	for _, submission := range r.submissions {
		if filter != nil && filter.TenantID != "" && submission.TenantID != filter.TenantID {
			continue
		}
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

func (r *InMemorySubmissionsRepository) FindByIdempotencyID(ctx context.Context, idempotencyID domain.IdempotencyID) (*domain.Submission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, submission := range r.submissions {
		if submission.IdempotencyID == idempotencyID {
			return submission, nil
		}
	}

	return nil, common.ErrNotFound
}

func (r *InMemorySubmissionsRepository) Upsert(ctx context.Context, submission *domain.Submission) (*domain.Submission, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.submissions[string(submission.ID)] = submission

	return submission, nil
}
