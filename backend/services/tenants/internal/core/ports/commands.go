package ports

import (
	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
	"github.com/go-playground/validator/v10"
)

type CreateTenantCommand struct {
	Name        string `validate:"required"`
	Description string `validate:"required"`
}

func NewCreateTenantCommand(name, description string) (*CreateTenantCommand, error) {
	command := &CreateTenantCommand{
		Name:        name,
		Description: description,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type UpdateTenantCommand struct {
	ID          domain.TenantID `validate:"required"`
	Name        string          `validate:"required"`
	Description string          `validate:"required"`
}

func NewUpdateTenantCommand(id domain.TenantID, name, description string) (*UpdateTenantCommand, error) {
	command := &UpdateTenantCommand{
		ID:          id,
		Name:        name,
		Description: description,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type CreateDataSourceCommand struct {
	TenantID   domain.TenantID
	Type       domain.DataSourceType
	Attributes domain.DataSourceAttributes
}

func NewCreateDataSourceCommand(
	tenantId domain.TenantID,
	sourceType domain.DataSourceType,
	attr domain.DataSourceAttributes,
) (*CreateDataSourceCommand, error) {
	command := &CreateDataSourceCommand{
		TenantID:   tenantId,
		Type:       sourceType,
		Attributes: attr,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}

type UpdateDataSourceCommand struct {
	ID         domain.DataSourceID
	TenantID   domain.TenantID
	Type       domain.DataSourceType
	Attributes domain.DataSourceAttributes
}

func NewUpdateDataSourceCommand(
	id domain.DataSourceID,
	tenantId domain.TenantID,
	sourceType domain.DataSourceType,
	attr domain.DataSourceAttributes,
) (*UpdateDataSourceCommand, error) {
	command := &UpdateDataSourceCommand{
		ID:         id,
		TenantID:   tenantId,
		Type:       sourceType,
		Attributes: attr,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(command); err != nil {
		return nil, err
	}

	return command, nil
}
