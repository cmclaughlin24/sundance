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

type baseVersionCommand struct {
	VersionID domain.VersionID
	FormID    domain.FormID
	TenantID  string
}

type CreateVersionCommand struct {
	baseVersionCommand
}

func NewCreateVersionCommand(id domain.VersionID, formId domain.FormID, tenantID string) (*CreateVersionCommand, error) {
	command := &CreateVersionCommand{
		baseVersionCommand{
			VersionID: id,
			FormID:    formId,
			TenantID:  tenantID,
		},
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type UpdateVersionCommand struct {
	baseVersionCommand
	Pages []*domain.Page
}

func NewUpdateVersionCommand(id domain.VersionID, formId domain.FormID, tenantID string, pages []*domain.Page) (*UpdateVersionCommand, error) {
	command := &UpdateVersionCommand{
		baseVersionCommand: baseVersionCommand{
			VersionID: id,
			FormID:    formId,
			TenantID:  tenantID,
		},
		Pages: pages,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type PublishVersionCommand struct {
	baseVersionCommand
	UserID string
}

func NewPublishVersionCommand(formID domain.FormID, tenantID string, versionID domain.VersionID, userID string) (*PublishVersionCommand, error) {
	command := &PublishVersionCommand{
		baseVersionCommand: baseVersionCommand{
			VersionID: versionID,
			FormID:    formID,
			TenantID:  tenantID,
		},
		UserID: userID,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type RetireVersionCommand struct {
	baseVersionCommand
	UserID string
}

func NewRetireVersionCommand(formID domain.FormID, tenantID string, versionID domain.VersionID, userID string) (*RetireVersionCommand, error) {
	command := &RetireVersionCommand{
		baseVersionCommand: baseVersionCommand{
			VersionID: versionID,
			FormID:    formID,
			TenantID:  tenantID,
		},
		UserID: userID,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}
