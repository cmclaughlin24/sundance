package ports

import (
	"context"

	"sundance/backend/services/forms/internal/core/domain"
)

type API struct {
	Tags           TagsAPI
	Forms          FormsAPI
	Submissions    SubmissionsAPI
	SubmissionJobs SubmissionJobsAPI
}

type FormsAPI interface {
	Find(context.Context, FindFormsQuery) ([]*domain.Form, error)
	FindByID(context.Context, FindByIDQuery[domain.FormID]) (*domain.Form, error)
	Create(context.Context, CreateFormCommand) (*domain.Form, error)
	Update(context.Context, UpdateFormCommand) (*domain.Form, error)
	Delete(context.Context, DeleteCommand[domain.FormID]) error
	FindVersions(context.Context, FindFormVersionsQuery) ([]*domain.FormVersion, error)
	FindVersion(context.Context, FindFormVersionByIDQuery) (*domain.FormVersion, error)
	CreateVersion(context.Context, *CreateFormVersionCommand) (*domain.FormVersion, error)
	UpdateVersion(context.Context, *UpdateFormVersionCommand) (*domain.FormVersion, error)
	PublishVersion(context.Context, PublishFormVersionCommand) (*domain.FormVersion, error)
	RetireVersion(context.Context, RetireFormVersionCommand) (*domain.FormVersion, error)
}

type SubmissionsAPI interface {
	Find(context.Context, FindSubmissionsQuery) ([]*domain.Submission, error)
	FindByID(context.Context, FindSubmissionByIDQuery[domain.SubmissionID]) (*domain.Submission, error)
	FindByReferenceID(context.Context, FindSubmissionByIDQuery[domain.ReferenceID]) (*domain.Submission, error)
	Create(context.Context, *CreateSubmissionCommand) (*domain.Submission, error)
	Replay(context.Context, ReplaySubmissionCommand) error
}

type SubmissionJobsAPI interface {
	Find(context.Context, FindSubmissionJobsQuery) ([]domain.SubmissionID, error)
	Process(context.Context, domain.SubmissionID) error
}

type TagsAPI interface {
	Find(context.Context, FindTagsQuery) ([]*domain.Tag, error)
	FindById(context.Context, FindByIDQuery[domain.TagID]) (*domain.Tag, error)
	Create(context.Context, CreateTagCommand) (*domain.Tag, error)
	Update(context.Context, UpdateTagCommand) (*domain.Tag, error)
	Delete(context.Context, DeleteCommand[domain.TagID]) error
	FindVersions(context.Context, FindTagVersionsQuery) ([]*domain.TagVersion, error)
	FindVersion(context.Context, FindTagVersionQuery) (*domain.TagVersion, error)
	CreateVersion(context.Context, CreateTagVersionCommand) (*domain.TagVersion, error)
	PublishVersion(context.Context, TransitionTagVersionCommand) (*domain.TagVersion, error)
	DeprecateVersion(context.Context, TransitionTagVersionCommand) (*domain.TagVersion, error)
	RetireVersion(context.Context, TransitionTagVersionCommand) (*domain.TagVersion, error)
}
