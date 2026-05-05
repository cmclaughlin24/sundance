package services

import (
	"context"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
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
	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	return s.dataSourcesRepository.Find(ctx, query.TenantID)
}

func (s *DataSourcesService) FindByID(ctx context.Context, query *ports.FindDataSourceByIDQuery) (*domain.DataSource, error) {
	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	if err := s.tenantExists(ctx, query.TenantID); err != nil {
		return nil, err
	}

	return s.dataSourcesRepository.FindByID(ctx, query.TenantID, query.ID)
}

func (s *DataSourcesService) Create(ctx context.Context, command *ports.CreateDataSourceCommand) (*domain.DataSource, error) {
	if err := validate.ValidateStruct(command); err != nil {
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
		return nil, err
	}

	return s.dataSourcesRepository.Upsert(ctx, ds)
}

func (s *DataSourcesService) Update(ctx context.Context, command *ports.UpdateDataSourceCommand) (*domain.DataSource, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	if err := s.tenantExists(ctx, command.TenantID); err != nil {
		return nil, err
	}

	ds, err := s.dataSourcesRepository.FindByID(ctx, command.TenantID, command.ID)
	if err != nil {
		return nil, err
	}

	if err := ds.Update(command.Name, command.Description, command.Type, command.Attributes); err != nil {
		return nil, err
	}

	return s.dataSourcesRepository.Upsert(ctx, ds)
}

func (s *DataSourcesService) Delete(ctx context.Context, command *ports.RemoveDataSourceCommand) error {
	if err := validate.ValidateStruct(command); err != nil {
		return err
	}

	if err := s.tenantExists(ctx, command.TenantID); err != nil {
		return err
	}

	exists, err := s.dataSourcesRepository.Exists(ctx, command.TenantID, command.ID)
	if err != nil {
		return err
	}

	if !exists {
		return common.ErrNotFound
	}

	return s.dataSourcesRepository.Delete(ctx, command.TenantID, command.ID)
}

func (s *DataSourcesService) Lookup(ctx context.Context, command *ports.GetDataSourceLookupsQuery) ([]*domain.Lookup, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	ds, err := s.dataSourcesRepository.FindByID(ctx, command.TenantID, command.ID)
	if err != nil {
		return nil, err
	}

	strategy, err := s.lookupStrategies.Get(ds.Type)
	if err != nil {
		return nil, err
	}

	return strategy.Lookup(ctx, ds)
}

func (s *DataSourcesService) tenantExists(ctx context.Context, tenantID domain.TenantID) error {
	exists, err := s.tenantsRepository.Exists(ctx, tenantID)

	if err != nil {
		return err
	}

	if !exists {
		return common.ErrNotFound
	}

	return nil
}
