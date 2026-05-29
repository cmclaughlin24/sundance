package domain

import (
	"sundance/backend/pkg/common/validate"
	"time"
)

type CanonicalTagVersionID string

type CanonicalTagType string

type CanonicalTagVersion struct {
	ID             CanonicalTagVersionID
	CanonicalTagID CanonicalTagID
	Version        int
	Type           CanonicalTagType
	CreatedAt      time.Time
	RetiredAt      time.Time
}

func NewCanonicalTagVersion(canonicalTagID CanonicalTagID, version int, tagType CanonicalTagType) (*CanonicalTagVersion, error) {
	ctv := &CanonicalTagVersion{
		ID:             CanonicalTagVersionID(NewID()),
		CanonicalTagID: canonicalTagID,
		Version:        version,
		Type:           tagType,
		CreatedAt:      Now(),
	}

	if err := validate.ValidateStruct(ctv); err != nil {
		return nil, err
	}

	return ctv, nil
}

func HydrateCanonicalTagVersion(id CanonicalTagVersionID, canonicalTagID CanonicalTagID, version int, tagType CanonicalTagType, createdAt, retiredAt time.Time) *CanonicalTagVersion {
	return &CanonicalTagVersion{
		ID:             id,
		CanonicalTagID: canonicalTagID,
		Version:        version,
		Type:           tagType,
		CreatedAt:      createdAt,
		RetiredAt:      retiredAt,
	}
}
