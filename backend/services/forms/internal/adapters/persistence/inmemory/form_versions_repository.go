package inmemory

import (
	"context"
	"log/slog"
	"sync"

	"sundance/backend/pkg/common"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type inMemoryFormVersionsRepository struct {
	mu       sync.RWMutex
	versions map[string]*domain.FormVersion
	logger   *slog.Logger
}

func newInMemoryFormVersionsRepository(logger *slog.Logger) ports.FormVersionRepository {
	return &inMemoryFormVersionsRepository{
		versions: make(map[string]*domain.FormVersion),
		logger:   logger,
	}
}

func (r *inMemoryFormVersionsRepository) Find(ctx context.Context, formID domain.FormID) ([]*domain.FormVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	versions := make([]*domain.FormVersion, 0)

	for _, v := range r.versions {
		if v.FormID == formID {
			versions = append(versions, v)
		}
	}

	return versions, nil
}

func (r *inMemoryFormVersionsRepository) FindByID(ctx context.Context, versionID domain.FormVersionID) (*domain.FormVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	version, ok := r.versions[string(versionID)]

	if !ok {
		return nil, common.ErrNotFound
	}

	return version, nil
}

func (r *inMemoryFormVersionsRepository) FindNextVersionNumber(ctx context.Context, formID domain.FormID) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	maxVersion := 0
	for _, v := range r.versions {
		if v.FormID == formID && v.Version > maxVersion {
			maxVersion = v.Version
		}
	}

	return maxVersion + 1, nil
}

func (r *inMemoryFormVersionsRepository) Upsert(ctx context.Context, version *domain.FormVersion) (*domain.FormVersion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.versions {
		if existing.FormID == version.FormID &&
			existing.ID != version.ID &&
			existing.Version == version.Version {
			return nil, domain.ErrDuplicateVersion
		}
	}

	r.versions[string(version.ID)] = version

	return version, nil
}
