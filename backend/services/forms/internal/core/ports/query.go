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

type FindTagsQuery struct {
	TenantID string `validate:"required"`
}

func NewTagsQuery(tenantID string) FindTagsQuery {
	return FindTagsQuery{
		TenantID: tenantID,
	}
}

func (q FindTagsQuery) Validate() error {
	return validate.ValidateStruct(q)
}

type FindTagVersionsQuery struct {
	FindByIDQuery[domain.TagID]
}

func NewFindTagVersionsQuery(tenantID string, tagID domain.TagID) FindTagVersionsQuery {
	return FindTagVersionsQuery{
		FindByIDQuery[domain.TagID]{
			TenantID: tenantID,
			ID:       tagID,
		},
	}
}

func (q FindTagVersionsQuery) Validate() error {
	return validate.ValidateStruct(q)
}

type FindTagVersionQuery struct {
	VersionID domain.TagVersionID `validate:"required"`
	FindByIDQuery[domain.TagID]
}

func NewFindTagVersionQuery(tenantID string, tagID domain.TagID, versionID domain.TagVersionID) FindTagVersionQuery {
	return FindTagVersionQuery{
		versionID,
		FindByIDQuery[domain.TagID]{
			TenantID: tenantID,
			ID:       tagID,
		},
	}
}

func (q FindTagVersionQuery) Validate() error {
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

type FindFormVersionsQuery struct {
	TenantID string        `validate:"required"`
	ID       domain.FormID `validate:"required"`
}

func NewFindFormVersionsQuery(tenantID string, formID domain.FormID) *FindFormVersionsQuery {
	return &FindFormVersionsQuery{
		TenantID: tenantID,
		ID:       formID,
	}
}

func (q *FindFormVersionsQuery) Validate() error {
	return validate.ValidateStruct(q)
}

type FindFormVersionByIDQuery struct {
	TenantID  string               `validate:"required"`
	ID        domain.FormID        `validate:"required"`
	VersionID domain.FormVersionID `validate:"required"`
}

func NewFindFormVersionByIDQuery(tenantID string, formID domain.FormID, versionID domain.FormVersionID) *FindFormVersionByIDQuery {
	return &FindFormVersionByIDQuery{
		TenantID:  tenantID,
		ID:        formID,
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
