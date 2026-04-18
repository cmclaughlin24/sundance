package dto

import (
	"time"

	"github.com/cmclaughlin24/sundance/submissions/internal/core/domain"
)

type SubmissionResponse struct {
	ID          domain.SubmissionID     `json:"id"`
	TenantID    string                  `json:"tenantId"`
	FormID      string                  `json:"formId"`
	VersionID   string                  `json:"versionId"`
	ReferenceID domain.ReferenceID      `json:"referenceId"`
	Status      domain.SubmissionStatus `json:"status"`
	Payload     any                     `json:"payload"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"createdBy"`
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
