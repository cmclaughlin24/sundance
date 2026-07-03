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

type FindEventsFilter struct {
	Statuses     []domain.EventStatus
	RetryLimit   int
	CreatedAfter time.Time
	Take         int
}
