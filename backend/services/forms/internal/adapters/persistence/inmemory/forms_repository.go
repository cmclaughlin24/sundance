package inmemory

import (
	"context"
	"log/slog"
	"sync"

	"sundance/backend/pkg/common"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type inMemoryFormsRepository struct {
	mu     sync.RWMutex
	forms  map[string]*domain.Form
	logger *slog.Logger
}

func newInMemoryFormsRepository(logger *slog.Logger) ports.FormsRepository {
	return &inMemoryFormsRepository{
		forms:  make(map[string]*domain.Form),
		logger: logger,
	}
}

func (r *inMemoryFormsRepository) Find(ctx context.Context, f *ports.FormFilters) ([]*domain.Form, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	forms := make([]*domain.Form, 0, len(r.forms))

	for _, form := range r.forms {
		if f != nil && f.TenantID != "" && form.TenantID != f.TenantID {
			continue
		}
		forms = append(forms, form)
	}

	return forms, nil
}

func (r *inMemoryFormsRepository) FindByID(ctx context.Context, id domain.FormID) (*domain.Form, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	form, ok := r.forms[string(id)]

	if !ok {
		return nil, common.ErrNotFound
	}

	return form, nil
}

func (r *inMemoryFormsRepository) Upsert(ctx context.Context, form *domain.Form) (*domain.Form, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.forms[string(form.ID)] = form

	return form, nil
}

func (r *inMemoryFormsRepository) Delete(ctx context.Context, id domain.FormID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.forms[string(id)]; !ok {
		return common.ErrNotFound
	}

	delete(r.forms, string(id))

	return nil
}
