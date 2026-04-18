package ports

import (
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
)

type baseTenantCommand struct {
	Name        string `validate:"required,max=75"`
	Description string `validate:"required,max=250"`
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

type baseDataSourceCommand struct {
	TenantID    domain.TenantID             `validate:"required"`
	Name        string                      `validate:"required,max=75"`
	Description string                      `validate:"required,max=250"`
	Type        domain.DataSourceType       `validate:"oneof=static scheduled query"`
	Attributes  domain.DataSourceAttributes `validate:"required"`
}

type CreateDataSourceCommand struct {
	baseDataSourceCommand
}

func NewCreateDataSourceCommand(
	tenantId domain.TenantID,
	name,
	description string,
	sourceType domain.DataSourceType,
	attr domain.DataSourceAttributes,
) *CreateDataSourceCommand {
	return &CreateDataSourceCommand{
		baseDataSourceCommand: baseDataSourceCommand{
			TenantID:    tenantId,
			Name:        name,
			Description: description,
			Type:        sourceType,
			Attributes:  attr,
		},
	}
}

type UpdateDataSourceCommand struct {
	baseDataSourceCommand
	ID domain.DataSourceID `validate:"required"`
}

func NewUpdateDataSourceCommand(
	id domain.DataSourceID,
	tenantId domain.TenantID,
	name,
	description string,
	sourceType domain.DataSourceType,
	attr domain.DataSourceAttributes,
) *UpdateDataSourceCommand {
	return &UpdateDataSourceCommand{
		ID: id,
		baseDataSourceCommand: baseDataSourceCommand{
			TenantID:   tenantId,
			Name:        name,
			Description: description,
			Type:       sourceType,
			Attributes: attr,
		},
	}
}
