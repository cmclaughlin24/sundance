package ports

import (
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/validate"
)

type baseTenantCommand struct {
	Name        string `validate:"required,max=75"`
	Description string `validate:"required,max=250"`
}

type CreateTenantCommand struct {
	baseTenantCommand
}

func NewCreateTenantCommand(name, description string) (*CreateTenantCommand, error) {
	command := &CreateTenantCommand{
		baseTenantCommand{
			Name:        name,
			Description: description,
		},
	}

	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type UpdateTenantCommand struct {
	ID domain.TenantID `validate:"required"`
	baseTenantCommand
}

func NewUpdateTenantCommand(id domain.TenantID, name, description string) (*UpdateTenantCommand, error) {
	command := &UpdateTenantCommand{
		id,
		baseTenantCommand{
			Name:        name,
			Description: description,
		},
	}

	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type baseDataSourceCommand struct {
	TenantID   domain.TenantID             `validate:"required"`
	Type       domain.DataSourceType       `validate:"oneof=static scheduled query"`
	Attributes domain.DataSourceAttributes `validate:"required"`
}

type CreateDataSourceCommand struct {
	baseDataSourceCommand
}

func NewCreateDataSourceCommand(
	tenantId domain.TenantID,
	sourceType domain.DataSourceType,
	attr domain.DataSourceAttributes,
) (*CreateDataSourceCommand, error) {
	command := &CreateDataSourceCommand{
		baseDataSourceCommand: baseDataSourceCommand{
			TenantID:   tenantId,
			Type:       sourceType,
			Attributes: attr,
		},
	}

	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type UpdateDataSourceCommand struct {
	baseDataSourceCommand
	ID domain.DataSourceID `validate:"required"`
}

func NewUpdateDataSourceCommand(
	id domain.DataSourceID,
	tenantId domain.TenantID,
	sourceType domain.DataSourceType,
	attr domain.DataSourceAttributes,
) (*UpdateDataSourceCommand, error) {
	command := &UpdateDataSourceCommand{
		ID: id,
		baseDataSourceCommand: baseDataSourceCommand{
			TenantID:   tenantId,
			Type:       sourceType,
			Attributes: attr,
		},
	}

	if err := validate.ValidateStruct(command); err != nil {
		return nil, err
	}

	return command, nil
}
