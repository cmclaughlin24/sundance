package ports

import "github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"

type FormFilters struct {
	TenantID string
}

type FindSubmissionsFilter struct {
	TenantID string
	Statuses []domain.SubmissionStatus
}
