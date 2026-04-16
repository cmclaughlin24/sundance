package dto

import (
	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
)

type UpsertTenantRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpsertDataSourceRequest struct {
	Type       domain.DataSourceType `json:"type" validate:"required"`
	Attributes any                   `json:"attributes"`
}
