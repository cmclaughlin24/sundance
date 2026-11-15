package processors

import (
	"context"
	"log/slog"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type submissionProcessor struct {
	logger                 *slog.Logger
	resolver               *submissionResolver
	validator              *submissionValidator
	normalizer             *submissionNormalizer
	formVersionsRepository ports.FormVersionRepository
}

func newSubmissionProcessor(
	logger *slog.Logger,
	resolver *submissionResolver,
	validator *submissionValidator,
	normalizer *submissionNormalizer,
	repository *ports.Repository,
) ports.SubmissionProcessor {
	return &submissionProcessor{
		logger:                 logger,
		resolver:               resolver,
		validator:              validator,
		normalizer:             normalizer,
		formVersionsRepository: repository.FormVersions,
	}
}

func (p *submissionProcessor) Process(ctx context.Context, submission *domain.Submission) ([]*domain.CanonicalFact, error) {
	l := p.logger.With("submission_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID)

	version, err := p.formVersionsRepository.FindByID(ctx, submission.VersionID)
	if err != nil {
		l.ErrorContext(ctx, "failed to retrieve version for submission job", "error", err)
		return nil, err
	}

	if err := p.validator.validateStructure(ctx, l, submission, version); err != nil {
		return nil, err
	}

	resolved, err := p.resolver.resolve(ctx, submission, version)
	if err != nil {
		return nil, err
	}

	if err := p.validator.validateResolved(ctx, l, resolved); err != nil {
		return nil, err
	}

	return p.normalizer.normalize(ctx, resolved)
}
