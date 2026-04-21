package ports

import (
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
)

type baseFindQuery struct {
	TenantID string        `validate:"required"`
	FormID   domain.FormID `validate:"required"`
}

type FindByIDQuery struct {
	baseFindQuery
}

func NewFindByIDQuery(tenantID string, formID domain.FormID) *FindByIDQuery {
	return &FindByIDQuery{
		baseFindQuery{
			TenantID: tenantID,
			FormID:   formID,
		},
	}
}

type FindVersionsQuery struct {
	baseFindQuery
}

func NewFindVersionsQuery(tenantID string, formID domain.FormID) *FindVersionsQuery {
	return &FindVersionsQuery{
		baseFindQuery{
			TenantID: tenantID,
			FormID:   formID,
		},
	}
}

type FindVersionByIDQuery struct {
	baseFindQuery
	VersionID domain.VersionID `validate:"required"`
}

func NewFindVersionByIDQuery(tenantID string, formID domain.FormID, versionID domain.VersionID) *FindVersionByIDQuery {
	return &FindVersionByIDQuery{
		baseFindQuery: baseFindQuery{
			TenantID: tenantID,
			FormID:   formID,
		},
		VersionID: versionID,
	}
}
