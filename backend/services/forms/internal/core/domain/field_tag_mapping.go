package domain

import (
	"sundance/backend/pkg/common/validate"
	"time"
)

type FieldTagMappingID string

type FieldTagMapping struct {
	ID           FieldTagMappingID
	FieldID      FieldID
	TagVersionID TagVersionID
	Priority     int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewFieldTagMapping(fieldID FieldID, tagVersionID TagVersionID, priority int) (*FieldTagMapping, error) {
	ftm := &FieldTagMapping{
		ID:           FieldTagMappingID(NewID()),
		FieldID:      fieldID,
		TagVersionID: tagVersionID,
		Priority:     priority,
		CreatedAt:    Now(),
	}

	if err := validate.ValidateStruct(ftm); err != nil {
		return nil, err
	}

	return ftm, nil
}

func HydrateFieldTagMapping(
	id FieldTagMappingID,
	fieldID FieldID,
	tagVersionID TagVersionID,
	priority int,
	createdAt time.Time,
	updatedAt time.Time,
) *FieldTagMapping {
	return &FieldTagMapping{
		ID:           id,
		FieldID:      fieldID,
		TagVersionID: tagVersionID,
		Priority:     priority,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}
