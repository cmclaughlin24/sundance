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

type CreateVersionCommand struct{}

func NewCreateVersionCommand() (*CreateVersionCommand, error) {
	command := &CreateVersionCommand{}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type UpdateVersionCommand struct {
}

func NewUpdateVersionCommand() (*UpdateVersionCommand, error) {
	command := &UpdateVersionCommand{}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type PublishVersionCommand struct {
	FormID    domain.FormID
	VersionID domain.VersionID
	UserId    string
}

func NewPublishVersionCommand(formId domain.FormID, versionId domain.VersionID, userId string) (*PublishVersionCommand, error) {
	command := &PublishVersionCommand{
		FormID:    formId,
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
	VersionID domain.VersionID
	UserId    string
}

func NewRetireVersionCommand(formId domain.FormID, versionId domain.VersionID, userId string) (*RetireVersionCommand, error) {
	command := &RetireVersionCommand{
		FormID:    formId,
		VersionID: versionId,
		UserId:    userId,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}
