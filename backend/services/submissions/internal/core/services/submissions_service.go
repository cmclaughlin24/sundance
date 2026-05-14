package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
)

type submissionsService struct {
	logger     *slog.Logger
	repository ports.SubmissionsRepository
}

func NewSubmissionsService(logger *slog.Logger, repository *ports.Repository) ports.SubmissionsService {
	return &submissionsService{
		logger:     logger,
		repository: repository.Submissions,
	}
}

func (s *submissionsService) Find(ctx context.Context, query *ports.FindSubmissionsQuery) ([]*domain.Submission, error) {
	s.logger.DebugContext(ctx, "listing submissions", "tenant_id", query.TenantID)

	if err := validate.ValidateStruct(query); err != nil {
		s.logger.WarnContext(ctx, "submission listing failed; invalid query", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	submissions, err := s.repository.Find(ctx, &ports.FindSubmissionsFilter{
		TenantID: query.TenantID,
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve submissions", "tenant_id", query.TenantID, "error", err)
		return nil, err
	}

	return submissions, nil
}

func (s *submissionsService) FindByID(ctx context.Context, query *ports.FindSubmissionByIDQuery[domain.SubmissionID]) (*domain.Submission, error) {
	s.logger.DebugContext(ctx, "finding submission", "tenant_id", query.TenantID, "submission_id", query.ID)

	if err := validate.ValidateStruct(query); err != nil {
		s.logger.WarnContext(ctx, "submission find failed; invalid query", "tenant_id", query.TenantID, "submission_id", query.ID, "error", err)
		return nil, err
	}

	submission, err := s.repository.FindByID(ctx, query.ID)
	if err != nil {
		s.logFindByIDError(ctx, err, query.ID)
		return nil, err
	}

	if submission.TenantID != query.TenantID {
		s.logger.WarnContext(ctx, "unauthorized submission access", "tenant_id", query.TenantID, "submission_id", query.ID)
		return nil, common.ErrUnauthorized
	}

	return submission, nil
}

func (s *submissionsService) FindByReferenceID(ctx context.Context, query *ports.FindSubmissionByIDQuery[domain.ReferenceID]) (*domain.Submission, error) {
	s.logger.DebugContext(ctx, "finding submission by reference", "tenant_id", query.TenantID, "reference_id", query.ID)

	if err := validate.ValidateStruct(query); err != nil {
		s.logger.WarnContext(ctx, "submission find failed; invalid query", "tenant_id", query.TenantID, "reference_id", query.ID, "error", err)
		return nil, err
	}

	submission, err := s.repository.FindByReferenceID(ctx, query.ID)
	if err != nil {
		s.logFindByReferenceIDError(ctx, err, query.ID)
		return nil, err
	}

	if submission.TenantID != query.TenantID {
		s.logger.WarnContext(ctx, "unauthorized submission access", "tenant_id", query.TenantID, "reference_id", query.ID)
		return nil, common.ErrUnauthorized
	}

	return submission, nil
}

func (s *submissionsService) Create(ctx context.Context, command *ports.CreateSubmissionCommand) (*domain.Submission, error) {
	s.logger.DebugContext(ctx, "creating submission", "tenant_id", command.TenantID)

	if err := validate.ValidateStruct(command); err != nil {
		s.logger.WarnContext(ctx, "submission creation failed; invalid command", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	submission, err := s.repository.FindByIdempotencyID(ctx, command.IdempotencyID)
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		s.logger.ErrorContext(ctx, "failed to check submission existence", "tenant_id", command.TenantID, "submission_idempotency_id", command.IdempotencyID, "error", err)
		return nil, err
	}

	if submission != nil {
		s.logger.InfoContext(ctx, "submission exists", "tenant_id", command.TenantID, "submission_id", submission.ID, "submission_idempotency_id", command.IdempotencyID)
		return submission, nil
	}

	submission, err = domain.NewSubmission(
		command.TenantID,
		command.FormID,
		command.VersionID,
		command.IdempotencyID,
		command.Payload,
	)
	if err != nil {
		s.logger.WarnContext(ctx, "submission creation failed; domain invariant violation", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	submission, err = s.repository.Upsert(ctx, submission)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist submission", "tenant_id", command.TenantID, "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "submission created", "tenant_id", command.TenantID, "submission_id", submission.ID)

	return submission, nil
}

func (s *submissionsService) Replay(ctx context.Context, command *ports.ReplaySubmissionCommand) error {
	s.logger.DebugContext(ctx, "replaying submission", "tenant_id", command.TenantID, "submission_id", command.ID)

	if err := validate.ValidateStruct(command); err != nil {
		s.logger.WarnContext(ctx, "submission replay failed; invalid command", "tenant_id", command.TenantID, "submission_id", command.ID, "error", err)
		return err
	}

	submission, err := s.repository.FindByID(ctx, command.ID)
	if err != nil {
		s.logFindByIDError(ctx, err, command.ID)
		return err
	}

	if submission.TenantID != command.TenantID {
		s.logger.WarnContext(ctx, "unauthorized form access", "tenant_id", command.TenantID, "submission_id", command.ID)
		return common.ErrUnauthorized
	}

	submission.Reset()

	if _, err := s.repository.Upsert(ctx, submission); err != nil {
		s.logger.ErrorContext(ctx, "submission replay failed", "tenant_id", command.TenantID, "submission_id", submission.ID, "error", err)
		return err
	}

	s.logger.InfoContext(ctx, "submission replayed", "tenant_id", command.TenantID, "submission_id", submission.ID)

	return nil
}

func (s *submissionsService) logFindByIDError(ctx context.Context, err error, id domain.SubmissionID) {
	if errors.Is(err, common.ErrNotFound) {
		s.logger.WarnContext(ctx, "submission not found", "submission_id", id)
	} else {
		s.logger.ErrorContext(ctx, "failed to retrieve submission", "submission_id", id, "error", err)
	}
}

func (s *submissionsService) logFindByReferenceIDError(ctx context.Context, err error, id domain.ReferenceID) {
	if errors.Is(err, common.ErrNotFound) {
		s.logger.WarnContext(ctx, "submission not found", "reference_id", id)
	} else {
		s.logger.ErrorContext(ctx, "failed to retrieve submission", "reference_id", id, "error", err)
	}
}
