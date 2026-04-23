package dto

import (
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
)

type LookupResponse struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func LookupToResponse(lookup *domain.Lookup) *LookupResponse {
	return &LookupResponse{
		Value: lookup.Value,
		Label: lookup.Label,
	}
}

func LookupsToResponse(lookups []*domain.Lookup) []*LookupResponse {
	dtos := make([]*LookupResponse, 0, len(lookups))

	for _, lookup := range lookups {
		dtos = append(dtos, LookupToResponse(lookup))
	}

	return dtos
}
