package dto

import (
	"time"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
)

type SubmissionRequest struct {
	FormID    string         `json:"formId" validate:"required,uuidv7"`
	VersionID string         `json:"versionId" validate:"required,uuidv7"`
	Payload   map[string]any `json:"payload" validate:"required" swaggertype:"object"`
}

type SubmissionResponse struct {
	ID          domain.SubmissionID     `json:"id"`
	TenantID    string                  `json:"tenantId"`
	FormID      domain.FormID           `json:"formId"`
	VersionID   domain.VersionID        `json:"versionId"`
	ReferenceID domain.ReferenceID      `json:"referenceId"`
	Status      domain.SubmissionStatus `json:"status"`
	Payload     any                     `json:"payload" swaggertype:"object"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"updatedAt"`
}

func SubmissionToResponse(s *domain.Submission) *SubmissionResponse {
	return &SubmissionResponse{
		ID:          s.ID,
		TenantID:    s.TenantID,
		FormID:      s.FormID,
		VersionID:   s.VersionID,
		ReferenceID: s.ReferenceID,
		Status:      s.Status,
		Payload:     s.Payload,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
