package ports

import (
	"context"

	"sundance/backend/services/tenants/internal/core/domain"
)

type API struct {
	Tenants        TenantsAPI
	DataSources    DataSourcesAPI
	DataSourceJobs DataSourceJobsAPI
}

type TenantsAPI interface {
	Find(context.Context) ([]*domain.Tenant, error)
	FindByID(context.Context, domain.TenantID) (*domain.Tenant, error)
	Create(context.Context, *CreateTenantCommand) (*domain.Tenant, error)
	Update(context.Context, *UpdateTenantCommand) (*domain.Tenant, error)
	Delete(context.Context, domain.TenantID) error
}

type DataSourcesAPI interface {
	Find(context.Context, *ListDataSourceQuery) ([]*domain.DataSource, error)
	FindByID(context.Context, *FindDataSourceByIDQuery) (*domain.DataSource, error)
	Create(context.Context, *CreateDataSourceCommand) (*domain.DataSource, error)
	Update(context.Context, *UpdateDataSourceCommand) (*domain.DataSource, error)
	Delete(context.Context, *RemoveDataSourceCommand) error
	Lookup(context.Context, *GetDataSourceLookupsQuery) ([]*domain.Lookup, error)
}

type DataSourceJobsAPI interface {
	Find(context.Context, *FindDataSourceJobsQuery) ([]*domain.DataSource, error)
	Process(context.Context, *ProcessDataSourceJobCommand) error
}
