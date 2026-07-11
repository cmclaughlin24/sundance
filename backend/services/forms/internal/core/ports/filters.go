package ports

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type FormFilters struct {
	TenantID string
}

type FindSubmissionsFilter struct {
	TenantID string
	Statuses []domain.SubmissionStatus
	Take     int
}

type TagFilters struct {
	TenantID string
}

type TagVersionFilters struct {
	TagID    domain.TagID
	Statuses []domain.TagStatus
}

type ClaimEventsOptions struct {
	RetryLimit    int
	BatchSize     int
	CreatedAfter  time.Time
	LeaseDuration time.Duration
}
