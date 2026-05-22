package ports

import (
	"sundance/backend/pkg/common/validate"
	"sundance/backend/services/forms/internal/core/domain"
)

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

type FindFormsByIDQuery struct {
	TenantID string        `validate:"required"`
	FormID   domain.FormID `validate:"required"`
}

func NewFindFormsByIDQuery(tenantID string, formID domain.FormID) *FindFormsByIDQuery {
	return &FindFormsByIDQuery{
		TenantID: tenantID,
		FormID:   formID,
	}
}

func (q *FindFormsByIDQuery) Validate() error {
	return validate.ValidateStruct(q)
}

type FindVersionsQuery struct {
	FindFormsByIDQuery
}

func NewFindVersionsQuery(tenantID string, formID domain.FormID) *FindVersionsQuery {
	return &FindVersionsQuery{
		FindFormsByIDQuery{
			TenantID: tenantID,
			FormID:   formID,
		},
	}
}

func (q *FindVersionsQuery) Validate() error {
	return validate.ValidateStruct(q)
}

type FindVersionByIDQuery struct {
	FindFormsByIDQuery
	VersionID domain.VersionID `validate:"required"`
}

func NewFindVersionByIDQuery(tenantID string, formID domain.FormID, versionID domain.VersionID) *FindVersionByIDQuery {
	return &FindVersionByIDQuery{
		FindFormsByIDQuery: FindFormsByIDQuery{
			TenantID: tenantID,
			FormID:   formID,
		},
		VersionID: versionID,
	}
}

func (q *FindVersionByIDQuery) Validate() error {
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

type FindSubmissionByIDQuery[T any] struct {
	TenantID string `validate:"required"`
	ID       T      `validate:"required"`
}

func NewFindSubmissionByIDQuery[T any](tenantID string, id T) *FindSubmissionByIDQuery[T] {
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

