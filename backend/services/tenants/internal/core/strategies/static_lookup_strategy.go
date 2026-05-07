package strategies

import (
	"context"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type StaticLookupStrategy struct {
	logger *slog.Logger
}

func NewStaticLookupStrategy(logger *slog.Logger) ports.LookupStrategy {
	return &StaticLookupStrategy{
		logger: logger,
	}
}

func (s *StaticLookupStrategy) Lookup(ctx context.Context, ds *domain.DataSource) ([]*domain.Lookup, error) {
	attr, err := getDataSourceAttributes[domain.StaticDataSourceAttributes](ds.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "data_source_id", ds.ID, "data_source_type", ds.Type, "error", err)
		return nil, err
	}

	return attr.Data, nil
}
