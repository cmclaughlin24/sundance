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

func NewCreateFormCommand(tenantId, name, description string) *CreateFormCommand {
	return &CreateFormCommand{
		TenantID:    tenantId,
		Name:        name,
		Description: description,
	}
}

type UpdateFormCommand struct {
	ID          domain.FormID
	TenantID    string
	Name        string
	Description string
}

func NewUpdateFormCommand(id domain.FormID, tenantId, name, description string) *UpdateFormCommand {
	return &UpdateFormCommand{
		ID:          id,
		TenantID:    tenantId,
		Name:        name,
		Description: description,
	}
}

type baseVersionCommand struct {
	VersionID domain.VersionID
	FormID    domain.FormID
	TenantID  string
}

type CreateVersionCommand struct {
	FormID   domain.FormID
	TenantID string
}

func NewCreateVersionCommand(formId domain.FormID, tenantID string) *CreateVersionCommand {
	return &CreateVersionCommand{
		FormID:   formId,
		TenantID: tenantID,
	}
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

func NewPublishVersionCommand(formID domain.FormID, tenantID string, versionID domain.VersionID, userID string) *PublishVersionCommand {
	return &PublishVersionCommand{
		baseVersionCommand: baseVersionCommand{
			VersionID: versionID,
			FormID:    formID,
			TenantID:  tenantID,
		},
		UserID: userID,
	}
}

type RetireVersionCommand struct {
	baseVersionCommand
	UserID string
}

func NewRetireVersionCommand(formID domain.FormID, tenantID string, versionID domain.VersionID, userID string) *RetireVersionCommand {
	return &RetireVersionCommand{
		baseVersionCommand: baseVersionCommand{
			VersionID: versionID,
			FormID:    formID,
			TenantID:  tenantID,
		},
		UserID: userID,
	}
}
