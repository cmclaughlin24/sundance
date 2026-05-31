package inmemory

import (
	"context"
	"log/slog"
	"sync"

	"sundance/backend/pkg/common"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type inMemoryTagVersionsRepository struct {
	mu       sync.RWMutex
	versions map[string]map[string]*domain.TagVersion
	logger   *slog.Logger
}

func newInMemoryTagVersionsRepository(logger *slog.Logger) ports.TagVersionsRepository {
	return &inMemoryTagVersionsRepository{
		versions: make(map[string]map[string]*domain.TagVersion),
		logger:   logger,
	}
}

func (r *inMemoryTagVersionsRepository) Find(ctx context.Context, filters ports.TagVersionFilters) ([]*domain.TagVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if filters.TagID != "" {
		tagVersions, ok := r.versions[string(filters.TagID)]

		if !ok {
			return make([]*domain.TagVersion, 0), nil
		}

		versions := make([]*domain.TagVersion, 0, len(tagVersions))

		for _, version := range tagVersions {
			versions = append(versions, version)
		}

		return versions, nil
	}

	versions := make([]*domain.TagVersion, 0)

	for _, tagVersions := range r.versions {
		for _, version := range tagVersions {
			versions = append(versions, version)
		}
	}

	return versions, nil
}

func (r *inMemoryTagVersionsRepository) FindByID(ctx context.Context, id domain.TagVersionID) (*domain.TagVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, tagVersions := range r.versions {
		if version, ok := tagVersions[string(id)]; ok {
			return version, nil
		}
	}

	return nil, common.ErrNotFound
}

func (r *inMemoryTagVersionsRepository) FindNextVersionNumber(ctx context.Context, tagID domain.TagID) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	versions, ok := r.versions[string(tagID)]

	if !ok {
		return 1, nil
	}

	maxVersion := 0
	for _, version := range versions {
		if version.Version > maxVersion {
			maxVersion = version.Version
		}
	}

	return maxVersion + 1, nil
}

func (r *inMemoryTagVersionsRepository) Upsert(ctx context.Context, version *domain.TagVersion) (*domain.TagVersion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tagVersions, ok := r.versions[string(version.TagID)]

	if !ok {
		tagVersions = make(map[string]*domain.TagVersion)
		r.versions[string(version.TagID)] = tagVersions
	}

	for _, existing := range tagVersions {
		if existing.ID != version.ID && existing.Version == version.Version {
			return nil, domain.ErrDuplicateTagVersion
		}
	}

	tagVersions[string(version.ID)] = version

	return version, nil
}
