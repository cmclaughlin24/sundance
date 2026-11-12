package strategies

import (
	"context"
	"log/slog"

	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type ScheduledDataSourceDefinitionStrategy struct {
	logger *slog.Logger
}

func NewScheduledDataSourceDefinitionStrategy(logger *slog.Logger) ports.DataSourceDefinitionStrategy {
	return &ScheduledDataSourceDefinitionStrategy{
		logger: logger,
	}
}

func (s *ScheduledDataSourceDefinitionStrategy) Description() string {
	return "Options fetched from an HTTP endpoint and cached, supporting periodic refresh on a configurable hour interval."
}

func (s *ScheduledDataSourceDefinitionStrategy) Attributes() []ports.AttributeDescriptor {
	return reflectAttributes[domain.ScheduledDataSourceAttributes]()
}

func (s *ScheduledDataSourceDefinitionStrategy) Lookup(ctx context.Context, ds *domain.DataSource, _ map[string]any) ([]*domain.Lookup, error) {
	attr, err := domain.GetDataSourceAttributes[domain.ScheduledDataSourceAttributes](ds.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "data_source_id", ds.ID, "data_source_type", ds.Type, "error", err)
		return nil, err
	}

	return attr.Data, nil
}
