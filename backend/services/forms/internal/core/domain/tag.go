package domain

import (
	"errors"
	"sundance/backend/pkg/common/validate"
	"time"
)

var (
	ErrInvalidCanonicalTag = errors.New("invalid tag")
)

type TagID string

type Tag struct {
	ID          TagID
	TenantID    string
	Key         string
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTag(tenantID, key, displayName string) (*Tag, error) {
	ct := &Tag{
		ID:          TagID(NewID()),
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

func HydrateTag(id TagID, tenantID, key, displayName string, createdAt, updatedAt time.Time) *Tag {
	return &Tag{
		ID:          id,
		TenantID:    tenantID,
		Key:         key,
		DisplayName: displayName,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (t *Tag) Update(displayName string) error {
	if t == nil {
		return ErrInvalidCanonicalTag
	}

	cpy := *t
	cpy.DisplayName = displayName

	if err := validate.ValidateStruct(cpy); err != nil {
		return err
	}

	*t = cpy
	t.UpdatedAt = Now()

	return nil
}
