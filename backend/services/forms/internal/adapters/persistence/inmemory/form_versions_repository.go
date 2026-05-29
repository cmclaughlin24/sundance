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
	versions map[string]map[string]*domain.FormVersion
	logger   *slog.Logger
}

func newInMemoryFormVersionsRepository(logger *slog.Logger) ports.FormVersionRepository {
	return &inMemoryFormVersionsRepository{
		versions: make(map[string]map[string]*domain.FormVersion),
		logger:   logger,
	}
}

func (r *inMemoryFormVersionsRepository) Find(ctx context.Context, formID domain.FormID) ([]*domain.FormVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	formVersions, ok := r.versions[string(formID)]

	if !ok {
		return make([]*domain.FormVersion, 0), nil
	}

	versions := make([]*domain.FormVersion, 0, len(formVersions))

	for _, version := range formVersions {
		versions = append(versions, version)
	}

	return versions, nil
}

func (r *inMemoryFormVersionsRepository) FindByID(ctx context.Context, formID domain.FormID, versionID domain.FormVersionID) (*domain.FormVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	formVersions, ok := r.versions[string(formID)]

	if !ok {
		return nil, common.ErrNotFound
	}

	version, ok := formVersions[string(versionID)]

	if !ok {
		return nil, common.ErrNotFound
	}

	return version, nil
}

func (r *inMemoryFormVersionsRepository) FindNextVersionNumber(ctx context.Context, formID domain.FormID) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	versions, ok := r.versions[string(formID)]

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

func (r *inMemoryFormVersionsRepository) Upsert(ctx context.Context, version *domain.FormVersion) (*domain.FormVersion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	formVersions, ok := r.versions[string(version.FormID)]

	if !ok {
		formVersions = make(map[string]*domain.FormVersion)
		r.versions[string(version.FormID)] = formVersions
	}

	for _, existing := range formVersions {
		if existing.ID != version.ID && existing.Version == version.Version {
			return nil, domain.ErrDuplicateVersion
		}
	}

	formVersions[string(version.ID)] = version

	return version, nil
}
