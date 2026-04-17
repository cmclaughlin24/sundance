package ports

import (
	"github.com/cmclaughlin24/sundance/forms/internal/core/domain"
)

type baseFormCommand struct {
	TenantID    string `validate:"required"`
	Name        string `validate:"required,max=75"`
	Description string `validate:"required,max=250"`
}

type CreateFormCommand struct {
	baseFormCommand
}

func NewCreateFormCommand(tenantID, name, description string) *CreateFormCommand {
	return &CreateFormCommand{
		baseFormCommand: baseFormCommand{
			TenantID:    tenantID,
			Name:        name,
			Description: description,
		},
	}
}

type UpdateFormCommand struct {
	baseFormCommand
	ID domain.FormID `validate:"required"`
}

func NewUpdateFormCommand(id domain.FormID, tenantID, name, description string) *UpdateFormCommand {
	return &UpdateFormCommand{
		baseFormCommand: baseFormCommand{
			TenantID:    tenantID,
			Name:        name,
			Description: description,
		},
		ID: id,
	}
}

type baseVersionCommand struct {
	FormID   domain.FormID `validate:"required"`
	TenantID string        `validate:"required"`
}

type CreateVersionCommand struct {
	baseVersionCommand
}

func NewCreateVersionCommand(formId domain.FormID, tenantID string) *CreateVersionCommand {
	return &CreateVersionCommand{
		baseVersionCommand{
			FormID:   formId,
			TenantID: tenantID,
		},
	}
}

type UpdateVersionCommand struct {
	baseVersionCommand
	VersionID domain.VersionID `validate:"required"`
	Pages     []*domain.Page
}

func NewUpdateVersionCommand(id domain.VersionID, formId domain.FormID, tenantID string, pages []*domain.Page) *UpdateVersionCommand {
	return &UpdateVersionCommand{
		baseVersionCommand: baseVersionCommand{
			FormID:   formId,
			TenantID: tenantID,
		},
		VersionID: id,
		Pages:     pages,
	}
}

type PublishVersionCommand struct {
	baseVersionCommand
	VersionID domain.VersionID `validate:"required"`
	UserID    string           `validate:"required"`
}

func NewPublishVersionCommand(formID domain.FormID, tenantID string, versionID domain.VersionID, userID string) *PublishVersionCommand {
	return &PublishVersionCommand{
		baseVersionCommand: baseVersionCommand{
			FormID:   formID,
			TenantID: tenantID,
		},
		VersionID: versionID,
		UserID:    userID,
	}
}

type RetireVersionCommand struct {
	baseVersionCommand
	VersionID domain.VersionID `validate:"required"`
	UserID    string           `validate:"required"`
}

func NewRetireVersionCommand(formID domain.FormID, tenantID string, versionID domain.VersionID, userID string) *RetireVersionCommand {
	return &RetireVersionCommand{
		baseVersionCommand: baseVersionCommand{
			FormID:   formID,
			TenantID: tenantID,
		},
		VersionID: versionID,
		UserID:    userID,
	}
}
