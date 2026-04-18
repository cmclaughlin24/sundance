package ports

import (
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
)

type baseFindQuery struct {
	FormID   domain.FormID `validate:"required"`
	TenantID string        `validate:"required"`
}

type FindByIDQuery struct {
	baseFindQuery
}

func NewFindByIDQuery(formID domain.FormID, tenantID string) *FindByIDQuery {
	return &FindByIDQuery{
		baseFindQuery{
			FormID:   formID,
			TenantID: tenantID,
		},
	}
}

type FindVersionsQuery struct {
	baseFindQuery
}

func NewFindVersionsQuery(formID domain.FormID, tenantID string) *FindVersionsQuery {
	return &FindVersionsQuery{
		baseFindQuery{
			FormID:   formID,
			TenantID: tenantID,
		},
	}
}

type FindVersionByIDQuery struct {
	baseFindQuery
	VersionID domain.VersionID `validate:"required"`
}

func NewFindVersionByIDQuery(formID domain.FormID, tenantID string, versionID domain.VersionID) *FindVersionByIDQuery {
	return &FindVersionByIDQuery{
		baseFindQuery: baseFindQuery{
			FormID:   formID,
			TenantID: tenantID,
		},
		VersionID: versionID,
	}
}
