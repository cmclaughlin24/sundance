package inmemory

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
)

type InmemoryTenantRepository struct {
	mu      sync.RWMutex
	tenants map[string]*domain.Tenant
	logger  *log.Logger
}

func NewInmemoryTenantRepository(logger *log.Logger) *InmemoryTenantRepository {
	return &InmemoryTenantRepository{
		tenants: make(map[string]*domain.Tenant),
		logger:  logger,
	}
}

func (r *InmemoryTenantRepository) Find(ctx context.Context) ([]*domain.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tenants := make([]*domain.Tenant, 0, len(r.tenants))

	for _, tenant := range r.tenants {
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

func (r *InmemoryTenantRepository) FindById(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tenant, ok := r.tenants[string(id)]

	if !ok {
		return nil, nil
	}

	return tenant, nil
}

func (r *InmemoryTenantRepository) Upsert(ctx context.Context, tenant *domain.Tenant) (*domain.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenant.UpdatedAt = time.Now()
	r.tenants[string(tenant.ID)] = tenant

	return tenant, nil
}

func (r *InmemoryTenantRepository) Remove(ctx context.Context, id domain.TenantID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tenants, string(id))

	return nil
}
