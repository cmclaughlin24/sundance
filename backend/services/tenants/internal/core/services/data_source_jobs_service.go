package services

import (
	"context"
	"log/slog"

	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type dataSourcesJobService struct {
	logger     *slog.Logger
	repository ports.DataSourcesRepository
	client     ports.LookupClient
}

func NewDataSourcesJobService(logger *slog.Logger, repository *ports.Repository, clients *ports.Clients) ports.DataSourceJobsAPI {
	return &dataSourcesJobService{
		logger:     logger,
		repository: repository.DataSources,
		client:     clients.Lookups,
	}
}

func (s *dataSourcesJobService) Find(ctx context.Context, query *ports.FindDataSourceJobsQuery) ([]*domain.DataSource, error) {
	s.logger.DebugContext(ctx, "listing data source jobs")

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "data source job listing failed; invalid query", "error", err)
		return nil, err
	}

	sources, err := s.repository.FindJobs(ctx, &ports.FindDataSourceJobsFilter{
		Types:             []domain.DataSourceType{domain.DataSourceTypeScheduled},
		Take:              query.Take,
		ExpiredAtOrBefore: Now(),
		RetryLimit:        query.RetryLimit,
	})

	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve data source jobs", "error", err)
		return nil, err
	}

	return sources, nil
}

func (s *dataSourcesJobService) Process(ctx context.Context, command *ports.ProcessDataSourceJobCommand) error {
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

	data, err := s.client.FetchLookups(ctx, attr.DataSourceHTTPRequest, nil)
	if err != nil {
		attr.RecordAttempt()
		s.logger.ErrorContext(ctx, "failed to fetch lookups for data source", "data_source_id", ds.ID, "error", err, "attempts", attr.Attempts)
	} else {
		lookups := s.toLookups(ctx, data, attr.ValueField, attr.LabelField, ds.ID)
		attr.RefreshData(lookups)
	}

	ds.UpdateAttributes(attr)

	if _, err := s.repository.Upsert(ctx, ds); err != nil {
		return err
	}

	return nil
}

func (s *dataSourcesJobService) toLookups(ctx context.Context, rows []map[string]any, valueField, labelField string, dataSourceID domain.DataSourceID) []*domain.Lookup {
	lookups := make([]*domain.Lookup, 0, len(rows))

	for i, row := range rows {
		value, ok := row[valueField]
		if !ok {
			s.logger.WarnContext(ctx, "skipping lookup row missing value field", "data_source_id", dataSourceID, "row_index", i, "value_field", valueField)
			continue
		}

		label, ok := row[labelField]
		if !ok {
			s.logger.WarnContext(ctx, "skipping lookup row missing label field", "data_source_id", dataSourceID, "row_index", i, "label_field", labelField)
			continue
		}

		switch l := label.(type) {
		case string:
			lookups = append(lookups, domain.NewLookup(value, l))
		default:
			s.logger.WarnContext(ctx, "skipping lookup row label field is not string", "data_source_id", dataSourceID, "row_index", i, "label_field", labelField)
		}
	}
	return lookups
}
