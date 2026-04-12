package rest

import "github.com/cmclaughlin24/sundance/tenants/internal/core/domain"

type tenantDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type dataSourceDto struct {
	Type       domain.DataSourceType
	Attributes any
}
