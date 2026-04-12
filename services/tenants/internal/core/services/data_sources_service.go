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

func (s *DataSourcesService) Find(ctx context.Context) ([]*domain.DataSource, error) {
	return nil, nil
}

func (s *DataSourcesService) FindById(context.Context, domain.DataSourceID) (*domain.DataSource, error) {
	return nil, nil
}

func (s *DataSourcesService) Create(context.Context, ports.CreateDataSourceCommand) (*domain.DataSource, error) {
	return nil, nil
}

func (s *DataSourcesService) Update(context.Context, ports.UpdateDataSourceCommand) (*domain.DataSource, error) {
	return nil, nil
}

func (s *DataSourcesService) Remove(context.Context, domain.DataSourceID) error {
	return nil
}
