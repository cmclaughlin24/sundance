package domain

import "sundance/backend/pkg/common/validate"

type SubmissionValue struct {
	ElementID    ElementID
	Value        any
	SectionIndex *int `validate:"omitempty,min=0"`
	ItemIndex    *int `validate:"omitempty,min=0"`
}

func NewSubmissionValue(elementID ElementID, value any, sectionIndex, itemIndex *int) (*SubmissionValue, error) {
	sv := &SubmissionValue{
		ElementID:    elementID,
		Value:        value,
		SectionIndex: sectionIndex,
		ItemIndex:    itemIndex,
	}

	if err := validate.ValidateStruct(sv); err != nil {
		return nil, err
	}

	return sv, nil
}

func HydrateSubmissionValue(elementID ElementID, value any, sectionIndex, itemIndex *int) *SubmissionValue {
	return &SubmissionValue{
		ElementID:    elementID,
		Value:        value,
		SectionIndex: sectionIndex,
		ItemIndex:    itemIndex,
	}
}

func (sv *SubmissionValue) IndexPath() []int {
	if sv == nil {
		return nil
	}

	path := make([]int, 0, 2)

	if sv.SectionIndex != nil {
		path = append(path, *sv.SectionIndex)
	}

	if sv.ItemIndex != nil {
		path = append(path, *sv.ItemIndex)
	}

	return path
}
