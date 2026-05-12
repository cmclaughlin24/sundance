package services

import (
	"context"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
)

type submissionJobsService struct {
	logger     *slog.Logger
	repository ports.SubmissionsRepository
}

func NewSubmissionJobsService(logger *slog.Logger, repository *ports.Repository) ports.SubmissionJobsService {
	return &submissionJobsService{
		logger:     logger,
		repository: repository.Submissions,
	}
}

func (s *submissionJobsService) Find(ctx context.Context, query *ports.FindSubmissionJobsQuery) ([]*domain.Submission, error) {
	s.logger.DebugContext(ctx, "listing submission jobs")

	submissions, err := s.repository.Find(ctx, &ports.FindSubmissionsFilter{Statuses: []domain.SubmissionStatus{domain.SubmissionStatusPending}})
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve submission jobs", "error", err)
		return nil, err
	}

	return submissions, nil
}

func (s *submissionJobsService) Process(context.Context) error {
	return nil
}
