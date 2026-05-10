package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/stratreg"
	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
)

type Repository struct {
	Database    database.Database
	Tenants     TenantsRepository
	DataSources DataSourcesRepository
}

type TenantsRepository interface {
	Find(context.Context) ([]*domain.Tenant, error)
	FindByID(context.Context, domain.TenantID) (*domain.Tenant, error)
	Exists(context.Context, domain.TenantID) (bool, error)
	Upsert(context.Context, *domain.Tenant) (*domain.Tenant, error)
	Delete(context.Context, domain.TenantID) error
}

type DataSourcesRepository interface {
	Find(context.Context, domain.TenantID) ([]*domain.DataSource, error)
	FindByID(context.Context, domain.TenantID, domain.DataSourceID) (*domain.DataSource, error)
	FindJobs(context.Context, *FindDataSourceJobsFilter) ([]*domain.DataSource, error)
	Exists(context.Context, domain.TenantID, domain.DataSourceID) (bool, error)
	Upsert(context.Context, *domain.DataSource) (*domain.DataSource, error)
	Delete(context.Context, domain.TenantID, domain.DataSourceID) error
	DeleteAll(context.Context, domain.TenantID) error
}

type Clients struct {
	Lookups LookupClient
}

type LookupClient interface {
	FetchLookups(ctx context.Context, method, url string, headers map[string]string) ([]*domain.Lookup, error)
}

type Strategies struct {
	Lookups LookupStrategyRegistry
}

type LookupStrategy interface {
	Lookup(context.Context, *domain.DataSource) ([]*domain.Lookup, error)
}

type LookupStrategyRegistry = stratreg.StrategyRegistry[domain.DataSourceType, LookupStrategy]
