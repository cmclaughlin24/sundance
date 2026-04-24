package inmemory

import (
	"context"
	"log"
	"sync"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type inMemoryDataSourceRepository struct {
	mu          sync.RWMutex
	dataSources map[string]*domain.DataSource
	logger      *log.Logger
}

func newInMemoryDataSourceRepository(logger *log.Logger) ports.DataSourcesRepository {
	return &inMemoryDataSourceRepository{
		dataSources: make(map[string]*domain.DataSource),
		logger:      logger,
	}
}

func (r *inMemoryDataSourceRepository) Find(ctx context.Context, tenantID domain.TenantID) ([]*domain.DataSource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.DataSource, 0)
	prefix := string(tenantID) + "/"

	for key, ds := range r.dataSources {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			result = append(result, ds)
		}
	}

	return result, nil
}

func (r *inMemoryDataSourceRepository) FindById(ctx context.Context, tenantID domain.TenantID, id domain.DataSourceID) (*domain.DataSource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := getDataSourceKey(tenantID, id)
	ds, ok := r.dataSources[key]

	if !ok {
		return nil, common.ErrNotFound
	}

	return ds, nil
}

func (r *inMemoryDataSourceRepository) Exists(ctx context.Context, tenantID domain.TenantID, id domain.DataSourceID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := getDataSourceKey(tenantID, id)
	_, exists := r.dataSources[key]
	return exists, nil
}

func (r *inMemoryDataSourceRepository) Upsert(ctx context.Context, ds *domain.DataSource) (*domain.DataSource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := getDataSourceKey(ds.TenantID, ds.ID)
	r.dataSources[key] = ds

	return ds, nil
}

func (r *inMemoryDataSourceRepository) Remove(ctx context.Context, tenantID domain.TenantID, id domain.DataSourceID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.dataSources, getDataSourceKey(tenantID, id))

	return nil
}

func getDataSourceKey(tenantID domain.TenantID, id domain.DataSourceID) string {
	return string(tenantID) + "/" + string(id)
}
