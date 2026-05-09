package services

import (
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
)

type serviceOptions struct {
	logger     *slog.Logger
	repository *ports.Repository
	strategies *ports.Strategies
}

func Bootstrap(opts ...func(*serviceOptions)) *ports.Services {
	var so serviceOptions
	for _, opt := range opts {
		opt(&so)
	}

	return &ports.Services{
		Tenants:        NewTenantsService(so.logger, so.repository),
		DataSources:    NewDataSourcesService(so.logger, so.repository, so.strategies),
		DataSourceJobs: NewDataSourcesJobService(so.logger, so.repository),
	}
}

func WithLogger(logger *slog.Logger) func(*serviceOptions) {
	return func(so *serviceOptions) {
		so.logger = logger
	}
}

func WithRepository(repository *ports.Repository) func(*serviceOptions) {
	return func(so *serviceOptions) {
		so.repository = repository
	}
}

func WithStrategies(strategies *ports.Strategies) func(*serviceOptions) {
	return func(so *serviceOptions) {
		so.strategies = strategies
	}
}
