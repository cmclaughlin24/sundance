package inmemory

import (
	"context"
	"log/slog"
	"sync"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

type InMemoryVersionsRepository struct {
	mu       sync.RWMutex
	versions map[string]map[string]*domain.Version
	logger   *slog.Logger
}

func NewInMemoryVersionsRepository(logger *slog.Logger) ports.VersionRepository {
	return &InMemoryVersionsRepository{
		versions: make(map[string]map[string]*domain.Version),
		logger:   logger,
	}
}

func (r *InMemoryVersionsRepository) Find(ctx context.Context, formID domain.FormID) ([]*domain.Version, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	formVersions, ok := r.versions[string(formID)]

	if !ok {
		return make([]*domain.Version, 0), nil
	}

	versions := make([]*domain.Version, 0, len(formVersions))

	for _, version := range formVersions {
		versions = append(versions, version)
	}

	return versions, nil
}

func (r *InMemoryVersionsRepository) FindByID(ctx context.Context, formID domain.FormID, versionID domain.VersionID) (*domain.Version, error) {
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

func (r *InMemoryVersionsRepository) FindNextVersionNumber(ctx context.Context, formID domain.FormID) (int, error) {
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

func (r *InMemoryVersionsRepository) Upsert(ctx context.Context, version *domain.Version) (*domain.Version, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	formVersions, ok := r.versions[string(version.FormID)]

	if !ok {
		formVersions = make(map[string]*domain.Version)
		r.versions[string(version.FormID)] = formVersions
	}

	formVersions[string(version.ID)] = version

	return version, nil
}
