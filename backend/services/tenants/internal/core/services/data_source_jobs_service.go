package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type DataSourcesJobService struct {
	logger     *slog.Logger
	repository ports.DataSourcesRepository
	client     ports.LookupClient
}

func NewDataSourcesJobService(logger *slog.Logger, repository *ports.Repository, clients *ports.Clients) ports.DataSourceJobsService {
	return &DataSourcesJobService{
		logger:     logger,
		repository: repository.DataSources,
		client:     clients.Lookups,
	}
}

func (s *DataSourcesJobService) Find(ctx context.Context) ([]*domain.DataSource, error) {
	s.logger.DebugContext(ctx, "listing data source jobs")

	sources, err := s.repository.FindJobs(ctx, &ports.FindDataSourceJobsFilter{})
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve data source jobs", "error", err)
		return nil, err
	}

	return sources, nil
}

func (s *DataSourcesJobService) Process(ctx context.Context, command *ports.ProcessDataSourceJobCommand) error {
	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "data source job process failed; invalid command", "error", err)
		return err
	}

	ds := command.DataSource
	s.logger.DebugContext(ctx, "processing data source job", "data_source_id", ds.ID)

	attr, err := domain.GetDataSourceAttributes[domain.ScheduledDataSourceAttributes](ds.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to process data source job", "data_source_id", ds.ID, "type", ds.Type, "error", err)
		return err
	}

	lookups, err := s.client.FetchLookups(ctx, attr.Method, attr.URL, attr.Headers)
	if err != nil {
		return err
	}

	// FIXME: Refactor into a domain update method.
	attr.Data = lookups
	attr.ExpirationDate = time.Now().Add(time.Duration(attr.IntervalHours * float64(time.Hour)))
	ds.Attributes = attr

	if _, err := s.repository.Upsert(ctx, ds); err != nil {
		return err
	}

	return nil
}
