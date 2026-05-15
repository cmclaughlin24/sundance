package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

type submissionJobsService struct {
	logger               *slog.Logger
	versionRepository    ports.VersionRepository
	submissionRepository ports.SubmissionsRepository
}

func NewSubmissionJobsService(logger *slog.Logger, repository *ports.Repository) ports.SubmissionJobsService {
	return &submissionJobsService{
		logger:               logger,
		versionRepository:    repository.Versions,
		submissionRepository: repository.Submissions,
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
		s.logger.WarnContext(ctx, "failed to process submission job; invalid status", "submssion_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID, "version_status", version.Status)
		return fmt.Errorf("")
	}

	// 3) Dynamically create a form definition struct.
	// 4) Validate the submission against the form definition struct.

	return nil
}
