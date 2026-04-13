package ports

import "github.com/cmclaughlin24/sundance/forms/internal/core/domain"

type FindByIDQuery struct {
	FormID   domain.FormID
	TenantID string
}

func NewFindByIDQuery(formID domain.FormID, tenantID string) *FindByIDQuery {
	return &FindByIDQuery{
		FormID:   formID,
		TenantID: tenantID,
	}
}

type FindVersionsQuery struct {
	FormID   domain.FormID
	TenantID string
}

func NewFindVersionsQuery(formID domain.FormID, tenantID string) *FindVersionsQuery {
	return &FindVersionsQuery{
		FormID:   formID,
		TenantID: tenantID,
	}
}

type FindVersionByIDQuery struct {
	FormID    domain.FormID
	VersionID domain.VersionID
	TenantID  string
}

func NewFindVersionByIDQuery(formID domain.FormID, tenantID string, versionID domain.VersionID) *FindVersionByIDQuery {
	return &FindVersionByIDQuery{
		FormID:    formID,
		VersionID: versionID,
		TenantID:  tenantID,
	}
}
