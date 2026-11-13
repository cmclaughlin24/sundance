package domain

import "sundance/backend/pkg/common/validate"

type SubmissionValue struct {
	ElementID       ElementID
	Value           any
	CollectionIndex *int `validate:"omitempty,min=0"`
	InstanceIndex   *int `validate:"omitempty,min=0"`
}

func NewSubmissionValue(elementID ElementID, value any, collectionIndex, instanceIndex *int) (*SubmissionValue, error) {
	sv := &SubmissionValue{
		ElementID:       elementID,
		Value:           value,
		CollectionIndex: collectionIndex,
		InstanceIndex:   instanceIndex,
	}

	if err := validate.ValidateStruct(sv); err != nil {
		return nil, err
	}

	return sv, nil
}

func HydrateSubmissionValue(elementID ElementID, value any, collectionIndex, instanceIndex *int) *SubmissionValue {
	return &SubmissionValue{
		ElementID:       elementID,
		Value:           value,
		CollectionIndex: collectionIndex,
		InstanceIndex:   instanceIndex,
	}
}
