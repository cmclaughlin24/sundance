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
	dataSourcesRepository ports.DataSourcesRepository
}

func NewDataSourcesService(logger *log.Logger, repository *ports.Repository) *DataSourcesService {
	return &DataSourcesService{
		logger:                logger,
		dataSourcesRepository: repository.DataSources,
	}
}

func (s *DataSourcesService) Find(ctx context.Context, query *ports.ListDataSourceQuery) ([]*domain.DataSource, error) {
	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	return s.dataSourcesRepository.Find(ctx, query.TenantID)
}

func (s *DataSourcesService) FindById(ctx context.Context, tenantID domain.TenantID, sourceID domain.DataSourceID) (*domain.DataSource, error) {
	return s.dataSourcesRepository.FindById(ctx, tenantID, sourceID)
}

func (s *DataSourcesService) Create(ctx context.Context, command *ports.CreateDataSourceCommand) (*domain.DataSource, error) {
	if err := validate.ValidateStruct(command); err != nil {
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

func (s *DataSourcesService) Remove(ctx context.Context, tenantID domain.TenantID, sourceID domain.DataSourceID) error {
	exists, err := s.dataSourcesRepository.Exists(ctx, tenantID, sourceID)

	if err != nil {
		return err
	}

	if !exists {
		return common.ErrNotFound
	}

	return s.dataSourcesRepository.Remove(ctx, tenantID, sourceID)
}

func (s *DataSourcesService) Lookup(ctx context.Context, tenantID domain.TenantID, sourceID domain.DataSourceID) ([]*domain.DataSourceLookup, error) {
	_, err := s.dataSourcesRepository.FindById(ctx, tenantID, sourceID)

	if err != nil {
		return nil, err
	}

	// TODO: Implement data source lookup strategy pattern based on the type of data source.

	return nil, nil
}
