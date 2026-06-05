package services

import (
	"bytes"
	"context"
	"log/slog"

	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

var (
	buf    bytes.Buffer
	logger = slog.New(slog.NewTextHandler(&buf, nil))
)

type mockDatabase struct {
	closeFn      func() error
	beginTxFn    func(context.Context) (context.Context, error)
	commitTxFn   func(context.Context) error
	rollbackTxfn func(context.Context) error
}

func (d *mockDatabase) Close(context.Context) error {
	return d.closeFn()
}

func (d *mockDatabase) BeginTx(ctx context.Context) (context.Context, error) {
	return d.beginTxFn(ctx)
}

func (d *mockDatabase) CommitTx(ctx context.Context) error {
	return d.commitTxFn(ctx)
}

func (d *mockDatabase) RollbackTx(ctx context.Context) error {
	return d.rollbackTxfn(ctx)
}

type mockTenantsRepository struct {
	findFn     func(context.Context) ([]*domain.Tenant, error)
	findByIdFn func(context.Context, domain.TenantID) (*domain.Tenant, error)
	existsFn   func(context.Context, domain.TenantID) (bool, error)
	upsertFn   func(context.Context, *domain.Tenant) (*domain.Tenant, error)
	deleteFn   func(context.Context, domain.TenantID) error
}

func (r *mockTenantsRepository) Find(ctx context.Context) ([]*domain.Tenant, error) {
	return r.findFn(ctx)
}

func (r *mockTenantsRepository) FindByID(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
	return r.findByIdFn(ctx, id)
}

func (r *mockTenantsRepository) Exists(ctx context.Context, id domain.TenantID) (bool, error) {
	return r.existsFn(ctx, id)
}

func (r *mockTenantsRepository) Upsert(ctx context.Context, t *domain.Tenant) (*domain.Tenant, error) {
	return r.upsertFn(ctx, t)
}

func (r *mockTenantsRepository) Delete(ctx context.Context, id domain.TenantID) error {
	return r.deleteFn(ctx, id)
}

type mockDataSourcesRepository struct {
	findFn      func(context.Context, domain.TenantID) ([]*domain.DataSource, error)
	findByIdFn  func(context.Context, domain.TenantID, domain.DataSourceID) (*domain.DataSource, error)
	findJobsFn  func(context.Context, *ports.FindDataSourceJobsFilter) ([]*domain.DataSource, error)
	existsFn    func(context.Context, domain.TenantID, domain.DataSourceID) (bool, error)
	upsertFn    func(context.Context, *domain.DataSource) (*domain.DataSource, error)
	deleteFn    func(context.Context, domain.TenantID, domain.DataSourceID) error
	deleteAllFn func(context.Context, domain.TenantID) error
}

func (r *mockDataSourcesRepository) Find(ctx context.Context, tenantID domain.TenantID) ([]*domain.DataSource, error) {
	return r.findFn(ctx, tenantID)
}

func (r *mockDataSourcesRepository) FindByID(ctx context.Context, tenantID domain.TenantID, id domain.DataSourceID) (*domain.DataSource, error) {
	return r.findByIdFn(ctx, tenantID, id)
}

func (r *mockDataSourcesRepository) FindJobs(ctx context.Context, filter *ports.FindDataSourceJobsFilter) ([]*domain.DataSource, error) {
	return r.findJobsFn(ctx, filter)
}

func (r *mockDataSourcesRepository) Exists(ctx context.Context, tenantID domain.TenantID, id domain.DataSourceID) (bool, error) {
	return r.existsFn(ctx, tenantID, id)
}

func (r *mockDataSourcesRepository) Upsert(ctx context.Context, ds *domain.DataSource) (*domain.DataSource, error) {
	return r.upsertFn(ctx, ds)
}

func (r *mockDataSourcesRepository) Delete(ctx context.Context, tenantID domain.TenantID, id domain.DataSourceID) error {
	return r.deleteFn(ctx, tenantID, id)
}

func (r *mockDataSourcesRepository) DeleteAll(ctx context.Context, tenantID domain.TenantID) error {
	return r.deleteAllFn(ctx, tenantID)
}

type mockLookupClient struct {
	fetchLookupsFn func(context.Context, domain.DataSourceHTTPRequest, map[string]any) ([]map[string]any, error)
}

func (c *mockLookupClient) FetchLookups(ctx context.Context, request domain.DataSourceHTTPRequest, params map[string]any) ([]map[string]any, error) {
	return c.fetchLookupsFn(ctx, request, params)
}

type mockLookupStrategy struct {
	lookupFn func(context.Context, *domain.DataSource, map[string]any) ([]*domain.Lookup, error)
}

func (s *mockLookupStrategy) Lookup(ctx context.Context, ds *domain.DataSource, query map[string]any) ([]*domain.Lookup, error) {
	return s.lookupFn(ctx, ds, query)
}
