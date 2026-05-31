package domain

import (
	"errors"
	"sundance/backend/pkg/common/validate"
	"time"
)

var (
	ErrInvalidCanonicalTag = errors.New("invalid canonical tag")
)

type CanonicalTagID string

type CanonicalTag struct {
	ID          CanonicalTagID
	TenantID    string
	Key         string
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewCanonicalTag(tenantID, key, displayName string) (*CanonicalTag, error) {
	ct := &CanonicalTag{
		ID:          CanonicalTagID(NewID()),
		TenantID:    tenantID,
		Key:         key,
		DisplayName: displayName,
		CreatedAt:   Now(),
	}

	if err := validate.ValidateStruct(ct); err != nil {
		return nil, err
	}

	return ct, nil
}

func HydrateCanonicalTag(id CanonicalTagID, tenantID, key, displayName string, createdAt, updatedAt time.Time) *CanonicalTag {
	return &CanonicalTag{
		ID:          id,
		TenantID:    tenantID,
		Key:         key,
		DisplayName: displayName,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (ct *CanonicalTag) Update(displayName string) error {
	if ct == nil {
		return ErrInvalidCanonicalTag
	}

	cpy := *ct
	cpy.DisplayName = displayName

	if err := validate.ValidateStruct(cpy); err != nil {
		return err
	}

	*ct = cpy
	ct.UpdatedAt = Now()

	return nil
}
