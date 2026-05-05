package inmemory

import (
	"context"
	"log/slog"
	"sync"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

type InMemoryFormsRepository struct {
	mu     sync.RWMutex
	forms  map[string]*domain.Form
	logger *slog.Logger
}

func NewInMemoryFormsRepository(logger *slog.Logger) ports.FormsRepository {
	return &InMemoryFormsRepository{
		forms:  make(map[string]*domain.Form),
		logger: logger,
	}
}

func (r *InMemoryFormsRepository) Find(ctx context.Context, f *ports.FormFilters) ([]*domain.Form, error) {
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

func (r *InMemoryFormsRepository) FindByID(ctx context.Context, id domain.FormID) (*domain.Form, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	form, ok := r.forms[string(id)]

	if !ok {
		return nil, common.ErrNotFound
	}

	return form, nil
}

func (r *InMemoryFormsRepository) Upsert(ctx context.Context, form *domain.Form) (*domain.Form, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.forms[string(form.ID)] = form

	return form, nil
}

func (r *InMemoryFormsRepository) Delete(ctx context.Context, id domain.FormID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.forms[string(id)]; !ok {
		return common.ErrNotFound
	}

	delete(r.forms, string(id))

	return nil
}
