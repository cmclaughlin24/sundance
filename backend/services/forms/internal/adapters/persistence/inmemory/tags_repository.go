package inmemory

import (
	"context"
	"log/slog"
	"sync"

	"sundance/backend/pkg/common"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type inMemoryTagRepository struct {
	mu     sync.RWMutex
	tags   map[string]*domain.Tag
	logger *slog.Logger
}

func newInMemoryTagRepository(logger *slog.Logger) ports.TagsRepository {
	return &inMemoryTagRepository{
		tags:   make(map[string]*domain.Tag),
		logger: logger,
	}
}

func (r *inMemoryTagRepository) Find(ctx context.Context, filters ports.TagFilters) ([]*domain.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tags := make([]*domain.Tag, 0, len(r.tags))

	for _, tag := range r.tags {
		if filters.TenantID != "" && tag.TenantID != filters.TenantID {
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *inMemoryTagRepository) FindByIDs(ctx context.Context, ids []domain.TagID) ([]*domain.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tags := make([]*domain.Tag, 0, len(ids))

	for _, id := range ids {
		if tag, ok := r.tags[string(id)]; ok {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

func (r *inMemoryTagRepository) FindByID(ctx context.Context, id domain.TagID) (*domain.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tag, ok := r.tags[string(id)]

	if !ok {
		return nil, common.ErrNotFound
	}

	return tag, nil
}

func (r *inMemoryTagRepository) Upsert(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tags[string(tag.ID)] = tag

	return tag, nil
}

func (r *inMemoryTagRepository) Delete(ctx context.Context, id domain.TagID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tags[string(id)]; !ok {
		return common.ErrNotFound
	}

	delete(r.tags, string(id))

	return nil
}
