package ports

import "github.com/cmclaughlin24/sundance/forms/internal/core/domain"

type FindByIdQuery struct {
	FormID       domain.FormID
	TenantID string
}
