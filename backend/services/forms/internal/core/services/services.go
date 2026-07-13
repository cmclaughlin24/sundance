package services

import (
	"log/slog"

	"sundance/backend/services/forms/internal/core/ports"
)

type serviceOptions struct {
	logger     *slog.Logger
	repository *ports.Repository
	processors *ports.Processors
}

func Bootstrap(opts ...func(*serviceOptions)) *ports.API {
	var so serviceOptions
	for _, opt := range opts {
		opt(&so)
	}

	return &ports.API{
		Tags:           NewTagsService(so.logger, so.repository),
		Forms:          NewFormsService(so.logger, so.repository),
		Submissions:    NewSubmissionsService(so.logger, so.processors, so.repository),
		SubmissionJobs: NewSubmissionJobsService(so.logger, so.processors, so.repository),
	}
}

func WithLogger(logger *slog.Logger) func(*serviceOptions) {
	return func(so *serviceOptions) {
		so.logger = logger
	}
}

func WithProcessors(processors *ports.Processors) func(*serviceOptions) {
	return func(so *serviceOptions) {
		so.processors = processors
	}
}

func WithRepository(repository *ports.Repository) func(*serviceOptions) {
	return func(so *serviceOptions) {
		so.repository = repository
	}
}
