package services

import (
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
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
		Forms:          NewFormsService(so.logger, so.repository),
		Submissions:    NewSubmissionsService(so.logger, so.repository),
		SubmissionJobs: NewSubmissionJobsService(so.logger, so.repository, so.strategies),
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
