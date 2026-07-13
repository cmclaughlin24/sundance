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
	version, err := p.formVersionsRepository.FindByID(ctx, submission.VersionID)
	if err != nil {
		p.logger.ErrorContext(ctx, "failed to retrieve version for submission job", "submission_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID, "error", err)
		return nil, err
	}

	resolved, err := p.resolver.resolve(ctx, submission, version)
	if err != nil {
		return nil, err
	}

	l := p.logger.With("submission_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID)
	if err := p.validator.validate(ctx, l, resolved); err != nil {
		return nil, err
	}

	return p.normalizer.normalize(ctx, resolved)
}
