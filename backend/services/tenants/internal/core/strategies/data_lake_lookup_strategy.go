package strategies

import (
	"context"
	"log/slog"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type DataLakeLookupStrategy struct {
	logger *slog.Logger
	client ports.DataLakeClient
}

func NewDataLakeLookupStrategy(logger *slog.Logger, client *ports.Clients) ports.LookupStrategy {
	return &DataLakeLookupStrategy{
		logger: logger,
		client: client.DataLake,
	}
}

func (s *DataLakeLookupStrategy) Lookup(ctx context.Context, ds *domain.DataSource) ([]*domain.Lookup, error) {
	attr, err := domain.GetDataSourceAttributes[domain.DataLakeDataSourceAttributes](ds.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "data_source_id", ds.ID, "data_source_type", ds.Type, "error", err)
		return nil, err
	}

	s.logger.DebugContext(ctx, "data lake lookup request", "data_source_id", ds.ID, "catalog", attr.Catalog, "schema", attr.Schema)

	lookups, err := s.client.Query(ctx, attr, nil)
	if err != nil {
		return nil, err
	}

	s.logger.DebugContext(ctx, "data lake lookup resolved", "data_source_id", ds.ID, "count", len(lookups))

	return lookups, nil
}
