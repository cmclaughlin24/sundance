package strategies

import (
	"context"
	"log/slog"

	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type ScheduledLookupStrategy struct {
	logger *slog.Logger
}

func NewScheduledLookupStrategy(logger *slog.Logger) ports.LookupStrategy {
	return &ScheduledLookupStrategy{
		logger: logger,
	}
}

func (s *ScheduledLookupStrategy) Lookup(ctx context.Context, ds *domain.DataSource) ([]*domain.Lookup, error) {
	attr, err := domain.GetDataSourceAttributes[domain.ScheduledDataSourceAttributes](ds.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "data_source_id", ds.ID, "data_source_type", ds.Type, "error", err)
		return nil, err
	}

	// TODO: Determine if making data lazy-loaded would make sense.
	return attr.Data, nil
}
