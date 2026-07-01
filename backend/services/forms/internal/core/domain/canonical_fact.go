package domain

type CanonicalFact struct {
	FieldID         FieldID
	TagVersionID    TagVersionID
	TagKey          string
	Value           any
	CollectionIndex *int
}

func NewCanonicalFact(fieldID FieldID, tagVersionID TagVersionID, tagKey string, value any, collectionIndex *int) *CanonicalFact {
	return &CanonicalFact{
		FieldID:         fieldID,
		TagVersionID:    tagVersionID,
		TagKey:          tagKey,
		Value:           value,
		CollectionIndex: collectionIndex,
	}
}

func HydrateCanonicalFact(fieldID FieldID, tagVersionID TagVersionID, tagKey string, value any, collectionIndex *int) *CanonicalFact {
	return &CanonicalFact{
		FieldID:         fieldID,
		TagVersionID:    tagVersionID,
		TagKey:          tagKey,
		Value:           value,
		CollectionIndex: collectionIndex,
	}
}
