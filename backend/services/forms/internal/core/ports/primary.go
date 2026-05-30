package ports

import (
	"context"

	"sundance/backend/services/forms/internal/core/domain"
)

type API struct {
	CanonicalTags  CanonicalTagAPI
	Forms          FormsAPI
	Submissions    SubmissionsAPI
	SubmissionJobs SubmissionJobsAPI
}

type CanonicalTagAPI interface {
	Find(context.Context) ([]*domain.CanonicalTag, error)
	FindById(context.Context) (*domain.CanonicalTag, error)
	Create(context.Context, any) (*domain.CanonicalTag, error)
	Delete(context.Context, any) error
	FindVersions(context.Context, any) error
	FindVersion(context.Context, any) error
	CreateVersion(context.Context, any) (*domain.CanonicalTagVersion, error)
	UpdateVersion(context.Context, any) (*domain.CanonicalTagVersion, error)
	PublishVersion(context.Context, any) (*domain.CanonicalTagVersion, error)
	DeprecateVersion(context.Context, any) (*domain.CanonicalTagVersion, error)
	RetireVersion(context.Context, any) (*domain.CanonicalTagVersion, error)
}

type FormsAPI interface {
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

type SubmissionsAPI interface {
	Find(context.Context, *FindSubmissionsQuery) ([]*domain.Submission, error)
	FindByID(context.Context, *FindSubmissionByIDQuery[domain.SubmissionID]) (*domain.Submission, error)
	FindByReferenceID(context.Context, *FindSubmissionByIDQuery[domain.ReferenceID]) (*domain.Submission, error)
	Create(context.Context, *CreateSubmissionCommand) (*domain.Submission, error)
	Replay(context.Context, *ReplaySubmissionCommand) error
}

type SubmissionJobsAPI interface {
	Find(context.Context, *FindSubmissionJobsQuery) ([]domain.SubmissionID, error)
	Process(context.Context, domain.SubmissionID) error
}
