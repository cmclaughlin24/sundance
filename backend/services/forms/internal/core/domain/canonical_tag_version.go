package domain

import (
	"sundance/backend/pkg/common/validate"
	"time"
)

type CanonicalTagVersionID string

type CanonicalTagType string

type CanonicalTagStatus string

const (
	CanonicalTagStatusDraft      CanonicalTagStatus = "draft"
	CanonicalTagStatusActive     CanonicalTagStatus = "active"
	CanonicalTagStatusDeprecated CanonicalTagStatus = "deprecated"
	CanonicalTagStatusRetired    CanonicalTagStatus = "retired"
)

type CanonicalTagVersion struct {
	ID             CanonicalTagVersionID
	CanonicalTagID CanonicalTagID
	Version        int
	Type           CanonicalTagType
	Status         CanonicalTagStatus
	CreatedAt      time.Time
	RetiredAt      time.Time
}

func NewCanonicalTagVersion(canonicalTagID CanonicalTagID, version int, tagType CanonicalTagType) (*CanonicalTagVersion, error) {
	ctv := &CanonicalTagVersion{
		ID:             CanonicalTagVersionID(NewID()),
		CanonicalTagID: canonicalTagID,
		Version:        version,
		Status:         CanonicalTagStatusDraft,
		Type:           tagType,
		CreatedAt:      Now(),
	}

	if err := validate.ValidateStruct(ctv); err != nil {
		return nil, err
	}

	return ctv, nil
}

func HydrateCanonicalTagVersion(
	id CanonicalTagVersionID,
	canonicalTagID CanonicalTagID,
	version int,
	tagType CanonicalTagType,
	status CanonicalTagStatus,
	createdAt, 
	retiredAt time.Time,
) *CanonicalTagVersion {
	return &CanonicalTagVersion{
		ID:             id,
		CanonicalTagID: canonicalTagID,
		Version:        version,
		Type:           tagType,
		Status:         status,
		CreatedAt:      createdAt,
		RetiredAt:      retiredAt,
	}
}
