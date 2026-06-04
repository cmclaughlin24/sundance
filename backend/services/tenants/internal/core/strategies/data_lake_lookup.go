package strategies

import (
	"context"
	"fmt"
	"log/slog"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type DataLakeLookupStrategy struct {
	baseLookupStrategy
	client ports.DataLakeClient
}

func NewDataLakeLookupStrategy(logger *slog.Logger, client *ports.Clients) ports.LookupStrategy {
	return &DataLakeLookupStrategy{
		client: client.DataLake,
		baseLookupStrategy: baseLookupStrategy{
			logger: logger,
		},
	}
}

func (s *DataLakeLookupStrategy) Lookup(ctx context.Context, ds *domain.DataSource, params map[string]any) ([]*domain.Lookup, error) {
	attr, err := domain.GetDataSourceAttributes[domain.DataLakeDataSourceAttributes](ds.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "data_source_id", ds.ID, "data_source_type", ds.Type, "error", err)
		return nil, err
	}

	if missing := s.missingRequiredKeys(attr.RequiredKeys, params); len(missing) > 0 {
		s.logger.ErrorContext(ctx, "data lake lookup missing required keys", "data_source_id", ds.ID, "missing_keys", missing)
		return nil, fmt.Errorf("%w: %v", domain.ErrMissingRequiredKeys, missing)
	}

	s.logger.DebugContext(ctx, "data lake lookup request", "data_source_id", ds.ID, "catalog", attr.Catalog, "schema", attr.Schema)

	lookups, err := s.client.Query(ctx, attr, params)
	if err != nil {
		return nil, err
	}

	s.logger.DebugContext(ctx, "data lake lookup resolved", "data_source_id", ds.ID, "count", len(lookups))

	return lookups, nil
}
