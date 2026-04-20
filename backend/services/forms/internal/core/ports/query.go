package ports

import (
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
)

type baseFindQuery struct {
	FormID domain.FormID `validate:"required"`
}

type FindByIDQuery struct {
	baseFindQuery
}

func NewFindByIDQuery(formID domain.FormID) *FindByIDQuery {
	return &FindByIDQuery{
		baseFindQuery{
			FormID: formID,
		},
	}
}

type FindVersionsQuery struct {
	baseFindQuery
}

func NewFindVersionsQuery(formID domain.FormID) *FindVersionsQuery {
	return &FindVersionsQuery{
		baseFindQuery{
			FormID: formID,
		},
	}
}

type FindVersionByIDQuery struct {
	baseFindQuery
	VersionID domain.VersionID `validate:"required"`
}

func NewFindVersionByIDQuery(formID domain.FormID, versionID domain.VersionID) *FindVersionByIDQuery {
	return &FindVersionByIDQuery{
		baseFindQuery: baseFindQuery{
			FormID: formID,
		},
		VersionID: versionID,
	}
}
