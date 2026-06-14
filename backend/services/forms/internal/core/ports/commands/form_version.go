package commands

import (
	"sundance/backend/pkg/common/validate"
	"sundance/backend/services/forms/internal/core/domain"
)

type baseFormVersionCommand struct {
	TenantID string        `validate:"required"`
	FormID   domain.FormID `validate:"required"`
}

type CreateFormVersionCommand struct {
	baseFormVersionCommand
	PagesData []PageData
}

func NewCreateFormVersionCommand(tenantID string, formID domain.FormID, pages []PageData) *CreateFormVersionCommand {
	return &CreateFormVersionCommand{
		baseFormVersionCommand: baseFormVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		PagesData: pages,
	}
}

func (c *CreateFormVersionCommand) Validate() error {
	return validate.ValidateStruct(*c)
}

type UpdateFormVersionCommand struct {
	baseFormVersionCommand
	VersionID domain.FormVersionID `validate:"required"`
	PagesData     []PageData
}

func NewUpdateFormVersionCommand(tenantID string, id domain.FormVersionID, formID domain.FormID, pages []PageData) *UpdateFormVersionCommand {
	return &UpdateFormVersionCommand{
		baseFormVersionCommand: baseFormVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		VersionID: id,
		PagesData:     pages,
	}
}

func (c *UpdateFormVersionCommand) Validate() error {
	return validate.ValidateStruct(*c)
}

type PublishFormVersionCommand struct {
	baseFormVersionCommand
	VersionID domain.FormVersionID `validate:"required"`
	UserID    string               `validate:"required"`
}

func NewPublishFormVersionCommand(tenantID string, formID domain.FormID, versionID domain.FormVersionID, userID string) PublishFormVersionCommand {
	return PublishFormVersionCommand{
		baseFormVersionCommand: baseFormVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		VersionID: versionID,
		UserID:    userID,
	}
}

func (c PublishFormVersionCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type RetireFormVersionCommand struct {
	baseFormVersionCommand
	VersionID domain.FormVersionID `validate:"required"`
	UserID    string               `validate:"required"`
}

func NewRetireFormVersionCommand(tenantID string, formID domain.FormID, versionID domain.FormVersionID, userID string) RetireFormVersionCommand {
	return RetireFormVersionCommand{
		baseFormVersionCommand: baseFormVersionCommand{
			TenantID: tenantID,
			FormID:   formID,
		},
		VersionID: versionID,
		UserID:    userID,
	}
}

func (c RetireFormVersionCommand) Validate() error {
	return validate.ValidateStruct(c)
}
