package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

var (
	ErrVersionStatus = errors.New("invalid version status")
)

type submissionJobsService struct {
	logger                   *slog.Logger
	versionRepository        ports.VersionRepository
	submissionRepository     ports.SubmissionsRepository
	fieldValidatorStrategies ports.FieldValidatorRegistry
}

func NewSubmissionJobsService(logger *slog.Logger, repository *ports.Repository, strategies *ports.Strategies) ports.SubmissionJobsService {
	return &submissionJobsService{
		logger:                   logger,
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

func (s *submissionJobsService) Process(ctx context.Context, command *ports.ProcessSubmissionJobCommand) error {
	if err := command.Validate(); err != nil {
		s.logger.WarnContext(ctx, "submission job process failed; invalid command", "error", err)
		return err
	}

	s.logger.DebugContext(ctx, "processing submission job", "submission_id", command.ID)

	submission, err := s.submissionRepository.FindByID(ctx, command.ID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve submission job", "submission_id", command.ID, "error", err)
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

	validationErr := make([]error, 0)

	for _, page := range version.GetPagesSlice() {
		// TODO: Check page rules and see if page should be evaluated.

		for _, section := range page.GetSectionsSlice() {
			// TODO: Check section rules and see if section should be evaluated.

			for _, field := range section.GetFieldsSlice() {
				// TODO: Check field rules and see if field should be evaluated.

				fieldValidator, err := s.fieldValidatorStrategies.Get(field.Type)
				if err != nil {
					s.logger.ErrorContext(ctx, "failed to prcoess submission; missing field validation strategy", "submission_id", submission.ID, "form_id", version.FormID, "version_id", version.ID, "field_id", field.ID, "field_type", field.Type)
					return err
				}

				if err = fieldValidator.Validate(ctx, field); err != nil {
					validationErr = append(validationErr, err)
				}
			}
		}
	}

	if len(validationErr) > 0 {
		// TODO: Return concat the list of errors into a single error and return.
	}

	return nil
}
