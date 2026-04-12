package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
)

type Repository struct {
	Database    Database
	Tenants     TenantsRepository
	DataSources DataSourcesRepository
}

type Database interface {
	Close() error
	BeginTx(context.Context) (context.Context, error)
	GetTx(context.Context) (any, error)
	CommitTx(context.Context) error
	RollbackTx(context.Context) error
}

type TenantsRepository interface {
	Exists(context.Context, domain.TenantID) (bool, error)
	Find(context.Context) ([]*domain.Tenant, error)
	FindById(context.Context, domain.TenantID) (*domain.Tenant, error)
	Upsert(context.Context, *domain.Tenant) (*domain.Tenant, error)
	Remove(context.Context, domain.TenantID) error
}

type DataSourcesRepository interface {
	Exists(context.Context, domain.DataSourceID) (bool, error)
	Find(context.Context) ([]*domain.DataSource, error)
	FindById(context.Context, domain.DataSourceID) (*domain.DataSource, error)
	Upsert(context.Context, *domain.DataSource) (*domain.DataSource, error)
	Remove(context.Context, domain.DataSourceID) error
}
