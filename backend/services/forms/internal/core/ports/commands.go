package ports

import (
	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
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

func (c *CreateFormCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type UpdateFormCommand struct {
	baseFormCommand
	ID domain.FormID `validate:"required"`
}

func NewUpdateFormCommand(tenantID string, id domain.FormID, name, description string) *UpdateFormCommand {
	return &UpdateFormCommand{
		baseFormCommand: baseFormCommand{
			TenantID:    tenantID,
			Name:        name,
			Description: description,
		},
		ID: id,
	}
}

func (c *UpdateFormCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type RemoveFormCommand struct {
	ID       domain.FormID `validate:"required"`
	TenantID string        `validate:"required"`
}

func NewRemoveFormCommand(tenantID string, id domain.FormID) *RemoveFormCommand {
	return &RemoveFormCommand{
		TenantID: tenantID,
		ID:       id,
	}
}

func (c *RemoveFormCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type baseVersionCommand struct {
	TenantID string        `validate:"required"`
	FormID   domain.FormID `validate:"required"`
}

type CreateVersionCommand struct {
	baseVersionCommand
	Pages []*domain.Page
}

func NewCreateVersionCommand(tenantID string, formID domain.FormID, pages []*domain.Page) *CreateVersionCommand {
	return &CreateVersionCommand{
		baseVersionCommand: baseVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		Pages: pages,
	}
}

func (c *CreateVersionCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type UpdateVersionCommand struct {
	baseVersionCommand
	VersionID domain.VersionID `validate:"required"`
	Pages     []*domain.Page
}

func NewUpdateVersionCommand(tenantID string, id domain.VersionID, formID domain.FormID, pages []*domain.Page) *UpdateVersionCommand {
	return &UpdateVersionCommand{
		baseVersionCommand: baseVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		VersionID: id,
		Pages:     pages,
	}
}

func (c *UpdateVersionCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type PublishVersionCommand struct {
	baseVersionCommand
	VersionID domain.VersionID `validate:"required"`
	UserID    string           `validate:"required"`
}

func NewPublishVersionCommand(tenantID string, formID domain.FormID, versionID domain.VersionID, userID string) *PublishVersionCommand {
	return &PublishVersionCommand{
		baseVersionCommand: baseVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		VersionID: versionID,
		UserID:    userID,
	}
}

func (c *PublishVersionCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type RetireVersionCommand struct {
	baseVersionCommand
	VersionID domain.VersionID `validate:"required"`
	UserID    string           `validate:"required"`
}

func NewRetireVersionCommand(tenantID string, formID domain.FormID, versionID domain.VersionID, userID string) *RetireVersionCommand {
	return &RetireVersionCommand{
		baseVersionCommand: baseVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		VersionID: versionID,
		UserID:    userID,
	}
}

func (c *RetireVersionCommand) Validate() error {
	return validate.ValidateStruct(c)
}
