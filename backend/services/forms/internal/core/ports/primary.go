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
	FindVersions(context.Context, *FindFormVersionsQuery) ([]*domain.FormVersion, error)
	FindVersion(context.Context, *FindFormVersionByIDQuery) (*domain.FormVersion, error)
	CreateVersion(context.Context, *CreateFormVersionCommand) (*domain.FormVersion, error)
	UpdateVersion(context.Context, *UpdateFormVersionCommand) (*domain.FormVersion, error)
	PublishVersion(context.Context, *PublishFormVersionCommand) (*domain.FormVersion, error)
	RetireVersion(context.Context, *RetireFormVersionCommand) (*domain.FormVersion, error)
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
