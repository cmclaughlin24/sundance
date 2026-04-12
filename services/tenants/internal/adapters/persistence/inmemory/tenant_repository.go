package inmemory

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
	"github.com/google/uuid"
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

	now := time.Now()

	if tenant.ID == "" {
		tenant.ID = domain.TenantID(uuid.New().String())
		tenant.CreatedAt = now
	} else {
		existing, exists := r.tenants[string(tenant.ID)]

		if !exists {
			// TODO: Implement improved error handling with custom types
			return nil, fmt.Errorf("tenant %q not found", tenant.ID)
		}

		tenant.CreatedAt = existing.CreatedAt
	}

	tenant.UpdatedAt = now
	r.tenants[string(tenant.ID)] = tenant

	return tenant, nil
}

func (r *InmemoryTenantRepository) Remove(ctx context.Context, id domain.TenantID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tenants, string(id))

	return nil
}
