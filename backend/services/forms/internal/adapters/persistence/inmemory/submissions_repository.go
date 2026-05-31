package inmemory

import (
	"context"
	"log/slog"
	"slices"
	"sync"

	"sundance/backend/pkg/common"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type inMemorySubmissionsRepository struct {
	mu          sync.RWMutex
	submissions map[string]*domain.Submission
	logger      *slog.Logger
}

func newInMemorySubmissionsRepository(logger *slog.Logger) ports.SubmissionsRepository {
	return &inMemorySubmissionsRepository{
		submissions: make(map[string]*domain.Submission),
		logger:      logger,
	}
}

func (r *inMemorySubmissionsRepository) Find(ctx context.Context, filter *ports.FindSubmissionsFilter) ([]*domain.Submission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	submissions := make([]*domain.Submission, 0, len(r.submissions))

	for _, submission := range r.submissions {
		if filter != nil && filter.TenantID != "" && submission.TenantID != filter.TenantID {
			continue
		}

		if filter != nil && len(filter.Statuses) > 0 && !slices.Contains(filter.Statuses, submission.Status) {
			continue
		}

		submissions = append(submissions, submission)
	}

	if filter != nil && filter.Take > 0 && len(submissions) > filter.Take {
		submissions = submissions[:filter.Take]
	}

	return submissions, nil
}

func (r *inMemorySubmissionsRepository) FindByID(ctx context.Context, id domain.SubmissionID) (*domain.Submission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	submission, ok := r.submissions[string(id)]

	if !ok {
		return nil, common.ErrNotFound
	}

	return submission, nil
}

func (r *inMemorySubmissionsRepository) FindByReferenceID(ctx context.Context, referenceID domain.ReferenceID) (*domain.Submission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, submission := range r.submissions {
		if submission.ReferenceID == referenceID {
			return submission, nil
		}
	}

	return nil, common.ErrNotFound
}

func (r *inMemorySubmissionsRepository) FindByIdempotencyID(ctx context.Context, idempotencyID domain.IdempotencyID) (*domain.Submission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, submission := range r.submissions {
		if submission.IdempotencyID == idempotencyID {
			return submission, nil
		}
	}

	return nil, common.ErrNotFound
}

func (r *inMemorySubmissionsRepository) FindJobs(ctx context.Context, filter *ports.FindSubmissionsFilter) ([]domain.SubmissionID, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]domain.SubmissionID, 0, len(r.submissions))

	for _, submission := range r.submissions {
		if filter != nil && filter.TenantID != "" && submission.TenantID != filter.TenantID {
			continue
		}

		if filter != nil && len(filter.Statuses) > 0 && !slices.Contains(filter.Statuses, submission.Status) {
			continue
		}

		ids = append(ids, submission.ID)
	}

	if filter != nil && filter.Take > 0 && len(ids) > filter.Take {
		ids = ids[:filter.Take]
	}

	return ids, nil
}

func (r *inMemorySubmissionsRepository) Upsert(ctx context.Context, submission *domain.Submission) (*domain.Submission, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.submissions[string(submission.ID)] = submission

	return submission, nil
}
