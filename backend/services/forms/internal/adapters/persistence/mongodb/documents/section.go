package documents

import "sundance/backend/services/forms/internal/core/domain"

type SectionDocument struct {
	ID       string                 `bson:"_id"`
	Key      string                 `bson:"key"`
	Name     string                 `bson:"name"`
	Position float32                `bson:"position"`
	Repeat   *SectionRepeatDocument `bson:"repeat,omitempty"`
	Elements []*ElementDocument     `bson:"elements"`
	Rules    []*ruleDocument        `bson:"rules"`
}

type SectionRepeatDocument struct {
	MinInstances int  `bson:"min_instances"`
	MaxInstances *int `bson:"max_instances"`
}

func ToSectionDocument(s *domain.Section) (*SectionDocument, error) {
	elements := s.GetElements()
	elementDocs := make([]*ElementDocument, 0, len(elements))

	for _, e := range elements {
		doc, err := ToElementDocument(e)

		if err != nil {
			return nil, err
		}

		elementDocs = append(elementDocs, doc)
	}

	rules := RulesToDocuments(s.GetRules())

	var repeat *SectionRepeatDocument
	if s.IsRepeatable() {
		r := s.GetRepeat()
		repeat = &SectionRepeatDocument{
			MinInstances: r.MinInstances,
			MaxInstances: r.MaxInstances,
		}
	}

	return &SectionDocument{
		ID:       string(s.ID),
		Key:      s.Key,
		Name:     s.Name,
		Position: s.GetPosition(),
		Repeat:   repeat,
		Elements: elementDocs,
		Rules:    rules,
	}, nil
}

func FromSectionDocument(s *SectionDocument) (*domain.Section, error) {
	var repeat *domain.SectionRepeat
	if s.Repeat != nil {
		repeat = &domain.SectionRepeat{
			MinInstances: s.Repeat.MinInstances,
			MaxInstances: s.Repeat.MaxInstances,
		}
	}

	section := domain.HydrateSection(
		domain.SectionID(s.ID),
		s.Key,
		s.Name,
		s.Position,
		repeat,
	)

	elements := make([]*domain.Element, 0, len(s.Elements))
	for _, e := range s.Elements {
		element, err := FromElementDocument(e)

		if err != nil {
			return nil, err
		}

		elements = append(elements, element)
	}

	if err := section.AddElements(elements...); err != nil {
		return nil, err
	}

	rules, err := documentsToRules(s.Rules)
	if err != nil {
		return nil, err
	}

	if err := section.SetRules(rules...); err != nil {
		return nil, err
	}

	return section, nil
}
