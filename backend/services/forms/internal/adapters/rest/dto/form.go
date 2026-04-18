package dto

import (
	"time"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
)

type UpsertFormRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FormResponse struct {
	ID          domain.FormID `json:"id"`
	TenantID    string        `json:"tenantId"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

func FormToResponse(form *domain.Form) *FormResponse {
	if form == nil {
		return nil
	}

	return &FormResponse{
		ID:          form.ID,
		TenantID:    form.TenantID,
		Name:        form.Name,
		Description: form.Description,
		CreatedAt:   form.CreatedAt,
		UpdatedAt:   form.UpdatedAt,
	}
}
