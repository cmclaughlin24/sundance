package processors

import (
	"log/slog"
	"sundance/backend/services/forms/internal/core/ports"
)

type processorOptions struct {
	logger     *slog.Logger
	evaluator  ports.RuleEvaluator
	repository *ports.Repository
	strategies *ports.Strategies
}

func Bootstrap(opts ...func(*processorOptions)) *ports.Processors {
	var po processorOptions
	for _, opt := range opts {
		opt(&po)
	}

	resolver := newSubmissionResolver(po.logger, po.evaluator)
	validator := newSubmissionValidator(po.strategies)
	normalizer := newSubmissionNormalizer(po.logger, po.repository)
	processor := newSubmissionProcessor(po.logger, resolver, validator, normalizer, po.repository)

	return &ports.Processors{
		Submission: processor,
	}
}

func WithLogger(logger *slog.Logger) func(*processorOptions) {
	return func(so *processorOptions) {
		so.logger = logger
	}
}

func WithRuleEvaluator(evaluator ports.RuleEvaluator) func(*processorOptions) {
	return func(so *processorOptions) {
		so.evaluator = evaluator
	}
}

func WithRepository(repository *ports.Repository) func(*processorOptions) {
	return func(so *processorOptions) {
		so.repository = repository
	}
}

func WithStrategies(strategies *ports.Strategies) func(*processorOptions) {
	return func(so *processorOptions) {
		so.strategies = strategies
	}
}
