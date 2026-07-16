package services

import (
	"context"
	"errors"
	"log/slog"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/processors"
	"sundance/backend/services/forms/internal/core/strategies"
)

type submissionJobsService struct {
	logger                *slog.Logger
	processors            ports.SubmissionProcessor
	submissionsRepository ports.SubmissionsRepository
}

func NewSubmissionJobsService(
	logger *slog.Logger,
	processors *ports.Processors,
	repository *ports.Repository,
) ports.SubmissionJobsAPI {
	return &submissionJobsService{
		logger:                logger,
		processors:            processors.Submission,
		submissionsRepository: repository.Submissions,
	}
}

func (s *submissionJobsService) Find(ctx context.Context, query ports.FindSubmissionJobsQuery) ([]domain.SubmissionID, error) {
	s.logger.DebugContext(ctx, "listing submission jobs")

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "submission job listing failed; invalid query", "error", err)
		return nil, err
	}

	ids, err := s.submissionsRepository.FindJobs(ctx, &ports.FindSubmissionsFilter{
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

	submission, err := s.submissionsRepository.FindByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve submission job", "submission_id", id, "error", err)
		return err
	}

	if submission.Status == domain.SubmissionStatusAccepted || submission.Status == domain.SubmissionStatusRejected {
		s.logger.WarnContext(ctx, "skipping submission job; invalid status", "submission_id", submission.ID, "status", submission.Status)
		return nil
	}

	facts, err := s.processors.Process(ctx, submission)
	if err == nil {
		submission.Accept(facts)
	} else if shouldReject(err) {
		submission.Reject(err)
	} else {
		submission.Fail(err)
	}

	if err := s.updateSubmission(ctx, submission); err != nil {
		return err
	}

	s.logger.InfoContext(ctx, "submission job recorded", "submission_id", submission.ID, "status", submission.Status)

	return nil
}

func (s *submissionJobsService) updateSubmission(ctx context.Context, submission *domain.Submission) error {
	_, err := s.submissionsRepository.Upsert(ctx, submission)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to upsert submission", "submission_id", submission.ID, "error", err)
		return err
	}

	return nil
}

func shouldReject(err error) bool {
	return errors.Is(err, domain.ErrInvalidVersionStatus) ||
		errors.Is(err, strategies.ErrFieldValidation) ||
		errors.Is(err, strategies.ErrFieldRequired) ||
		errors.Is(err, strategies.ErrFieldTypeValue) ||
		errors.Is(err, processors.ErrMissingCollectionIndex)
}
