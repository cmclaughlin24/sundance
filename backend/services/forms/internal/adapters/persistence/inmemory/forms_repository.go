package inmemory

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
	"github.com/google/uuid"
)

type InMemoryFormsRepository struct {
	mu       sync.RWMutex
	forms    map[string]*domain.Form
	versions map[string]map[string]*domain.Version
	logger   *log.Logger
}

func NewInMemoryFormsRepository(logger *log.Logger) ports.FormsRepository{
	return &InMemoryFormsRepository{
		forms:    make(map[string]*domain.Form),
		versions: make(map[string]map[string]*domain.Version),
		logger:   logger,
	}
}

func (r *InMemoryFormsRepository) Find(ctx context.Context) ([]*domain.Form, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	forms := make([]*domain.Form, 0, len(r.forms))

	for _, form := range r.forms {
		forms = append(forms, form)
	}

	return forms, nil
}

func (r *InMemoryFormsRepository) FindById(ctx context.Context, id domain.FormID) (*domain.Form, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	form, ok := r.forms[string(id)]

	if !ok {
		return nil, common.ErrNotFound
	}

	return form, nil
}

func (r *InMemoryFormsRepository) Create(ctx context.Context, form *domain.Form) (*domain.Form, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	if form.ID == "" {
		form.ID = domain.FormID(uuid.New().String())
		form.CreatedAt = now
	}

	form.UpdatedAt = now
	r.forms[string(form.ID)] = form

	return form, nil
}

func (r *InMemoryFormsRepository) Update(ctx context.Context, form *domain.Form) (*domain.Form, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.forms[string(form.ID)]

	if !exists {
		return nil, common.ErrNotFound
	}

	form.CreatedAt = existing.CreatedAt
	form.UpdatedAt = time.Now()
	r.forms[string(form.ID)] = form

	return form, nil
}

func (r *InMemoryFormsRepository) FindVersions(ctx context.Context, formID domain.FormID) ([]*domain.Version, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	formVersions, ok := r.versions[string(formID)]

	if !ok {
		return nil, common.ErrNotFound
	}

	versions := make([]*domain.Version, 0, len(formVersions))

	for _, version := range formVersions {
		versions = append(versions, version)
	}

	return versions, nil
}

func (r *InMemoryFormsRepository) FindVersion(ctx context.Context, formID domain.FormID, versionID domain.VersionID) (*domain.Version, error) {
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

func (r *InMemoryFormsRepository) FindNextVersionNumber(ctx context.Context, formID domain.FormID) (int, error) {
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

func (r *InMemoryFormsRepository) CreateVersion(ctx context.Context, version *domain.Version) (*domain.Version, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if version.ID == "" {
		version.ID = domain.VersionID(uuid.New().String())
	}

	formVersions, ok := r.versions[string(version.FormID)]

	if !ok {
		formVersions = make(map[string]*domain.Version)
		r.versions[string(version.FormID)] = formVersions
	}

	formVersions[string(version.ID)] = version

	return version, nil
}

func (r *InMemoryFormsRepository) UpdateVersion(ctx context.Context, version *domain.Version) (*domain.Version, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	formVersions, ok := r.versions[string(version.FormID)]

	if !ok {
		return nil, common.ErrNotFound
	}

	_, ok = formVersions[string(version.ID)]

	if !ok {
		return nil, common.ErrNotFound
	}

	formVersions[string(version.ID)] = version

	return version, nil
}
