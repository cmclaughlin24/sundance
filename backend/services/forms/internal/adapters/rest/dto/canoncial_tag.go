package dto

import "sundance/backend/services/forms/internal/core/domain"

type UpsertCanonicalTagRequest struct {}

type CanonicalTagResponse struct {}

func CanonicalTagToResponse(ct *domain.CanonicalTag) *CanonicalTagResponse {
	return &CanonicalTagResponse{}
}
