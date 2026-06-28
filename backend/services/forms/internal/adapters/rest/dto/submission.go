package dto

import (
	"time"

	"sundance/backend/services/forms/internal/core/domain"
)

type SubmissionRequest struct {
	FormID    string                    `json:"formId" validate:"required,uuidv7"`
	VersionID string                    `json:"versionId" validate:"required,uuidv7"`
	Values    []SubmissionFieldValueDto `json:"values" validate:"dive"`
}

type SubmissionFieldValueDto struct {
	FieldID         domain.FieldID `json:"fieldId" validate:"required"`
	Value           any            `json:"value" validate:"required"`
	CollectionIndex *int           `json:"collectionIndex,omitempty"`
}

type SubmissionResponse struct {
	ID          domain.SubmissionID       `json:"id"`
	TenantID    string                    `json:"tenantId"`
	FormID      domain.FormID             `json:"formId"`
	VersionID   domain.FormVersionID      `json:"versionId"`
	ReferenceID domain.ReferenceID        `json:"referenceId"`
	Status      domain.SubmissionStatus   `json:"status"`
	Values      []SubmissionFieldValueDto `json:"values"`
	CreatedAt   time.Time                 `json:"createdAt"`
	UpdatedAt   time.Time                 `json:"updatedAt"`
}

func SubmissionToResponse(s *domain.Submission) *SubmissionResponse {
	values := make([]SubmissionFieldValueDto, 0, len(s.Values))
	for _, value := range s.Values {
		values = append(values, SubmissionFieldValueDto{FieldID: value.FieldID, Value: value.Value, CollectionIndex: value.CollectionIndex})
	}

	return &SubmissionResponse{
		ID:          s.ID,
		TenantID:    s.TenantID,
		FormID:      s.FormID,
		VersionID:   s.VersionID,
		ReferenceID: s.ReferenceID,
		Status:      s.Status,
		Values:      values,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
