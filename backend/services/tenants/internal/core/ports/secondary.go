package ports

import (
	"context"

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
	FindById(context.Context, domain.TenantID) (*domain.Tenant, error)
	Exists(context.Context, domain.TenantID) (bool, error)
	Upsert(context.Context, *domain.Tenant) (*domain.Tenant, error)
	Remove(context.Context, domain.TenantID) error
}

type DataSourcesRepository interface {
	Find(context.Context, domain.TenantID) ([]*domain.DataSource, error)
	FindById(context.Context, domain.TenantID, domain.DataSourceID) (*domain.DataSource, error)
	Exists(context.Context, domain.TenantID, domain.DataSourceID) (bool, error)
	Upsert(context.Context, *domain.DataSource) (*domain.DataSource, error)
	Remove(context.Context, domain.TenantID, domain.DataSourceID) error
}
