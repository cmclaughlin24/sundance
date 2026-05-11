package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
)

type Services struct {
	Submissions    SubmissionsService
	SubmissionJobs SubmissionJobsService
}

type SubmissionsService interface {
	Find(context.Context, *FindSubmissionsQuery) ([]*domain.Submission, error)
	FindByID(context.Context, *FindSubmissionByIDQuery[domain.SubmissionID]) (*domain.Submission, error)
	FindByReferenceID(context.Context, *FindSubmissionByIDQuery[domain.ReferenceID]) (*domain.Submission, error)
	Create(context.Context, *CreateSubmissionCommand) (*domain.Submission, error)
	Replay(context.Context, *ReplaySubmissionCommand) error
}

type SubmissionJobsService interface {
	Find(context.Context, *FindSubmissionJobsQuery) ([]*domain.Submission, error)
	Process(context.Context) error
}
