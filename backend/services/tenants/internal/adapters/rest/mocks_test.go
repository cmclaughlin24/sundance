package rest

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type mockTenantsService struct {
	findFn     func(context.Context) ([]*domain.Tenant, error)
	findByIdFn func(context.Context, domain.TenantID) (*domain.Tenant, error)
	createFn   func(context.Context, *ports.CreateTenantCommand) (*domain.Tenant, error)
	updateFn   func(context.Context, *ports.UpdateTenantCommand) (*domain.Tenant, error)
	deleteFn   func(context.Context, domain.TenantID) error
}

func (s *mockTenantsService) Find(ctx context.Context) ([]*domain.Tenant, error) {
	return s.findFn(ctx)
}

func (s *mockTenantsService) FindByID(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
	return s.findByIdFn(ctx, id)
}

func (s *mockTenantsService) Create(ctx context.Context, command *ports.CreateTenantCommand) (*domain.Tenant, error) {
	return s.createFn(ctx, command)
}

func (s *mockTenantsService) Update(ctx context.Context, command *ports.UpdateTenantCommand) (*domain.Tenant, error) {
	return s.updateFn(ctx, command)
}

func (s *mockTenantsService) Delete(ctx context.Context, id domain.TenantID) error {
	return s.deleteFn(ctx, id)
}

type mockDataSourcesService struct {
	findFn     func(context.Context, *ports.ListDataSourceQuery) ([]*domain.DataSource, error)
	findByIdFn func(context.Context, *ports.FindDataSourceByIDQuery) (*domain.DataSource, error)
	createFn   func(context.Context, *ports.CreateDataSourceCommand) (*domain.DataSource, error)
	updateFn   func(context.Context, *ports.UpdateDataSourceCommand) (*domain.DataSource, error)
	deleteFn   func(context.Context, *ports.RemoveDataSourceCommand) error
	lookupFn   func(context.Context, *ports.GetDataSourceLookupsQuery) ([]*domain.Lookup, error)
}

func (s *mockDataSourcesService) Find(ctx context.Context, query *ports.ListDataSourceQuery) ([]*domain.DataSource, error) {
	return s.findFn(ctx, query)
}

func (s *mockDataSourcesService) FindByID(ctx context.Context, query *ports.FindDataSourceByIDQuery) (*domain.DataSource, error) {
	return s.findByIdFn(ctx, query)
}

func (s *mockDataSourcesService) Create(ctx context.Context, command *ports.CreateDataSourceCommand) (*domain.DataSource, error) {
	return s.createFn(ctx, command)
}

func (s *mockDataSourcesService) Update(ctx context.Context, command *ports.UpdateDataSourceCommand) (*domain.DataSource, error) {
	return s.updateFn(ctx, command)
}

func (s *mockDataSourcesService) Delete(ctx context.Context, command *ports.RemoveDataSourceCommand) error {
	return s.deleteFn(ctx, command)
}

func (s *mockDataSourcesService) Lookup(ctx context.Context, query *ports.GetDataSourceLookupsQuery) ([]*domain.Lookup, error) {
	return s.lookupFn(ctx, query)
}
