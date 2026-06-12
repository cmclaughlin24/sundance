package services

import (
	"context"
	"errors"
	"log/slog"

	"sundance/backend/pkg/common"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type dataSourcesService struct {
	logger                *slog.Logger
	tenantsRepository     ports.TenantsRepository
	dataSourcesRepository ports.DataSourcesRepository
	lookupStrategies      ports.LookupStrategyRegistry
}

func NewDataSourcesService(logger *slog.Logger, repository *ports.Repository, strategies *ports.Strategies) ports.DataSourcesAPI {
	return &dataSourcesService{
		logger:                logger,
		dataSourcesRepository: repository.DataSources,
		tenantsRepository:     repository.Tenants,
		lookupStrategies:      strategies.Lookups,
	}
}

func (s *dataSourcesService) Find(ctx context.Context, query *ports.ListDataSourceQuery) ([]*domain.DataSource, error) {
	s.logger.DebugContext(ctx, "listing data sources", "tenant_id", query.TenantID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "data source listing failed; invalid query", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	sources, err := s.dataSourcesRepository.Find(ctx, query.TenantID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve data sources", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	return sources, nil
}

func (s *dataSourcesService) FindByID(ctx context.Context, query *ports.FindDataSourceByIDQuery) (*domain.DataSource, error) {
	s.logger.DebugContext(ctx, "finding data source", "tenant_id", query.TenantID, "data_source_id", query.ID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "data source find failed; invalid query", "tenant_id", query.TenantID, "data_source_id", query.ID, "error", err)
		return nil, err
	}

	if err := s.tenantExists(ctx, query.TenantID); err != nil {
		return nil, err
	}

	ds, err := s.dataSourcesRepository.FindByID(ctx, query.TenantID, query.ID)
	if err != nil {
		s.logFindByIDError(ctx, err, query.TenantID, query.ID)
		return nil, err
	}

	return ds, nil
}

func (s *dataSourcesService) Create(ctx context.Context, cmd *ports.CreateDataSourceCommand) (*domain.DataSource, error) {
	s.logger.DebugContext(ctx, "creating data source", "tenant_id", cmd.TenantID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "data source creation failed; invalid command", "tenant_id", cmd.TenantID, "error", err)
		return nil, err
	}

	if err := s.tenantExists(ctx, cmd.TenantID); err != nil {
		return nil, err
	}

	ds, err := domain.NewDataSource(
		cmd.TenantID,
		cmd.Name,
		cmd.Description,
		cmd.Type,
		cmd.Attributes,
	)

	if err != nil {
		s.logger.WarnContext(ctx, "data source creation failed; domain invariant violation", "tenant_id", cmd.TenantID, "error", err)
		return nil, err
	}

	ds, err = s.dataSourcesRepository.Upsert(ctx, ds)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist data source", "tenant_id", cmd.TenantID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "data source created", "tenant_id", cmd.TenantID, "data_source_id", ds.ID)

	return ds, nil
}

func (s *dataSourcesService) Update(ctx context.Context, cmd *ports.UpdateDataSourceCommand) (*domain.DataSource, error) {
	s.logger.DebugContext(ctx, "updating data source", "tenant_id", cmd.TenantID, "data_source_id", cmd.ID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "data source update failed; invalid command", "tenant_id", cmd.TenantID, "data_source_id", cmd.ID, "error", err)
		return nil, err
	}

	if err := s.tenantExists(ctx, cmd.TenantID); err != nil {
		return nil, err
	}

	ds, err := s.dataSourcesRepository.FindByID(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		s.logFindByIDError(ctx, err, cmd.TenantID, cmd.ID)
		return nil, err
	}

	if err := ds.Update(cmd.Name, cmd.Description, cmd.Type, cmd.Attributes); err != nil {
		s.logger.WarnContext(ctx, "data source update failed; domain invariant violation", "tenant_id", cmd.TenantID, "data_source_id", cmd.ID, "error", err)
		return nil, err
	}

	ds, err = s.dataSourcesRepository.Upsert(ctx, ds)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist data source", "tenant_id", cmd.TenantID, "data_source_id", cmd.ID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "data source updated", "tenant_id", cmd.TenantID, "data_source_id", cmd.ID)

	return ds, nil
}

func (s *dataSourcesService) Delete(ctx context.Context, cmd *ports.RemoveDataSourceCommand) error {
	s.logger.DebugContext(ctx, "deleting data source", "tenant_id", cmd.TenantID, "data_source_id", cmd.ID)

	if err := cmd.Validate(); err != nil {
		s.logger.WarnContext(ctx, "data source deletion failed; invalid command", "tenant_id", cmd.TenantID, "data_source_id", cmd.ID, "error", err)
		return err
	}

	if err := s.tenantExists(ctx, cmd.TenantID); err != nil {
		return err
	}

	exists, err := s.dataSourcesRepository.Exists(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to check data source existence", "tenant_id", cmd.TenantID, "data_source_id", cmd.ID, "error", err)
		return err
	}

	if !exists {
		s.logger.WarnContext(ctx, "data source not found", "tenant_id", cmd.TenantID, "data_source_id", cmd.ID)
		return common.ErrNotFound
	}

	if err := s.dataSourcesRepository.Delete(ctx, cmd.TenantID, cmd.ID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete data source", "tenant_id", cmd.TenantID, "data_source_id", cmd.ID, "error", err)
		return err
	}

	s.logger.InfoContext(ctx, "data source deleted", "tenant_id", cmd.TenantID, "data_source_id", cmd.ID)

	return nil
}

func (s *dataSourcesService) Lookup(ctx context.Context, query *ports.GetDataSourceLookupsQuery) ([]*domain.Lookup, error) {
	s.logger.DebugContext(ctx, "looking up data source", "tenant_id", query.TenantID, "data_source_id", query.ID)

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "data source lookup failed; invalid query", "tenant_id", query.TenantID, "data_source_id", query.ID, "error", err)
		return nil, err
	}

	ds, err := s.dataSourcesRepository.FindByID(ctx, query.TenantID, query.ID)
	if err != nil {
		s.logFindByIDError(ctx, err, query.TenantID, query.ID)
		return nil, err
	}

	strategy, err := s.lookupStrategies.Get(ds.Type)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get lookup strategy", "tenant_id", query.TenantID, "data_source_id", query.ID, "type", ds.Type, "error", err)
		return nil, err
	}

	lookups, err := strategy.Lookup(ctx, ds, query.Params)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to execute lookup", "tenant_id", query.TenantID, "data_source_id", query.ID, "type", ds.Type, "error", err)
		return nil, err
	}

	return lookups, nil
}

func (s *dataSourcesService) tenantExists(ctx context.Context, tenantID domain.TenantID) error {
	exists, err := s.tenantsRepository.Exists(ctx, tenantID)

	if err != nil {
		s.logger.ErrorContext(ctx, "failed to check tenant existence", "tenant_id", tenantID, "error", err)
		return err
	}

	if !exists {
		s.logger.WarnContext(ctx, "tenant not found", "tenant_id", tenantID)
		return common.ErrNotFound
	}

	return nil
}

func (s *dataSourcesService) logFindByIDError(ctx context.Context, err error, tenantID domain.TenantID, id domain.DataSourceID) {
	if errors.Is(err, common.ErrNotFound) {
		s.logger.WarnContext(ctx, "data source not found", "tenant_id", tenantID, "data_source_id", id)
	} else {
		s.logger.ErrorContext(ctx, "failed to retrieve data source", "tenant_id", tenantID, "data_source_id", id, "error", err)
	}
}
