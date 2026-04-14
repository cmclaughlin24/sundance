package ports

import (
	"github.com/cmclaughlin24/sundance/forms/internal/core/domain"
	"github.com/go-playground/validator/v10"
)

type FindByIDQuery struct {
	FormID   domain.FormID
	TenantID string
}

func NewFindByIDQuery(formID domain.FormID, tenantID string) (*FindByIDQuery, error) {
	query := &FindByIDQuery{
		FormID:   formID,
		TenantID: tenantID,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(query); err != nil {
		return nil, err
	}

	return query, nil
}

type FindVersionsQuery struct {
	FormID   domain.FormID
	TenantID string
}

func NewFindVersionsQuery(formID domain.FormID, tenantID string) (*FindVersionsQuery, error) {
	query := &FindVersionsQuery{
		FormID:   formID,
		TenantID: tenantID,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(query); err != nil {
		return nil, err
	}

	return query, nil
}

type FindVersionByIDQuery struct {
	FormID    domain.FormID
	VersionID domain.VersionID
	TenantID  string
}

func NewFindVersionByIDQuery(formID domain.FormID, tenantID string, versionID domain.VersionID) (*FindVersionByIDQuery, error) {
	query := &FindVersionByIDQuery{
		FormID:    formID,
		VersionID: versionID,
		TenantID:  tenantID,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(query); err != nil {
		return nil, err
	}

	return query, nil
}
