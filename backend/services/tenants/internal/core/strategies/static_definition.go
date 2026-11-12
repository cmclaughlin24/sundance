package strategies

import (
	"context"
	"log/slog"

	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type StaticDataSourceDefinitionStrategy struct {
	logger *slog.Logger
}

func NewStaticDataSourceDefinitionStrategy(logger *slog.Logger) ports.DataSourceDefinitionStrategy {
	return &StaticDataSourceDefinitionStrategy{
		logger: logger,
	}
}

func (s *StaticDataSourceDefinitionStrategy) Description() string {
	return "Fixed list of value and label options defined directly on the data source."
}

func (s *StaticDataSourceDefinitionStrategy) Attributes() []ports.AttributeDescriptor {
	return reflectAttributes[domain.StaticDataSourceAttributes]()
}

func (s *StaticDataSourceDefinitionStrategy) Lookup(ctx context.Context, ds *domain.DataSource, _ map[string]any) ([]*domain.Lookup, error) {
	attr, err := domain.GetDataSourceAttributes[domain.StaticDataSourceAttributes](ds.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "data_source_id", ds.ID, "data_source_type", ds.Type, "error", err)
		return nil, err
	}

	return attr.Data, nil
}
