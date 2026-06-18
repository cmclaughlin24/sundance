package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"sundance/backend/pkg/database"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/strategies"
)

var (
	ErrVersionStatus = errors.New("invalid version status")
)

type ruleByTypeGetter interface {
	GetRuleByType(domain.RuleType) *domain.Rule
}

type submissionJobsService struct {
	logger                   *slog.Logger
	evaluator                ports.RuleEvaluator
	database                 database.Database
	versionRepository        ports.FormVersionRepository
	submissionRepository     ports.SubmissionsRepository
	fieldValidatorStrategies ports.FieldValidatorRegistry
}

func NewSubmissionJobsService(
	logger *slog.Logger,
	evaluator ports.RuleEvaluator,
	repository *ports.Repository,
	strategies *ports.Strategies,
) ports.SubmissionJobsAPI {
	return &submissionJobsService{
		logger:                   logger,
		evaluator:                evaluator,
		database:                 repository.Database,
		versionRepository:        repository.FormVersions,
		submissionRepository:     repository.Submissions,
		fieldValidatorStrategies: strategies.FieldValidator,
	}
}

func (s *submissionJobsService) Find(ctx context.Context, query ports.FindSubmissionJobsQuery) ([]domain.SubmissionID, error) {
	s.logger.DebugContext(ctx, "listing submission jobs")

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "submission job listing failed; invalid query", "error", err)
		return nil, err
	}

	ids, err := s.submissionRepository.FindJobs(ctx, &ports.FindSubmissionsFilter{
		Take:     query.Take,
		Statuses: []domain.SubmissionStatus{domain.SubmissionStatusPending},
	})

	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve submission jobs", "error", err)
		return nil, err
	}

	return ids, nil
}

func (s *submissionJobsService) Process(ctx context.Context, id domain.SubmissionID) error {
	s.logger.DebugContext(ctx, "processing submission job", "submission_id", id)

	submission, err := s.submissionRepository.FindByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve submission job", "submission_id", id, "error", err)
		return err
	}

	if submission.Status == domain.SubmissionStatusAccepted || submission.Status == domain.SubmissionStatusRejected {
		s.logger.WarnContext(ctx, "skipping submission job; invalid status", "submission_id", submission.ID, "status", submission.Status)
		return nil
	}

	err = s.sanitize(ctx, submission)

	if err := s.recordAttempt(ctx, submission, err); err != nil {
		return err
	}

	s.logger.InfoContext(ctx, "submission job recorded", "submission_id", submission.ID, "status", submission.Status)

	return err
}

func (s *submissionJobsService) sanitize(ctx context.Context, submission *domain.Submission) error {
	version, err := s.versionRepository.FindByID(ctx, submission.VersionID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve version for submission job", "submission_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID, "error", err)
		return err
	}

	if version.Status == domain.FormVersionStatusDraft {
		s.logger.WarnContext(ctx, "skipping submission job; invalid status", "submission_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID, "version_status", version.Status)
		return ErrVersionStatus
	}

	evalCtx := make(ports.RuleEvaluationContext, len(submission.Values))
	for _, field := range version.FlatFields() {
		if val, ok := submission.GetFieldValue(field.ID); ok {
			evalCtx[field.Key] = val.Value
		}
	}

pageLoop:
	for _, page := range version.GetPages() {
		visible, err := s.shouldValidate(ctx, page, evalCtx)

		if err != nil {
			return err
		}

		if !visible {
			continue pageLoop
		}

	sectionLoop:
		for _, section := range page.GetSections() {
			visible, err := s.shouldValidate(ctx, section, evalCtx)

			if err != nil {
				return err
			}

			if !visible {
				continue sectionLoop
			}

		fieldLoop:
			for _, field := range section.GetFields() {
				visible, err := s.shouldValidate(ctx, field, evalCtx)

				if err != nil {
					return err
				}

				if !visible {
					continue fieldLoop
				}

				required, err := s.isRequired(ctx, field, evalCtx)
				if err != nil {
					return err
				}

				if required != nil {
					field.Attributes.SetIsRequired(*required)
				}

				if err := s.validateField(ctx, field, submission); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *submissionJobsService) validateField(ctx context.Context, field *domain.Field, submission *domain.Submission) error {
	fieldValidator, err := s.fieldValidatorStrategies.Get(field.Type)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to process submission; missing field validation strategy", "submission_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID, "field_id", field.ID, "field_type", field.Type)
		return err
	}

	value, ok := submission.GetFieldValue(field.ID)
	if !ok {
		if field.Attributes.GetIsRequired() {
			s.logger.WarnContext(ctx, "submission validation failed; required field missing", "submission_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID, "field_id", field.ID, "field_key", field.Key)
			return fmt.Errorf("%w; id=%s key=%s", strategies.ErrFieldRequired, field.ID, field.Key)
		}

		return nil
	}

	if err = fieldValidator.Validate(ctx, *field, *value); err != nil {
		return err
	}

	return nil
}

func (s *submissionJobsService) isRequired(ctx context.Context, rg ruleByTypeGetter, evalCtx ports.RuleEvaluationContext) (*bool, error) {
	rule := rg.GetRuleByType(domain.RuleTypeRequired)
	if rule == nil {
		return nil, nil
	}

	result, err := s.evaluator.Evaluate(ctx, rule, evalCtx)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *submissionJobsService) shouldValidate(ctx context.Context, rg ruleByTypeGetter, evalCtx ports.RuleEvaluationContext) (bool, error) {
	rule := rg.GetRuleByType(domain.RuleTypeVisible)

	if rule == nil {
		return true, nil
	}

	return s.evaluator.Evaluate(ctx, rule, evalCtx)
}

func (s *submissionJobsService) recordAttempt(ctx context.Context, submission *domain.Submission, err error) error {
	if err == nil {
		submission.Accept()
	} else if shouldReject(err) {
		submission.Reject(err)
	} else {
		submission.Fail(err)
	}

	txCtx, err := s.database.BeginTx(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to begin transaction", "submission_id", submission.ID, "error", err)
		return err
	}

	defer s.database.RollbackTx(txCtx)

	_, err = s.submissionRepository.Upsert(txCtx, submission)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to upsert submission", "submission_id", submission.ID, "error", err)
		return err
	}

	if err := s.database.CommitTx(txCtx); err != nil {
		s.logger.ErrorContext(ctx, "failed to commit transaction", "submission_id", submission.ID, "error", err)
		return err
	}

	return nil
}

func shouldReject(err error) bool {
	return errors.Is(err, ErrVersionStatus) ||
		errors.Is(err, strategies.ErrFieldValidation) ||
		errors.Is(err, strategies.ErrFieldRequired) ||
		errors.Is(err, strategies.ErrFieldTypeValue)
}
