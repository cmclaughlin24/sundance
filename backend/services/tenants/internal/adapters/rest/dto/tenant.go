package dto

import (
	"time"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
)

type TenantRequest struct {
	Name        string `json:"name" validate:"required,max=75"`
	Description string `json:"description" validate:"max=500"`
}

type TenantResponse struct {
	ID          domain.TenantID `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func TenantToResponse(tenant *domain.Tenant) *TenantResponse {
	return &TenantResponse{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Description: tenant.Description,
		CreatedAt:   tenant.CreatedAt,
		UpdatedAt:   tenant.UpdatedAt,
	}
}
