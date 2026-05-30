package ports

import (
	"context"

	"sundance/backend/pkg/common/stratreg"
	"sundance/backend/pkg/database"
	"sundance/backend/services/forms/internal/core/domain"
)

type Repository struct {
	Database      database.Database
	CanonicalTags CanonicalTagRepository
	Forms         FormsRepository
	FormVersions  FormVersionRepository
	Submissions   SubmissionsRepository
}

type CanonicalTagRepository interface{
	Find(context.Context, CanonicalTagFilters) ([]*domain.CanonicalTag, error)
	FindByID(context.Context, domain.CanonicalTagID) (*domain.CanonicalTag, error)
	Upsert(context.Context, *domain.CanonicalTag) (*domain.CanonicalTag, error)
	Delete(context.Context, domain.CanonicalTagID) error
}

type FormsRepository interface {
	Find(context.Context, *FormFilters) ([]*domain.Form, error)
	FindByID(context.Context, domain.FormID) (*domain.Form, error)
	Upsert(context.Context, *domain.Form) (*domain.Form, error)
	Delete(context.Context, domain.FormID) error
}

type FormVersionRepository interface {
	Find(context.Context, domain.FormID) ([]*domain.FormVersion, error)
	FindByID(context.Context, domain.FormID, domain.FormVersionID) (*domain.FormVersion, error)
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
