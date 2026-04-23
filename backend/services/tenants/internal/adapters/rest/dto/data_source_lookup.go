package dto

import (
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
)

type DataSourceLookupResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func DataSourceLookupToResponse(lookup *domain.DataSourceLookup) *DataSourceLookupResponse {
	return &DataSourceLookupResponse{
		Code:        lookup.Code,
		Description: lookup.Description,
	}
}
