package domain

type SubmissionValue struct {
	ElementID       ElementID
	Value           any
	CollectionIndex *int
}

func NewSubmissionValue(elementID ElementID, value any, collectionIndex *int) *SubmissionValue {
	return &SubmissionValue{
		ElementID:       elementID,
		Value:           value,
		CollectionIndex: collectionIndex,
	}
}

func HydrateSubmissionValue(elementID ElementID, value any, collectionIndex *int) *SubmissionValue {
	return &SubmissionValue{
		ElementID:       elementID,
		Value:           value,
		CollectionIndex: collectionIndex,
	}
}
