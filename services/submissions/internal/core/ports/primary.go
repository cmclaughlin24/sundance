package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/submissions/internal/core/domain"
)

type Services struct {
	Submissions SubmissionsService
}

type SubmissionsService interface {
	Find(context.Context) ([]*domain.Submission, error)
	FindById(context.Context, *FindByIdQuery[domain.SubmissionID]) (*domain.Submission, error)
	FindByReferenceId(context.Context, *FindByIdQuery[domain.ReferenceID]) (*domain.Submission, error)
	FindAttempts(context.Context) ([]*domain.SubmissionAttempt, error)
	Replay(context.Context, ReplaySubmissionCommand) error 
}
