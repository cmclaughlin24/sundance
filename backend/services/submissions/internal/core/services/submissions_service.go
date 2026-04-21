package services

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
)

type SubmissionsService struct {
	logger                *log.Logger
	submissionsRepository ports.SubmissionsRepository
}

func NewSubmissionsService(logger *log.Logger, repository *ports.Repository) *SubmissionsService {
	return &SubmissionsService{
		logger:                logger,
		submissionsRepository: repository.Submissions,
	}
}

func (s *SubmissionsService) Find(ctx context.Context) ([]*domain.Submission, error) {
	return s.submissionsRepository.Find(ctx)
}

func (s *SubmissionsService) FindById(ctx context.Context, query *ports.FindByIdQuery[domain.SubmissionID]) (*domain.Submission, error) {
	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	submission, err := s.submissionsRepository.FindById(ctx, query.ID)

	if err != nil {
		return nil, err
	}

	if submission.TenantID != query.TenantID {
		return nil, common.ErrUnauthorized
	}

	return submission, nil
}

func (s *SubmissionsService) FindByReferenceId(ctx context.Context, query *ports.FindByIdQuery[domain.ReferenceID]) (*domain.Submission, error) {
	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	submission, err := s.submissionsRepository.FindByReferenceId(ctx, query.ID)

	if err != nil {
		return nil, err
	}

	if submission.TenantID != query.TenantID {
		return nil, common.ErrUnauthorized
	}

	return submission, nil
}

func (s *SubmissionsService) FindAttempts(context.Context) ([]*domain.SubmissionAttempt, error) {
	return nil, nil
}

func (s *SubmissionsService) Replay(context.Context, ports.ReplaySubmissionCommand) error {
	return nil
}
