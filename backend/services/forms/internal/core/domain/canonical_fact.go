package domain

import "slices"

type CanonicalFact struct {
	ElementID       ElementID
	TagVersionID    TagVersionID
	TagKey          string
	Value           any
	CollectionIndex []int
}

func NewCanonicalFact(elementID ElementID, tagVersionID TagVersionID, tagKey string, value any, collectionIndex []int) *CanonicalFact {
	return &CanonicalFact{
		ElementID:       elementID,
		TagVersionID:    tagVersionID,
		TagKey:          tagKey,
		Value:           value,
		CollectionIndex: slices.Clone(collectionIndex),
	}
}

func HydrateCanonicalFact(elementID ElementID, tagVersionID TagVersionID, tagKey string, value any, collectionIndex []int) *CanonicalFact {
	return &CanonicalFact{
		ElementID:       elementID,
		TagVersionID:    tagVersionID,
		TagKey:          tagKey,
		Value:           value,
		CollectionIndex: slices.Clone(collectionIndex),
	}
}
