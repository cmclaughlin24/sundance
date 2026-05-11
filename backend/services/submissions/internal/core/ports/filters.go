package ports

import "github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"

type FindSubmissionsFilter struct {
	TenantID string
	Statuses []domain.SubmissionStatus
}
