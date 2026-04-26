package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
)

type Services struct {
	Tenants     TenantsService
	DataSources DataSourcesService
}

type TenantsService interface {
	Find(context.Context) ([]*domain.Tenant, error)
	FindByID(context.Context, domain.TenantID) (*domain.Tenant, error)
	Create(context.Context, *CreateTenantCommand) (*domain.Tenant, error)
	Update(context.Context, *UpdateTenantCommand) (*domain.Tenant, error)
	Delete(context.Context, domain.TenantID) error
}

type DataSourcesService interface {
	Find(context.Context, *ListDataSourceQuery) ([]*domain.DataSource, error)
	FindByID(context.Context, *FindDataSourceByIDQuery) (*domain.DataSource, error)
	Create(context.Context, *CreateDataSourceCommand) (*domain.DataSource, error)
	Update(context.Context, *UpdateDataSourceCommand) (*domain.DataSource, error)
	Delete(context.Context, *RemoveDataSourceCommand) error
	Lookup(context.Context, *GetDataSourceLookupsQuery) ([]*domain.Lookup, error)
}
