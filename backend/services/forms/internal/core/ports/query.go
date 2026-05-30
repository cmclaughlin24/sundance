package ports

import (
	"sundance/backend/pkg/common/validate"
	"sundance/backend/services/forms/internal/core/domain"
)

type FindByIDQuery[T comparable] struct {
	TenantID string `validate:"required"`
	ID       T      `validate:"required"`
}

func NewFindByIDQuery[T comparable](tenantID string, id T) FindByIDQuery[T] {
	return FindByIDQuery[T]{
		TenantID: tenantID,
		ID:       id,
	}
}

func (q FindByIDQuery[T]) Validate() error {
	return validate.ValidateStruct(q)
}

type FindCanonicalTagsQuery struct {
	TenantID string `validate:"required"`
}

func NewCanonicalTagsQuery(tenantID string) FindCanonicalTagsQuery {
	return FindCanonicalTagsQuery{
		TenantID: tenantID,
	}
}

func (q FindCanonicalTagsQuery) Validate() error {
	return validate.ValidateStruct(q)
}

type FindFormsQuery struct {
	// TODO: Add pagination support through embedded struct.
	TenantID string `validate:"required"`
}

func NewFindFormsQuery(tenantID string) *FindFormsQuery {
	return &FindFormsQuery{
		TenantID: tenantID,
	}
}

func (q *FindFormsQuery) Validate() error {
	return validate.ValidateStruct(q)
}

type FindFormByIDQuery = FindByIDQuery[domain.FormID]

func NewFindFormByIDQuery(tenantID string, formID domain.FormID) *FindFormByIDQuery {
	return &FindFormByIDQuery{
		TenantID: tenantID,
		ID:   formID,
	}
}

type FindFormVersionsQuery struct {
	FindFormByIDQuery
}

func NewFindFormVersionsQuery(tenantID string, formID domain.FormID) *FindFormVersionsQuery {
	return &FindFormVersionsQuery{
		FindFormByIDQuery{
			TenantID: tenantID,
			ID:   formID,
		},
	}
}

func (q *FindFormVersionsQuery) Validate() error {
	return validate.ValidateStruct(q)
}

type FindFormVersionByIDQuery struct {
	FindFormByIDQuery
	VersionID domain.FormVersionID `validate:"required"`
}

func NewFindFormVersionByIDQuery(tenantID string, formID domain.FormID, versionID domain.FormVersionID) *FindFormVersionByIDQuery {
	return &FindFormVersionByIDQuery{
		FindFormByIDQuery: FindFormByIDQuery{
			TenantID: tenantID,
			ID:   formID,
		},
		VersionID: versionID,
	}
}

func (q *FindFormVersionByIDQuery) Validate() error {
	return validate.ValidateStruct(q)
}

type FindSubmissionsQuery struct {
	TenantID string `validate:"required"`
}

func NewFindSubmissionsQuery(tenantID string) *FindSubmissionsQuery {
	return &FindSubmissionsQuery{
		TenantID: tenantID,
	}
}

func (q *FindSubmissionsQuery) Validate() error {
	return validate.ValidateStruct(q)
}

type FindSubmissionByIDQuery[T comparable] = FindByIDQuery[T]

func NewFindSubmissionByIDQuery[T comparable](tenantID string, id T) *FindSubmissionByIDQuery[T] {
	query := &FindSubmissionByIDQuery[T]{
		TenantID: tenantID,
		ID:       id,
	}

	return query
}

type FindSubmissionJobsQuery struct {
	Take int `validate:"min=0"`
}

func NewFindSubmissionJobsQuery(take int) *FindSubmissionJobsQuery {
	return &FindSubmissionJobsQuery{take}
}

func (q *FindSubmissionJobsQuery) Validate() error {
	return validate.ValidateStruct(q)
}
