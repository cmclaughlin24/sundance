package ports

import (
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
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

type FormFilters struct {
	TenantID string
}
