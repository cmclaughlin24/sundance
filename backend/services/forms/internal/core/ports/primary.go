package ports

import (
	"context"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports/commands"
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
	Create(context.Context, commands.CreateFormCommand) (*domain.Form, error)
	Update(context.Context, commands.UpdateFormCommand) (*domain.Form, error)
	Delete(context.Context, commands.DeleteCommand[domain.FormID]) error
	FindVersions(context.Context, FindFormVersionsQuery) ([]*domain.FormVersion, error)
	FindVersion(context.Context, FindFormVersionByIDQuery) (*domain.FormVersion, error)
	CreateVersion(context.Context, *commands.CreateFormVersionCommand) (*domain.FormVersion, error)
	UpdateVersion(context.Context, *commands.UpdateFormVersionCommand) (*domain.FormVersion, error)
	PublishVersion(context.Context, commands.PublishFormVersionCommand) (*domain.FormVersion, error)
	RetireVersion(context.Context, commands.RetireFormVersionCommand) (*domain.FormVersion, error)
}

type SubmissionsAPI interface {
	Find(context.Context, FindSubmissionsQuery) ([]*domain.Submission, error)
	FindByID(context.Context, FindSubmissionByIDQuery[domain.SubmissionID]) (*domain.Submission, error)
	FindByReferenceID(context.Context, FindSubmissionByIDQuery[domain.ReferenceID]) (*domain.Submission, error)
	Create(context.Context, *commands.CreateSubmissionCommand) (*domain.Submission, error)
	Normalize(context.Context, *commands.NormalizeSubmissionCommand) (domain.FactMap, error)
	Replay(context.Context, commands.ReplaySubmissionCommand) error
}

type SubmissionJobsAPI interface {
	Find(context.Context, FindSubmissionJobsQuery) ([]domain.SubmissionID, error)
	Process(context.Context, domain.SubmissionID) error
}

type TagsAPI interface {
	Find(context.Context, FindTagsQuery) ([]*domain.Tag, error)
	FindById(context.Context, FindByIDQuery[domain.TagID]) (*domain.Tag, error)
	Create(context.Context, commands.CreateTagCommand) (*domain.Tag, error)
	Update(context.Context, commands.UpdateTagCommand) (*domain.Tag, error)
	Delete(context.Context, commands.DeleteCommand[domain.TagID]) error
	FindVersions(context.Context, FindTagVersionsQuery) ([]*domain.TagVersion, error)
	FindVersion(context.Context, FindTagVersionQuery) (*domain.TagVersion, error)
	CreateVersion(context.Context, commands.CreateTagVersionCommand) (*domain.TagVersion, error)
	PublishVersion(context.Context, commands.TransitionTagVersionCommand) (*domain.TagVersion, error)
	DeprecateVersion(context.Context, commands.TransitionTagVersionCommand) (*domain.TagVersion, error)
	RetireVersion(context.Context, commands.TransitionTagVersionCommand) (*domain.TagVersion, error)
}
