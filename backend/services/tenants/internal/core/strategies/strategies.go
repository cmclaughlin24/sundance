package strategies

import (
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/stratreg"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
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

	lookupStrategies := stratreg.New[domain.DataSourceType, ports.LookupStrategy]().
		Set(domain.DataSourceTypeStatic, NewStaticLookupStrategy(so.logger)).
		Set(domain.DataSourceTypeScheduled, NewScheduledLookupStrategy(so.logger)).
		Set(domain.DataSourceTypeWebhook, NewWebhookLookupStrategy(so.logger, so.clients))

	return &ports.Strategies{
		Lookups: lookupStrategies,
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
