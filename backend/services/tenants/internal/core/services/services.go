package services

import (
	"log/slog"
	"time"

	"sundance/backend/services/tenants/internal/core/ports"
)

// Package declaration for the current time function. Allows for easier testing by enabling the injection of a
// mock time function.
var Now = time.Now

type serviceOptions struct {
	logger     *slog.Logger
	repository *ports.Repository
	strategies *ports.Strategies
	clients    *ports.Clients
}

func Bootstrap(opts ...func(*serviceOptions)) *ports.Services {
	var so serviceOptions
	for _, opt := range opts {
		opt(&so)
	}

	return &ports.Services{
		Tenants:        NewTenantsService(so.logger, so.repository),
		DataSources:    NewDataSourcesService(so.logger, so.repository, so.strategies),
		DataSourceJobs: NewDataSourcesJobService(so.logger, so.repository, so.clients),
	}
}

func WithClients(clients *ports.Clients) func(*serviceOptions) {
	return func(so *serviceOptions) {
		so.clients = clients
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
