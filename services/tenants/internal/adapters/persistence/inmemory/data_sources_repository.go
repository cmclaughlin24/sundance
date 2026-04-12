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

type InmemoryDataSourceRepository struct {
	mu          sync.RWMutex
	dataSources map[string]*domain.DataSource
	logger      *log.Logger
}

func NewInmemoryDataSourceRepository(logger *log.Logger) *InmemoryDataSourceRepository {
	return &InmemoryDataSourceRepository{
		dataSources: make(map[string]*domain.DataSource),
		logger:      logger,
	}
}

func (r *InmemoryDataSourceRepository) Find(ctx context.Context, tenantID domain.TenantID) ([]*domain.DataSource, error) {
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

func (r *InmemoryDataSourceRepository) FindById(ctx context.Context, tenantID domain.TenantID, id domain.DataSourceID) (*domain.DataSource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := getDataSourceKey(tenantID, id)
	ds, ok := r.dataSources[key]

	if !ok {
		return nil, nil
	}

	return ds, nil
}

func (r *InmemoryDataSourceRepository) Upsert(ctx context.Context, ds *domain.DataSource) (*domain.DataSource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	key := getDataSourceKey(ds.TenantID, ds.ID)

	if ds.ID == "" {
		ds.ID = domain.DataSourceID(uuid.New().String())
		key = getDataSourceKey(ds.TenantID, ds.ID)
		ds.CreatedAt = now
	} else {
		existing, exists := r.dataSources[key]

		if !exists {
			return nil, fmt.Errorf("data source %q not found", ds.ID)
		}

		ds.CreatedAt = existing.CreatedAt
	}

	ds.UpdatedAt = now
	r.dataSources[key] = ds

	return ds, nil
}

func (r *InmemoryDataSourceRepository) Remove(ctx context.Context, tenantID domain.TenantID, id domain.DataSourceID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.dataSources, getDataSourceKey(tenantID, id))

	return nil
}

func getDataSourceKey(tenantID domain.TenantID, id domain.DataSourceID) string {
	return string(tenantID) + "/" + string(id)
}
