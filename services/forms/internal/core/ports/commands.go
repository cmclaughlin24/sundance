package ports

import (
	"github.com/cmclaughlin24/sundance/forms/internal/core/domain"
	"github.com/go-playground/validator/v10"
)

type CreateFormCommand struct {
	TenantID    string
	Name        string
	Description string
}

func NewCreateFormCommand(tenantId, name, description string) (*CreateFormCommand, error) {
	command := &CreateFormCommand{
		TenantID:    tenantId,
		Name:        name,
		Description: description,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type UpdateFormCommand struct {
	ID          domain.FormID
	TenantID    string
	Name        string
	Description string
}

func NewUpdateFormCommand(id domain.FormID, tenantId, name, description string) (*UpdateFormCommand, error) {
	command := &UpdateFormCommand{
		ID:          id,
		TenantID:    tenantId,
		Name:        name,
		Description: description,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type CreateVersionCommand struct {
	ID       domain.VersionID
	FormID   domain.FormID
	TenantID string
}

func NewCreateVersionCommand(id domain.VersionID, formId domain.FormID, tenantID string) (*CreateVersionCommand, error) {
	command := &CreateVersionCommand{
		ID:       id,
		FormID:   formId,
		TenantID: tenantID,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type UpdateVersionCommand struct {
	ID       domain.VersionID
	FormID   domain.FormID
	TenantID string
	Pages    []*domain.Page
}

func NewUpdateVersionCommand(id domain.VersionID, formId domain.FormID, tenantID string, pages []*domain.Page) (*UpdateVersionCommand, error) {
	command := &UpdateVersionCommand{
		ID:       id,
		FormID:   formId,
		TenantID: tenantID,
		Pages:    pages,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type PublishVersionCommand struct {
	FormID    domain.FormID
	TenantID  string
	VersionID domain.VersionID
	UserId    string
}

func NewPublishVersionCommand(formId domain.FormID, tenantId string, versionId domain.VersionID, userId string) (*PublishVersionCommand, error) {
	command := &PublishVersionCommand{
		FormID:    formId,
		TenantID:  tenantId,
		VersionID: versionId,
		UserId:    userId,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type RetireVersionCommand struct {
	FormID    domain.FormID
	TenantID  string
	VersionID domain.VersionID
	UserId    string
}

func NewRetireVersionCommand(formId domain.FormID, tenantId string, versionId domain.VersionID, userId string) (*RetireVersionCommand, error) {
	command := &RetireVersionCommand{
		FormID:    formId,
		TenantID:  tenantId,
		VersionID: versionId,
		UserId:    userId,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}
