package domain

import (
	"sundance/backend/pkg/common/validate"
	"time"
)

type ElementTagMappingID string

type ElementTagMappingConfig struct {
	TagVersionID   TagVersionID
	Priority       int
	HasStaticValue bool
	StaticValue    any
}

type ElementTagMapping struct {
	ID        ElementTagMappingID
	ElementID ElementID
	CreatedAt time.Time
	UpdatedAt time.Time
	ElementTagMappingConfig
}

func NewElementTagMapping(
	elementID ElementID,
	tagVersionID TagVersionID,
	priority int,
	hasStaticValue bool,
	staticValue any,
) (*ElementTagMapping, error) {
	etm := &ElementTagMapping{
		ID:        ElementTagMappingID(NewID()),
		ElementID: elementID,
		ElementTagMappingConfig: ElementTagMappingConfig{
			TagVersionID:   tagVersionID,
			Priority:       priority,
			HasStaticValue: hasStaticValue,
			StaticValue:    staticValue,
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
	hasStaticValue bool,
	staticValue any,
	createdAt time.Time,
	updatedAt time.Time,
) *ElementTagMapping {
	return &ElementTagMapping{
		ID:        id,
		ElementID: elementID,
		ElementTagMappingConfig: ElementTagMappingConfig{
			TagVersionID:   tagVersionID,
			Priority:       priority,
			HasStaticValue: hasStaticValue,
			StaticValue:    staticValue,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
