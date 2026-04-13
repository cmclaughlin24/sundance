package rest

import (
	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
)

type upsertTenantDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type upsertDataSourceDto struct {
	Type       domain.DataSourceType `json:"type"`
	Attributes any                   `json:"attributes"`
}

