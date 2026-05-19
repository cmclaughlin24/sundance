package ports

import (
	"context"

	"sundance/backend/services/forms/internal/core/domain"
)

type Services struct {
	Forms          FormsService
	Submissions    SubmissionsService
	SubmissionJobs SubmissionJobsService
}

type FormsService interface {
	Find(context.Context, *FindFormsQuery) ([]*domain.Form, error)
	FindByID(context.Context, *FindFormsByIDQuery) (*domain.Form, error)
	Create(context.Context, *CreateFormCommand) (*domain.Form, error)
	Update(context.Context, *UpdateFormCommand) (*domain.Form, error)
	Delete(context.Context, *RemoveFormCommand) error
	FindVersions(context.Context, *FindVersionsQuery) ([]*domain.Version, error)
	FindVersion(context.Context, *FindVersionByIDQuery) (*domain.Version, error)
	CreateVersion(context.Context, *CreateVersionCommand) (*domain.Version, error)
	UpdateVersion(context.Context, *UpdateVersionCommand) (*domain.Version, error)
	PublishVersion(context.Context, *PublishVersionCommand) (*domain.Version, error)
	RetireVersion(context.Context, *RetireVersionCommand) (*domain.Version, error)
}

type SubmissionsService interface {
	Find(context.Context, *FindSubmissionsQuery) ([]*domain.Submission, error)
	FindByID(context.Context, *FindSubmissionByIDQuery[domain.SubmissionID]) (*domain.Submission, error)
	FindByReferenceID(context.Context, *FindSubmissionByIDQuery[domain.ReferenceID]) (*domain.Submission, error)
	Create(context.Context, *CreateSubmissionCommand) (*domain.Submission, error)
	Replay(context.Context, *ReplaySubmissionCommand) error
}

type SubmissionJobsService interface {
	Find(context.Context, *FindSubmissionJobsQuery) ([]domain.SubmissionID, error)
	Process(context.Context, domain.SubmissionID) error
}
