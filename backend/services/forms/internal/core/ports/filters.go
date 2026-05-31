package ports

import "sundance/backend/services/forms/internal/core/domain"

type TagFilters struct {
	TenantID string
}

type FormFilters struct {
	TenantID string
}

type FindSubmissionsFilter struct {
	TenantID string
	Statuses []domain.SubmissionStatus
	Take     int
}
