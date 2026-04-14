package services

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/tenants/internal/core/ports"
)

type DataSourcesService struct {
	logger     *log.Logger
	repository *ports.Repository
}

func NewDataSourcesService(logger *log.Logger, repository *ports.Repository) *DataSourcesService {
	return &DataSourcesService{
		logger:     logger,
		repository: repository,
	}
}

func (s *DataSourcesService) Find(ctx context.Context, query *ports.ListDataSourceQuery) ([]*domain.DataSource, error) {
	return s.repository.DataSources.Find(ctx, query.TenantID)
}

func (s *DataSourcesService) FindById(ctx context.Context, tenantId domain.TenantID, sourceId domain.DataSourceID) (*domain.DataSource, error) {
	return s.repository.DataSources.FindById(ctx, tenantId, sourceId)
}

func (s *DataSourcesService) Create(ctx context.Context, command *ports.CreateDataSourceCommand) (*domain.DataSource, error) {
	dataSource, err := s.repository.DataSources.Upsert(
		ctx,
		domain.NewDataSource("", command.TenantID, command.Type, command.Attributes),
	)

	if err != nil {
		return nil, err
	}

	return dataSource, nil
}

func (s *DataSourcesService) Update(ctx context.Context, command *ports.UpdateDataSourceCommand) (*domain.DataSource, error) {
	dataSource, err := s.repository.DataSources.Upsert(
		ctx,
		domain.NewDataSource(command.ID, command.TenantID, command.Type, command.Attributes),
	)

	if err != nil {
		return nil, err
	}

	return dataSource, nil
}

func (s *DataSourcesService) Remove(ctx context.Context, tenantId domain.TenantID, sourceId domain.DataSourceID) error {
	return s.repository.DataSources.Remove(ctx, tenantId, sourceId)
}

func (s *DataSourcesService) Lookup(ctx context.Context, tenantId domain.TenantID, sourceId domain.DataSourceID) ([]*domain.DataSourceLookup, error) {
	_, err := s.repository.DataSources.FindById(ctx, tenantId, sourceId)

	if err != nil {
		return nil, err
	}

	// TODO: Implement data source lookup strategy pattern based on the type of data source.

	return nil, nil
}
