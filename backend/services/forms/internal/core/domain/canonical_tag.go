package domain

import (
	"sundance/backend/pkg/common/validate"
	"time"
)

type CanonicalTagID string

type CanonicalTagStatus string

type CanonicalTag struct {
	ID          CanonicalTagID
	TenantID    string
	Key         string
	DisplayName string
	Status      CanonicalTagStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewCanonicalTag(tenantID, key, displayName string, status CanonicalTagStatus) (*CanonicalTag, error) {
	ct := &CanonicalTag{
		ID:          CanonicalTagID(NewID()),
		TenantID:    tenantID,
		Key:         key,
		DisplayName: displayName,
		Status:      status,
		CreatedAt:   Now(),
	}

	if err := validate.ValidateStruct(ct); err != nil {
		return nil, err
	}

	return ct, nil
}

func HydrateCanonicalTag(id CanonicalTagID, tenantID, key, displayName string, status CanonicalTagStatus, createdAt, updatedAt time.Time) *CanonicalTag {
	return &CanonicalTag{
		ID:          id,
		TenantID:    tenantID,
		Key:         key,
		DisplayName: displayName,
		Status:      status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
