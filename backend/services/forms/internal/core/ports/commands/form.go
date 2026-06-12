package commands

import (
	"sundance/backend/pkg/common/validate"
	"sundance/backend/services/forms/internal/core/domain"
)

type baseFormCommand struct {
	TenantID    string `validate:"required"`
	Name        string `validate:"required,max=75"`
	Description string `validate:"required,max=500"`
}

type CreateFormCommand struct {
	baseFormCommand
}

func NewCreateFormCommand(tenantID, name, description string) CreateFormCommand {
	return CreateFormCommand{
		baseFormCommand: baseFormCommand{
			TenantID:    tenantID,
			Name:        name,
			Description: description,
		},
	}
}

func (c CreateFormCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type UpdateFormCommand struct {
	baseFormCommand
	ID domain.FormID `validate:"required"`
}

func NewUpdateFormCommand(tenantID string, id domain.FormID, name, description string) UpdateFormCommand {
	return UpdateFormCommand{
		baseFormCommand: baseFormCommand{
			TenantID:    tenantID,
			Name:        name,
			Description: description,
		},
		ID: id,
	}
}

func (c UpdateFormCommand) Validate() error {
	return validate.ValidateStruct(c)
}
