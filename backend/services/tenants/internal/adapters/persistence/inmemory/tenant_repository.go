package inmemory

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"github.com/google/uuid"
)

type inMemoryTenantRepository struct {
	mu      sync.RWMutex
	tenants map[string]*domain.Tenant
	logger  *log.Logger
}

func newInMemoryTenantRepository(logger *log.Logger) ports.TenantsRepository {
	return &inMemoryTenantRepository{
		tenants: make(map[string]*domain.Tenant),
		logger:  logger,
	}
}

func (r *inMemoryTenantRepository) Find(ctx context.Context) ([]*domain.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tenants := make([]*domain.Tenant, 0, len(r.tenants))

	for _, tenant := range r.tenants {
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

func (r *inMemoryTenantRepository) FindById(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tenant, ok := r.tenants[string(id)]

	if !ok {
		return nil, common.ErrNotFound
	}

	return tenant, nil
}

func (r *inMemoryTenantRepository) Exists(ctx context.Context, id domain.TenantID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.tenants[string(id)]
	return exists, nil
}

func (r *inMemoryTenantRepository) Upsert(ctx context.Context, tenant *domain.Tenant) (*domain.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	if tenant.ID == "" {
		tenant.ID = domain.TenantID(uuid.New().String())
		tenant.CreatedAt = now
	} else {
		existing, exists := r.tenants[string(tenant.ID)]

		if !exists {
			return nil, common.ErrNotFound
		}

		tenant.CreatedAt = existing.CreatedAt
	}

	tenant.UpdatedAt = now
	r.tenants[string(tenant.ID)] = tenant

	return tenant, nil
}

func (r *inMemoryTenantRepository) Remove(ctx context.Context, id domain.TenantID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tenants, string(id))

	return nil
}
