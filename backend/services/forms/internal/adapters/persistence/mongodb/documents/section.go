package documents

import "sundance/backend/services/forms/internal/core/domain"

type SectionDocument struct {
	ID       string           `bson:"_id"`
	Key      string           `bson:"key"`
	Name     string           `bson:"name"`
	Position float32          `bson:"position"`
	Fields   []*FieldDocument `bson:"fields"`
	Rules    []*ruleDocument  `bson:"rules"`
}

func ToSectionDocument(s *domain.Section) (*SectionDocument, error) {
	fields := s.GetFields()
	fieldDocs := make([]*FieldDocument, 0, len(fields))

	for _, f := range fields {
		doc, err := ToFieldDocument(f)

		if err != nil {
			return nil, err
		}

		fieldDocs = append(fieldDocs, doc)
	}

	rules := RulesToDocuments(s.GetRules())

	return &SectionDocument{
		ID:       string(s.ID),
		Key:      s.Key,
		Name:     s.Name,
		Position: s.GetPosition(),
		Fields:   fieldDocs,
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

	fields := make([]*domain.Field, 0, len(s.Fields))
	for _, f := range s.Fields {
		field, err := FromFieldDocument(f)

		if err != nil {
			return nil, err
		}

		fields = append(fields, field)
	}

	if err := section.AddFields(fields...); err != nil {
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
