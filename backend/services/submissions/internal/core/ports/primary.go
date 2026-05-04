package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
)

type Services struct {
	Submissions SubmissionsService
}

type SubmissionsService interface {
	Find(context.Context) ([]*domain.Submission, error)
	FindByID(context.Context, *FindByIDQuery[domain.SubmissionID]) (*domain.Submission, error)
	FindByReferenceID(context.Context, *FindByIDQuery[domain.ReferenceID]) (*domain.Submission, error)
	Replay(context.Context, ReplaySubmissionCommand) error
}
