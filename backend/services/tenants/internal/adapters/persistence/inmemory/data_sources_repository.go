package inmemory

import (
	"context"
	"log/slog"
	"slices"
	"sync"

	"sundance/backend/pkg/common"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type inMemoryDataSourceRepository struct {
	mu          sync.RWMutex
	dataSources map[string]*domain.DataSource
	logger      *slog.Logger
}

func newInMemoryDataSourceRepository(logger *slog.Logger) ports.DataSourcesRepository {
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

func (r *inMemoryDataSourceRepository) FindByID(ctx context.Context, tenantID domain.TenantID, id domain.DataSourceID) (*domain.DataSource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := getDataSourceKey(tenantID, id)
	ds, ok := r.dataSources[key]

	if !ok {
		return nil, common.ErrNotFound
	}

	return ds, nil
}

func (r *inMemoryDataSourceRepository) FindJobs(ctx context.Context, filters *ports.FindDataSourceJobsFilter) ([]*domain.DataSource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.DataSource, 0)
	for _, ds := range r.dataSources {
		if !slices.Contains(filters.Types, ds.Type) {
			continue
		}

		if attr, ok := ds.Attributes.(*domain.ScheduledDataSourceAttributes); ok {
			if !attr.ExpirationDate.IsZero() && attr.ExpirationDate.After(filters.ExpiredAtOrBefore) {
				continue
			}
		}

		result = append(result, ds)

		if filters.Limit > 0 && len(result) >= filters.Limit {
			break
		}
	}

	return result, nil
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

func (r *inMemoryDataSourceRepository) Delete(ctx context.Context, tenantID domain.TenantID, id domain.DataSourceID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.dataSources, getDataSourceKey(tenantID, id))

	return nil
}

func (r *inMemoryDataSourceRepository) DeleteAll(ctx context.Context, tenantID domain.TenantID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	prefix := string(tenantID) + "/"

	for key := range r.dataSources {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			delete(r.dataSources, key)
		}
	}

	return nil
}

func getDataSourceKey(tenantID domain.TenantID, id domain.DataSourceID) string {
	return string(tenantID) + "/" + string(id)
}
