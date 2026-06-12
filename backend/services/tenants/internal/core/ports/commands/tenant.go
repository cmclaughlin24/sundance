package commands

import (
	"sundance/backend/pkg/common/validate"
	"sundance/backend/services/tenants/internal/core/domain"
)


type baseTenantCommand struct {
	Name        string `validate:"required"`
	Description string
}

type CreateTenantCommand struct {
	baseTenantCommand
}

func NewCreateTenantCommand(name, description string) *CreateTenantCommand {
	return &CreateTenantCommand{
		baseTenantCommand{
			Name:        name,
			Description: description,
		},
	}
}

func (c *CreateTenantCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type UpdateTenantCommand struct {
	ID domain.TenantID `validate:"required"`
	baseTenantCommand
}

func NewUpdateTenantCommand(id domain.TenantID, name, description string) *UpdateTenantCommand {
	return &UpdateTenantCommand{
		id,
		baseTenantCommand{
			Name:        name,
			Description: description,
		},
	}
}

func (c *UpdateTenantCommand) Validate() error {
	return validate.ValidateStruct(c)
}
