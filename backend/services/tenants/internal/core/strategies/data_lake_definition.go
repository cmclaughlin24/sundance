package strategies

import (
	"context"
	"fmt"
	"log/slog"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type DataLakeDataSourceDefinitionStrategy struct {
	baseDefinitionStrategy
	client ports.DataLakeClient
}

func NewDataLakeDataSourceDefinitionStrategy(logger *slog.Logger, client *ports.Clients) ports.DataSourceDefinitionStrategy {
	return &DataLakeDataSourceDefinitionStrategy{
		client: client.DataLake,
		baseDefinitionStrategy: baseDefinitionStrategy{
			logger: logger,
		},
	}
}

func (s *DataLakeDataSourceDefinitionStrategy) Description() string {
	return "Options resolved by running a parameterized query against the data lake, supporting catalog and schema selection, required and optional keys, and query timeouts."
}

func (s *DataLakeDataSourceDefinitionStrategy) Attributes() []ports.AttributeDescriptor {
	return reflectAttributes[domain.DataLakeDataSourceAttributes]()
}

func (s *DataLakeDataSourceDefinitionStrategy) Lookup(ctx context.Context, ds *domain.DataSource, params map[string]any) ([]*domain.Lookup, error) {
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
