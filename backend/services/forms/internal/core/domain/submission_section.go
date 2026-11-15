package domain

import (
	"errors"
)

var (
	ErrDuplicateItemIndex = errors.New("invalid submission instances; duplicate item indices")
	ErrMixedItemIndex     = errors.New("an element can either be scalar or a list; cannot submission values for an element with nil and int")
)

type SubmissionSection struct {
	SectionID SectionID
	instances []SectionInstance
}

func NewSubmissionSection(id SectionID, instances []SectionInstance) (*SubmissionSection, error) {
	s := &SubmissionSection{
		SectionID: id,
		instances: instances,
	}

	if err := s.validate(); err != nil {
		return nil, err
	}

	return s, nil
}

func HydrateSubmissionSection(id SectionID, instances []SectionInstance) *SubmissionSection {
	return &SubmissionSection{
		SectionID: id,
		instances: instances,
	}
}

func (s *SubmissionSection) Instances() []SectionInstance {
	if s == nil {
		return nil
	}

	var instances []SectionInstance
	for _, i := range s.instances {
		instances = append(instances, i.clone())
	}

	return instances
}

func (s *SubmissionSection) clone() *SubmissionSection {
	return &SubmissionSection{
		SectionID: s.SectionID,
		instances: s.Instances(),
	}
}

func (s *SubmissionSection) validate() error {
	for _, i := range s.instances {
		if err := i.validate(); err != nil {
			return err
		}
	}

	return nil
}

type SectionInstance []*SubmissionValue

func (si SectionInstance) Values(elementID ElementID) []*SubmissionValue {
	if si == nil {
		return nil
	}

	var values []*SubmissionValue
	for _, v := range si {
		if v.ElementID == elementID {
			values = append(values, v)
		}
	}

	return values
}

func (si SectionInstance) clone() SectionInstance {
	if si == nil {
		return nil
	}

	cpy := make(SectionInstance, 0, len(si))
	for _, v := range si {
		// NOTE: Error is ignored as the SubmissionValue should already be in valid state.
		sv, _ := NewSubmissionValue(v.ElementID, v.Value, v.ItemIndex)
		cpy = append(cpy, sv)
	}

	return cpy
}

func (si SectionInstance) validate() error {
	if si == nil {
		return nil
	}

	type valueKey struct {
		elementID ElementID
		itemIndex int
	}

	seen := make(map[valueKey]struct{})
	kinds := make(map[ElementID]bool)

	for _, v := range si {
		idx := -1
		if v.ItemIndex != nil {
			idx = *v.ItemIndex
		}

		k := valueKey{v.ElementID, idx}
		if _, ok := seen[k]; ok {
			return ErrDuplicateItemIndex
		}

		isList := idx != -1
		if prev, ok := kinds[v.ElementID]; ok && prev != isList {
			return ErrMixedItemIndex
		}

		kinds[v.ElementID] = isList
		seen[k] = struct{}{}
	}

	return nil
}
