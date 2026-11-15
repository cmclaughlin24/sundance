package dto

import (
	"time"

	"sundance/backend/services/forms/internal/core/domain"
)

type SubmissionRequest struct {
	FormID    string                 `json:"formId" validate:"required,uuidv7"`
	VersionID string                 `json:"versionId" validate:"required,uuidv7"`
	Sections  []SubmissionSectionDto `json:"sections" validate:"dive"`
}

type SubmissionResponse struct {
	ID          domain.SubmissionID     `json:"id"`
	TenantID    string                  `json:"tenantId"`
	FormID      domain.FormID           `json:"formId"`
	VersionID   domain.FormVersionID    `json:"versionId"`
	ReferenceID domain.ReferenceID      `json:"referenceId"`
	Status      domain.SubmissionStatus `json:"status"`
	Sections    []SubmissionSectionDto  `json:"sections"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"updatedAt"`
}

func SubmissionToResponse(s *domain.Submission) *SubmissionResponse {
	sections := make([]SubmissionSectionDto, 0, len(s.Sections))
	for _, section := range s.Sections {
		sections = append(
			sections,
			toSubmissionSectionResponse(section),
		)
	}

	return &SubmissionResponse{
		ID:          s.ID,
		TenantID:    s.TenantID,
		FormID:      s.FormID,
		VersionID:   s.VersionID,
		ReferenceID: s.ReferenceID,
		Status:      s.Status,
		Sections:    sections,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
