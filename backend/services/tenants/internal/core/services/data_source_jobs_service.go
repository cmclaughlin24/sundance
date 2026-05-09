package services

import (
	"context"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type DataSourcesJobService struct {
	logger     *slog.Logger
	repository ports.DataSourcesRepository
}

func NewDataSourcesJobService(logger *slog.Logger, repository *ports.Repository) ports.DataSourceJobsService {
	return &DataSourcesJobService{
		logger:     logger,
		repository: repository.DataSources,
	}
}

func (s *DataSourcesJobService) Find(ctx context.Context) ([]*domain.DataSource, error) {
	s.logger.DebugContext(ctx, "listing data source jobs")

	return nil, nil
}

func (s *DataSourcesJobService) Process(ctx context.Context, command *ports.ProcessDataSourceJobCommand) error {
	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "data source job process failed; invalid command", "error", err)
		return err
	}

	ds := command.DataSource
	s.logger.DebugContext(ctx, "processing data source", "data_source_id", ds.ID)

	return nil
}
