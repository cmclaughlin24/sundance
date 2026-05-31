package inmemory

import (
	"context"
	"log/slog"
	"sync"

	"sundance/backend/pkg/common"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type inMemoryCanonicalTagRepository struct {
	mu     sync.RWMutex
	tags   map[string]*domain.CanonicalTag
	logger *slog.Logger
}

func newInMemoryCanonicalTagRepository(logger *slog.Logger) ports.CanonicalTagRepository {
	return &inMemoryCanonicalTagRepository{
		tags:   make(map[string]*domain.CanonicalTag),
		logger: logger,
	}
}

func (r *inMemoryCanonicalTagRepository) Find(ctx context.Context, filters ports.CanonicalTagFilters) ([]*domain.CanonicalTag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tags := make([]*domain.CanonicalTag, 0, len(r.tags))

	for _, tag := range r.tags {
		if filters.TenantID != "" && tag.TenantID != filters.TenantID {
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *inMemoryCanonicalTagRepository) FindByID(ctx context.Context, id domain.CanonicalTagID) (*domain.CanonicalTag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tag, ok := r.tags[string(id)]

	if !ok {
		return nil, common.ErrNotFound
	}

	return tag, nil
}

func (r *inMemoryCanonicalTagRepository) Upsert(ctx context.Context, tag *domain.CanonicalTag) (*domain.CanonicalTag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tags[string(tag.ID)] = tag

	return tag, nil
}

func (r *inMemoryCanonicalTagRepository) Delete(ctx context.Context, id domain.CanonicalTagID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tags[string(id)]; !ok {
		return common.ErrNotFound
	}

	delete(r.tags, string(id))

	return nil
}
