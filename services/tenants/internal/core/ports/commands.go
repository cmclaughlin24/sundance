package ports

import "github.com/cmclaughlin24/sundance/tenants/internal/core/domain"

type CreateTenantCommand struct {
	Name        string
	Description string
}

func NewCreateTenantCommand(name, description string) CreateTenantCommand {
	command := CreateTenantCommand{
		Name:        name,
		Description: description,
	}

	// TODO: Add validation.

	return command
}

type UpdateTenantCommand struct {
	ID          domain.TenantID
	Name        string
	Description string
}

func NewUpdateTenantCommand(id domain.TenantID, name, description string) UpdateTenantCommand {
	command := UpdateTenantCommand{
		ID:          id,
		Name:        name,
		Description: description,
	}

	// TODO: Add validation.

	return command
}

type CreateDataSourceCommand struct {
	TenantID   domain.TenantID
	Type       domain.DataSourceType
	Attributes domain.DataSourceAttributes
}

func NewCreateDataSourceCommand(tenantId domain.TenantID, sourceType domain.DataSourceType, attr domain.DataSourceAttributes) CreateDataSourceCommand {
	command := CreateDataSourceCommand{
		TenantID:   tenantId,
		Type:       sourceType,
		Attributes: attr,
	}

	// TODO: Add validation.

	return command
}

type UpdateDataSourceCommand struct {
	ID         domain.DataSourceID
	TenantID   domain.TenantID
	Type       domain.DataSourceType
	Attributes domain.DataSourceAttributes
}

func NewUpdateDataSourceCommand(id domain.DataSourceID, tenantId domain.TenantID, sourceType domain.DataSourceType, attr domain.DataSourceAttributes) UpdateDataSourceCommand {
	command := UpdateDataSourceCommand{
		ID:         id,
		TenantID:   tenantId,
		Type:       sourceType,
		Attributes: attr,
	}

	// TODO: Add validation.

	return command
}
