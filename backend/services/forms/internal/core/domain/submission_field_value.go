package domain

type SubmissionFieldValue struct {
	FieldID         FieldID
	Value           any
	CollectionIndex *int
}

func NewSubmissionFieldValue(fieldID FieldID, value any, collectionIndex *int) *SubmissionFieldValue {
	return &SubmissionFieldValue{
		FieldID:         fieldID,
		Value:           value,
		CollectionIndex: collectionIndex,
	}
}

func HydrateSubmissionFieldValue(fieldID FieldID, value any, collectionIndex *int) *SubmissionFieldValue {
	return &SubmissionFieldValue{
		FieldID:         fieldID,
		Value:           value,
		CollectionIndex: collectionIndex,
	}
}
