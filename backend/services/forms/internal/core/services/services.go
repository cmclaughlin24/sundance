package services

import (
	"log/slog"

	"sundance/backend/services/forms/internal/core/ports"
)

type serviceOptions struct {
	logger     *slog.Logger
	evaluator  ports.RuleEvaluator
	repository *ports.Repository
	strategies *ports.Strategies
}

func Bootstrap(opts ...func(*serviceOptions)) *ports.API {
	var so serviceOptions
	for _, opt := range opts {
		opt(&so)
	}

	return &ports.API{
		CanonicalTags:  NewCanonicalTagService(so.logger, so.repository),
		Forms:          NewFormsService(so.logger, so.repository),
		Submissions:    NewSubmissionsService(so.logger, so.repository),
		SubmissionJobs: NewSubmissionJobsService(so.logger, so.evaluator, so.repository, so.strategies),
	}
}

func WithLogger(logger *slog.Logger) func(*serviceOptions) {
	return func(so *serviceOptions) {
		so.logger = logger
	}
}

func WithRuleEvaluator(evaluator ports.RuleEvaluator) func(*serviceOptions) {
	return func(so *serviceOptions) {
		so.evaluator = evaluator
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
