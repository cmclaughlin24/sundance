package domain

type CanonicalFact struct {
	FieldID      FieldID
	TagVersionID TagVersionID
	TagKey       string
	Value        any
}

func NewCanonicalFact(fieldID FieldID, tagVersionID TagVersionID, tagKey string, value any) *CanonicalFact {
	return &CanonicalFact{
		FieldID:      fieldID,
		TagVersionID: tagVersionID,
		TagKey:       tagKey,
		Value:        value,
	}
}
