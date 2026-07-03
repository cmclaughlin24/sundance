package ports

import (
	"context"

	"sundance/backend/pkg/common/stratreg"
	"sundance/backend/pkg/database"
	"sundance/backend/services/forms/internal/core/domain"
)

type Repository struct {
	Database     database.Database
	Outbox       OutboxRepository
	Tags         TagsRepository
	TagVersions  TagVersionsRepository
	Forms        FormsRepository
	FormVersions FormVersionRepository
	Submissions  SubmissionsRepository
}

type FormsRepository interface {
	Find(context.Context, *FormFilters) ([]*domain.Form, error)
	FindByID(context.Context, domain.FormID) (*domain.Form, error)
	Upsert(context.Context, *domain.Form) (*domain.Form, error)
	Delete(context.Context, domain.FormID) error
}

type FormVersionRepository interface {
	Find(context.Context, domain.FormID) ([]*domain.FormVersion, error)
	FindByID(context.Context, domain.FormVersionID) (*domain.FormVersion, error)
	FindNextVersionNumber(context.Context, domain.FormID) (int, error)
	Upsert(context.Context, *domain.FormVersion) (*domain.FormVersion, error)
}

type SubmissionsRepository interface {
	Find(context.Context, *FindSubmissionsFilter) ([]*domain.Submission, error)
	FindJobs(context.Context, *FindSubmissionsFilter) ([]domain.SubmissionID, error)
	FindByID(context.Context, domain.SubmissionID) (*domain.Submission, error)
	FindByReferenceID(context.Context, domain.ReferenceID) (*domain.Submission, error)
	FindByIdempotencyID(context.Context, domain.IdempotencyID) (*domain.Submission, error)
	Upsert(context.Context, *domain.Submission) (*domain.Submission, error)
}

type TagsRepository interface {
	Find(context.Context, TagFilters) ([]*domain.Tag, error)
	FindByIDs(context.Context, []domain.TagID) ([]*domain.Tag, error)
	FindByID(context.Context, domain.TagID) (*domain.Tag, error)
	Upsert(context.Context, *domain.Tag) (*domain.Tag, error)
	Delete(context.Context, domain.TagID) error
}

type TagVersionsRepository interface {
	Find(context.Context, TagVersionFilters) ([]*domain.TagVersion, error)
	FindByIDs(context.Context, []domain.TagVersionID) ([]*domain.TagVersion, error)
	FindByID(context.Context, domain.TagVersionID) (*domain.TagVersion, error)
	FindNextVersionNumber(context.Context, domain.TagID) (int, error)
	Upsert(context.Context, *domain.TagVersion) (*domain.TagVersion, error)
}

type OutboxRepository interface {
	Find(context.Context, FindEventsFilter) ([]*domain.Event, error)
	Upsert(context.Context, *domain.Event) (*domain.Event, error)
}

type Publisher interface {
	Publish(context.Context, domain.Event) error
}

type Strategies struct {
	FieldValidator FieldValidatorRegistry
}

type FieldValidatorStrategy interface {
	Validate(context.Context, domain.Field, domain.SubmissionFieldValue) error
}

type FieldValidatorRegistry = stratreg.StrategyRegistry[domain.FieldType, FieldValidatorStrategy]

type RuleEvaluationContext = map[string]any

type RuleEvaluator interface {
	Evaluate(context.Context, *domain.Rule, RuleEvaluationContext) (bool, error)
}
