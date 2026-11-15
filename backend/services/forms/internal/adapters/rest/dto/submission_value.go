package dto

import "sundance/backend/services/forms/internal/core/domain"

type SubmissionValueDto struct {
	ElementID domain.ElementID `json:"elementId" validate:"required"`
	Value     any              `json:"value" validate:"required"`
	ItemIndex *int             `json:"itemIndex,omitempty"`
}
