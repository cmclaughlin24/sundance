package domain

import "sundance/backend/pkg/common/validate"

type SubmissionValue struct {
	ElementID ElementID
	Value     any
	ItemIndex *int `validate:"omitempty,min=0"`
}

func NewSubmissionValue(elementID ElementID, value any, itemIndex *int) (*SubmissionValue, error) {
	sv := &SubmissionValue{
		ElementID: elementID,
		Value:     value,
		ItemIndex: itemIndex,
	}

	if err := validate.ValidateStruct(sv); err != nil {
		return nil, err
	}

	return sv, nil
}

func HydrateSubmissionValue(elementID ElementID, value any, itemIndex *int) *SubmissionValue {
	return &SubmissionValue{
		ElementID: elementID,
		Value:     value,
		ItemIndex: itemIndex,
	}
}
