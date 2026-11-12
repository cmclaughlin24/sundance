package strategies

import (
	"log/slog"

	"sundance/backend/pkg/common/stratreg"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type strategyOptions struct {
	logger     *slog.Logger
	repository *ports.Repository
	clients    *ports.Clients
}

func Bootstrap(opts ...func(*strategyOptions)) *ports.Strategies {
	var so strategyOptions
	for _, opt := range opts {
		opt(&so)
	}

	definitionStrategies := stratreg.New[domain.DataSourceType, ports.DataSourceDefinitionStrategy]().
		Set(domain.DataSourceTypeStatic, NewStaticDataSourceDefinitionStrategy(so.logger)).
		Set(domain.DataSourceTypeScheduled, NewScheduledDataSourceDefinitionStrategy(so.logger)).
		Set(domain.DataSourceTypeWebhook, NewWebhookDataSourceDefinitionStrategy(so.logger, so.clients)).
		Set(domain.DataSourceTypeDataLake, NewDataLakeDataSourceDefinitionStrategy(so.logger, so.clients))

	return &ports.Strategies{
		DataSourceDefinition: definitionStrategies,
	}
}

func WithLogger(logger *slog.Logger) func(*strategyOptions) {
	return func(so *strategyOptions) {
		so.logger = logger
	}
}

func WithRepository(repository *ports.Repository) func(*strategyOptions) {
	return func(so *strategyOptions) {
		so.repository = repository
	}
}

func WithClients(clients *ports.Clients) func(*strategyOptions) {
	return func(so *strategyOptions) {
		so.clients = clients
	}
}
