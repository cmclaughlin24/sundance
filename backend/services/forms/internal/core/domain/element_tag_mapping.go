package domain

import (
	"sundance/backend/pkg/common/validate"
	"time"
)

type ElementTagMappingID string

type ElementTagMappingConfig struct {
	TagVersionID TagVersionID
	Priority     int
}

type ElementTagMapping struct {
	ID        ElementTagMappingID
	ElementID ElementID
	CreatedAt time.Time
	UpdatedAt time.Time
	ElementTagMappingConfig
}

func NewElementTagMapping(elementID ElementID, tagVersionID TagVersionID, priority int) (*ElementTagMapping, error) {
	etm := &ElementTagMapping{
		ID:        ElementTagMappingID(NewID()),
		ElementID: elementID,
		ElementTagMappingConfig: ElementTagMappingConfig{
			TagVersionID: tagVersionID,
			Priority:     priority,
		},
		CreatedAt: Now(),
	}

	if err := validate.ValidateStruct(etm); err != nil {
		return nil, err
	}

	return etm, nil
}

func HydrateElementTagMapping(
	id ElementTagMappingID,
	elementID ElementID,
	tagVersionID TagVersionID,
	priority int,
	createdAt time.Time,
	updatedAt time.Time,
) *ElementTagMapping {
	return &ElementTagMapping{
		ID:        id,
		ElementID: elementID,
		ElementTagMappingConfig: ElementTagMappingConfig{
			TagVersionID: tagVersionID,
			Priority:     priority,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
