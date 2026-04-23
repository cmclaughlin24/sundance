package services

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type DataSourcesService struct {
	logger                *log.Logger
	tenantsRepository     ports.TenantsRepository
	dataSourcesRepository ports.DataSourcesRepository
}

func NewDataSourcesService(logger *log.Logger, repository *ports.Repository) *DataSourcesService {
	return &DataSourcesService{
		logger:                logger,
		dataSourcesRepository: repository.DataSources,
		tenantsRepository:     repository.Tenants,
	}
}

func (s *DataSourcesService) Find(ctx context.Context, query *ports.ListDataSourceQuery) ([]*domain.DataSource, error) {
	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	return s.dataSourcesRepository.Find(ctx, query.TenantID)
}

func (s *DataSourcesService) FindById(ctx context.Context, query *ports.FindDataSourceByIDQuery) (*domain.DataSource, error) {
	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	if err := s.tenantExists(ctx, query.TenantID); err != nil {
		return nil, err
	}

	return s.dataSourcesRepository.FindById(ctx, query.TenantID, query.ID)
}

func (s *DataSourcesService) Create(ctx context.Context, command *ports.CreateDataSourceCommand) (*domain.DataSource, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	if err := s.tenantExists(ctx, command.TenantID); err != nil {
		return nil, err
	}

	ds, err := domain.NewDataSource(
		"",
		command.TenantID,
		command.Name,
		command.Description,
		command.Type,
		command.Attributes,
	)

	if err != nil {
		return nil, err
	}

	dataSource, err := s.dataSourcesRepository.Upsert(ctx, ds)

	if err != nil {
		return nil, err
	}

	return dataSource, nil
}

func (s *DataSourcesService) Update(ctx context.Context, command *ports.UpdateDataSourceCommand) (*domain.DataSource, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	if err := s.tenantExists(ctx, command.TenantID); err != nil {
		return nil, err
	}

	ds, err := domain.NewDataSource(
		command.ID,
		command.TenantID,
		command.Name,
		command.Description,
		command.Type,
		command.Attributes,
	)

	if err != nil {
		return nil, err
	}

	dataSource, err := s.dataSourcesRepository.Upsert(ctx, ds)

	if err != nil {
		return nil, err
	}

	return dataSource, nil
}

func (s *DataSourcesService) Remove(ctx context.Context, command *ports.RemoveDataSourceCommand) error {
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

	return s.dataSourcesRepository.Remove(ctx, command.TenantID, command.ID)
}

func (s *DataSourcesService) Lookup(ctx context.Context, command *ports.GetDataSourceLookupsCommand) ([]*domain.DataSourceLookup, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	_, err := s.dataSourcesRepository.FindById(ctx, command.TenantID, command.ID)

	if err != nil {
		return nil, err
	}

	// TODO: Implement data source lookup strategy pattern based on the type of data source.

	return nil, nil
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
