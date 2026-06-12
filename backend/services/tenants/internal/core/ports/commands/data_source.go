package commands

import (
	"sundance/backend/pkg/common/validate"
	"sundance/backend/services/tenants/internal/core/domain"
)

type baseDataSourceCommand struct {
	TenantID    domain.TenantID `validate:"required"`
	Name        string          `validate:"required"`
	Description string
	Type        domain.DataSourceType       `validate:"required"`
	Attributes  domain.DataSourceAttributes `validate:"required"`
}

type CreateDataSourceCommand struct {
	baseDataSourceCommand
}

func NewCreateDataSourceCommand(
	tenantID domain.TenantID,
	name,
	description string,
	sourceType domain.DataSourceType,
	attr domain.DataSourceAttributes,
) *CreateDataSourceCommand {
	return &CreateDataSourceCommand{
		baseDataSourceCommand: baseDataSourceCommand{
			TenantID:    tenantID,
			Name:        name,
			Description: description,
			Type:        sourceType,
			Attributes:  attr,
		},
	}
}

func (c *CreateDataSourceCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type UpdateDataSourceCommand struct {
	baseDataSourceCommand
	ID domain.DataSourceID `validate:"required"`
}

func NewUpdateDataSourceCommand(
	tenantID domain.TenantID,
	id domain.DataSourceID,
	name,
	description string,
	sourceType domain.DataSourceType,
	attr domain.DataSourceAttributes,
) *UpdateDataSourceCommand {
	return &UpdateDataSourceCommand{
		ID: id,
		baseDataSourceCommand: baseDataSourceCommand{
			TenantID:    tenantID,
			Name:        name,
			Description: description,
			Type:        sourceType,
			Attributes:  attr,
		},
	}
}

func (c *UpdateDataSourceCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type RemoveDataSourceCommand struct {
	ID       domain.DataSourceID `validate:"required"`
	TenantID domain.TenantID     `validate:"required"`
}

func NewRemoveDataSourceCommand(tenantID domain.TenantID, sourceID domain.DataSourceID) *RemoveDataSourceCommand {
	return &RemoveDataSourceCommand{
		TenantID: tenantID,
		ID:       sourceID,
	}
}

func (c *RemoveDataSourceCommand) Validate() error {
	return validate.ValidateStruct(c)
}

type ProcessDataSourceJobCommand struct {
	DataSource *domain.DataSource `validate:"required"`
}

func NewProcessDataSourceJobCommand(ds *domain.DataSource) *ProcessDataSourceJobCommand {
	return &ProcessDataSourceJobCommand{ds}
}

func (c *ProcessDataSourceJobCommand) Validate() error {
	return validate.ValidateStruct(c)
}
