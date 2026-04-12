package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
)

type Services struct {
	Tenants     TenantsService
	DataSources DataSourcesService
}

type TenantsService interface {
	Find(context.Context) ([]*domain.Tenant, error)
	FindById(context.Context, domain.TenantID) (*domain.Tenant, error)
	Create(context.Context, CreateTenantCommand) (*domain.Tenant, error)
	Update(context.Context, UpdateTenantCommand) (*domain.Tenant, error)
	Remove(context.Context, domain.TenantID) error
}

type DataSourcesService interface {
	Find(context.Context, ListDataSourceQuery) ([]*domain.DataSource, error)
	FindById(context.Context, domain.TenantID, domain.DataSourceID) (*domain.DataSource, error)
	Create(context.Context, CreateDataSourceCommand) (*domain.DataSource, error)
	Update(context.Context, UpdateDataSourceCommand) (*domain.DataSource, error)
	Remove(context.Context, domain.TenantID, domain.DataSourceID) error
}
