package services

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/submissions/internal/core/domain"
	"github.com/cmclaughlin24/sundance/submissions/internal/core/ports"
)

type SubmissionsService struct {
	logger     *log.Logger
	repository *ports.Repository
}

func NewSubmissionsService(logger *log.Logger, repository *ports.Repository) *SubmissionsService {
	return &SubmissionsService{
		logger:     logger,
		repository: repository,
	}
}

func (s *SubmissionsService) Find(context.Context) ([]*domain.Submission, error) {
	return nil, nil
}

func (s *SubmissionsService) FindById(context.Context, ports.FindByIdQuery[domain.SubmissionID]) (*domain.Submission, error) {
	return nil, nil
}

func (s *SubmissionsService) FindByReferenceId(context.Context, ports.FindByIdQuery[domain.ReferenceID]) (*domain.Submission, error) {
	return nil, nil
}

func (s *SubmissionsService) FindAttempts(context.Context) ([]*domain.SubmissionAttempt, error) {
	return nil, nil
}

func (s *SubmissionsService) Replay(context.Context, ports.ReplaySubmissionCommand) error {
	return nil
}
