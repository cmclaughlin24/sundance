package rest

import (
	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
)

type upsertTenantRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type upsertDataSourceRequest struct {
	Type       domain.DataSourceType `json:"type"`
	Attributes any                   `json:"attributes"`
}

