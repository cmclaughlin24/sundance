package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

var (
	ErrVersionStatus = errors.New("invalid version status")
)

type ruleGetter interface {
	GetRule(domain.RuleType) *domain.Rule
}

type submissionJobsService struct {
	logger                   *slog.Logger
	evaluator                ports.RuleEvaluator
	versionRepository        ports.VersionRepository
	submissionRepository     ports.SubmissionsRepository
	fieldValidatorStrategies ports.FieldValidatorRegistry
}

func NewSubmissionJobsService(
	logger *slog.Logger,
	evaluator ports.RuleEvaluator,
	repository *ports.Repository,
	strategies *ports.Strategies,
) ports.SubmissionJobsService {
	return &submissionJobsService{
		logger:                   logger,
		evaluator:                evaluator,
		versionRepository:        repository.Versions,
		submissionRepository:     repository.Submissions,
		fieldValidatorStrategies: strategies.FieldValidator,
	}
}

func (s *submissionJobsService) Find(ctx context.Context, query *ports.FindSubmissionJobsQuery) ([]domain.SubmissionID, error) {
	s.logger.DebugContext(ctx, "listing submission jobs")

	ids, err := s.submissionRepository.FindJobs(ctx, &ports.FindSubmissionsFilter{Statuses: []domain.SubmissionStatus{domain.SubmissionStatusPending}})
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

	if submission.Status != domain.SubmissionStatusPending {
		s.logger.WarnContext(ctx, "skipping submission job; invalid status", "submission_id", submission.ID, "status", submission.Status)
		return nil
	}

	version, err := s.versionRepository.FindByID(ctx, submission.FormID, submission.VersionID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve version for submission job", "submission_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID, "error", err)
		return err
	}

	if version.Status == domain.VersionStatusDraft {
		s.logger.WarnContext(ctx, "failed to process submission job; invalid status", "submission_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID, "version_status", version.Status)
		return ErrVersionStatus
	}

	if err := s.validate(ctx, version, submission); err != nil {
		// TODO: Decide how to handle errors based on type.
		return err
	}

	return nil
}

func (s *submissionJobsService) validate(ctx context.Context, version *domain.Version, submission *domain.Submission) error {
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
			return fmt.Errorf("field id=%s key=%s is required", field.ID, field.Key)
		}

		return nil
	}

	if err = fieldValidator.Validate(ctx, *field, *value); err != nil {
		return err
	}

	return nil
}

func (s submissionJobsService) shouldValidate(ctx context.Context, getter ruleGetter, evalCtx ports.RuleEvaluationContext) (bool, error) {
	rule := getter.GetRule(domain.RuleTypeVisible)
	if rule == nil {
		return true, nil
	}
	return s.evaluator.Evaluate(ctx, rule, evalCtx)
}
