package ports

import (
	"context"

	"sundance/backend/pkg/common/stratreg"
	"sundance/backend/pkg/database"
	"sundance/backend/services/tenants/internal/core/domain"
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
	DataLake DataLakeClient
	Lookups  LookupClient
}

type DataLakeClient interface {
	Query(context.Context, domain.DataLakeDataSourceAttributes, map[string]any) ([]*domain.Lookup, error)
}

type LookupClient interface {
	FetchLookups(context.Context, domain.DataSourceHTTPRequest, map[string]any) ([]map[string]any, error)
}

type Strategies struct {
	Lookups LookupStrategyRegistry
}

type LookupStrategy interface {
	Lookup(context.Context, *domain.DataSource, map[string]any) ([]*domain.Lookup, error)
}

type LookupStrategyRegistry = stratreg.StrategyRegistry[domain.DataSourceType, LookupStrategy]
