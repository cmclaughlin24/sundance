package services

import (
	"context"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type DataSourcesService struct {
	logger                *slog.Logger
	tenantsRepository     ports.TenantsRepository
	dataSourcesRepository ports.DataSourcesRepository
	lookupStrategies      ports.LookupStrategyRegistry
}

func NewDataSourcesService(logger *slog.Logger, repository *ports.Repository, strategies *ports.Strategies) ports.DataSourcesService {
	return &DataSourcesService{
		logger:                logger,
		dataSourcesRepository: repository.DataSources,
		tenantsRepository:     repository.Tenants,
		lookupStrategies:      strategies.Lookups,
	}
}

func (s *DataSourcesService) Find(ctx context.Context, query *ports.ListDataSourceQuery) ([]*domain.DataSource, error) {
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

func (s *DataSourcesService) FindByID(ctx context.Context, query *ports.FindDataSourceByIDQuery) (*domain.DataSource, error) {
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

func (s *DataSourcesService) Create(ctx context.Context, command *ports.CreateDataSourceCommand) (*domain.DataSource, error) {
	s.logger.DebugContext(ctx, "creating data source", "tenant_id", command.TenantID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "data source creation failed; invalid command", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	if err := s.tenantExists(ctx, command.TenantID); err != nil {
		return nil, err
	}

	ds, err := domain.NewDataSource(
		command.TenantID,
		command.Name,
		command.Description,
		command.Type,
		command.Attributes,
	)

	if err != nil {
		s.logger.WarnContext(ctx, "data source creation failed; domain invariant violation", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	ds, err = s.dataSourcesRepository.Upsert(ctx, ds)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist data source", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "data source created", "tenant_id", command.TenantID, "data_source_id", ds.ID)

	return ds, nil
}

func (s *DataSourcesService) Update(ctx context.Context, command *ports.UpdateDataSourceCommand) (*domain.DataSource, error) {
	s.logger.DebugContext(ctx, "updating data source", "tenant_id", command.TenantID, "data_source_id", command.ID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "data source update failed; invalid command", "tenant_id", command.TenantID, "data_source_id", command.ID, "error", err)
		return nil, err
	}

	if err := s.tenantExists(ctx, command.TenantID); err != nil {
		return nil, err
	}

	ds, err := s.dataSourcesRepository.FindByID(ctx, command.TenantID, command.ID)
	if err != nil {
		s.logFindByIDError(ctx, err, command.TenantID, command.ID)
		return nil, err
	}

	if err := ds.Update(command.Name, command.Description, command.Type, command.Attributes); err != nil {
		s.logger.WarnContext(ctx, "data source update failed; domain invariant violation", "tenant_id", command.TenantID, "data_source_id", command.ID, "error", err)
		return nil, err
	}

	ds, err = s.dataSourcesRepository.Upsert(ctx, ds)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist data source", "tenant_id", command.TenantID, "data_source_id", command.ID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "data source updated", "tenant_id", command.TenantID, "data_source_id", command.ID)

	return ds, nil
}

func (s *DataSourcesService) Delete(ctx context.Context, command *ports.RemoveDataSourceCommand) error {
	s.logger.DebugContext(ctx, "deleting data source", "tenant_id", command.TenantID, "data_source_id", command.ID)

	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "data source deletion failed; invalid command", "tenant_id", command.TenantID, "data_source_id", command.ID, "error", err)
		return err
	}

	if err := s.tenantExists(ctx, command.TenantID); err != nil {
		return err
	}

	exists, err := s.dataSourcesRepository.Exists(ctx, command.TenantID, command.ID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to check data source existence", "tenant_id", command.TenantID, "data_source_id", command.ID, "error", err)
		return err
	}

	if !exists {
		s.logger.WarnContext(ctx, "data source not found", "tenant_id", command.TenantID, "data_source_id", command.ID)
		return common.ErrNotFound
	}

	if err := s.dataSourcesRepository.Delete(ctx, command.TenantID, command.ID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete data source", "tenant_id", command.TenantID, "data_source_id", command.ID, "error", err)
		return err
	}

	s.logger.InfoContext(ctx, "data source deleted", "tenant_id", command.TenantID, "data_source_id", command.ID)

	return nil
}

func (s *DataSourcesService) Lookup(ctx context.Context, query *ports.GetDataSourceLookupsQuery) ([]*domain.Lookup, error) {
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

	lookups, err := strategy.Lookup(ctx, ds)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to execute lookup", "tenant_id", query.TenantID, "data_source_id", query.ID, "type", ds.Type, "error", err)
		return nil, err
	}

	return lookups, nil
}

func (s *DataSourcesService) tenantExists(ctx context.Context, tenantID domain.TenantID) error {
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

func (s *DataSourcesService) logFindByIDError(ctx context.Context, err error, tenantID domain.TenantID, id domain.DataSourceID) {
	if err == common.ErrNotFound {
		s.logger.WarnContext(ctx, "data source not found", "tenant_id", tenantID, "data_source_id", id)
	} else {
		s.logger.ErrorContext(ctx, "failed to retrieve data source", "tenant_id", tenantID, "data_source_id", id, "error", err)
	}
}
