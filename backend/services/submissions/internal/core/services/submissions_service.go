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
	logger     *log.Logger
	repository ports.SubmissionsRepository
}

func NewSubmissionsService(logger *log.Logger, repository *ports.Repository) ports.SubmissionsService {
	return &SubmissionsService{
		logger:     logger,
		repository: repository.Submissions,
	}
}

func (s *SubmissionsService) Find(ctx context.Context) ([]*domain.Submission, error) {
	return s.repository.Find(ctx)
}

func (s *SubmissionsService) FindByID(ctx context.Context, query *ports.FindByIDQuery[domain.SubmissionID]) (*domain.Submission, error) {
	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	submission, err := s.repository.FindByID(ctx, query.ID)

	if err != nil {
		return nil, err
	}

	if submission.TenantID != query.TenantID {
		return nil, common.ErrUnauthorized
	}

	return submission, nil
}

func (s *SubmissionsService) FindByReferenceID(ctx context.Context, query *ports.FindByIDQuery[domain.ReferenceID]) (*domain.Submission, error) {
	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	submission, err := s.repository.FindByReferenceID(ctx, query.ID)

	if err != nil {
		return nil, err
	}

	if submission.TenantID != query.TenantID {
		return nil, common.ErrUnauthorized
	}

	return submission, nil
}

func (s *SubmissionsService) Create(ctx context.Context, command ports.CreateSubmissionCommand) (*domain.Submission, error) {
	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	submission, err := domain.NewSubmission(command.TenantID, command.FormID, command.VersionID, command.Payload)
	if err != nil {
		return nil, err
	}

	submission, err = s.repository.Upsert(ctx, submission)
	if err != nil {
		return nil, err
	}

	return submission, nil
}

func (s *SubmissionsService) Replay(context.Context, ports.ReplaySubmissionCommand) error {
	return nil
}
