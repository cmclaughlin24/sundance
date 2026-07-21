package documents

import "sundance/backend/services/forms/internal/core/domain"

type SectionDocument struct {
	ID       string              `bson:"_id"`
	Key      string              `bson:"key"`
	Name     string              `bson:"name"`
	Position float32             `bson:"position"`
	Elements []*ElementDocument  `bson:"elements"`
	Rules    []*ruleDocument     `bson:"rules"`
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

	return &SectionDocument{
		ID:       string(s.ID),
		Key:      s.Key,
		Name:     s.Name,
		Position: s.GetPosition(),
		Elements: elementDocs,
		Rules:    rules,
	}, nil
}

func FromSectionDocument(s *SectionDocument) (*domain.Section, error) {
	section := domain.HydrateSection(
		domain.SectionID(s.ID),
		s.Key,
		s.Name,
		s.Position,
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
