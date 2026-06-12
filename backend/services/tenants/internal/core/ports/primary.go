package ports

import (
	"context"

	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports/commands"
)

type API struct {
	Tenants        TenantsAPI
	DataSources    DataSourcesAPI
	DataSourceJobs DataSourceJobsAPI
}

type TenantsAPI interface {
	Find(context.Context) ([]*domain.Tenant, error)
	FindByID(context.Context, domain.TenantID) (*domain.Tenant, error)
	Create(context.Context, *commands.CreateTenantCommand) (*domain.Tenant, error)
	Update(context.Context, *commands.UpdateTenantCommand) (*domain.Tenant, error)
	Delete(context.Context, domain.TenantID) error
}

type DataSourcesAPI interface {
	Find(context.Context, *ListDataSourceQuery) ([]*domain.DataSource, error)
	FindByID(context.Context, *FindDataSourceByIDQuery) (*domain.DataSource, error)
	Create(context.Context, *commands.CreateDataSourceCommand) (*domain.DataSource, error)
	Update(context.Context, *commands.UpdateDataSourceCommand) (*domain.DataSource, error)
	Delete(context.Context, *commands.RemoveDataSourceCommand) error
	Lookup(context.Context, *GetDataSourceLookupsQuery) ([]*domain.Lookup, error)
}

type DataSourceJobsAPI interface {
	Find(context.Context, *FindDataSourceJobsQuery) ([]*domain.DataSource, error)
	Process(context.Context, *commands.ProcessDataSourceJobCommand) error
}
