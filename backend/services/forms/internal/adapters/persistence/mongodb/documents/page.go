package documents

import "sundance/backend/services/forms/internal/core/domain"

type PageDocument struct {
	ID       string             `bson:"_id"`
	Key      string             `bson:"key"`
	Name     string             `bson:"name"`
	Position float32            `bson:"position"`
	Sections []*SectionDocument `bson:"sections"`
	Rules    []*ruleDocument    `bson:"rules"`
}

func ToPageDocument(p *domain.Page) (*PageDocument, error) {
	sections := p.GetSections()
	sectionDocs := make([]*SectionDocument, 0, len(sections))

	for _, s := range sections {
		doc, err := ToSectionDocument(s)

		if err != nil {
			return nil, err
		}

		sectionDocs = append(sectionDocs, doc)
	}

	rules := RulesToDocuments(p.GetRules())

	return &PageDocument{
		ID:       string(p.ID),
		Key:      p.Key,
		Name:     p.Name,
		Position: p.GetPosition(),
		Sections: sectionDocs,
		Rules:    rules,
	}, nil
}

func FromPageDocument(p *PageDocument) (*domain.Page, error) {
	page := domain.HydratePage(
		domain.PageID(p.ID),
		p.Key,
		p.Name,
		p.Position,
	)

	sections := make([]*domain.Section, 0, len(p.Sections))
	for _, s := range p.Sections {
		section, err := FromSectionDocument(s)

		if err != nil {
			return nil, err
		}

		sections = append(sections, section)
	}

	if err := page.AddSections(sections...); err != nil {
		return nil, err
	}

	rules, err := documentsToRules(p.Rules)
	if err != nil {
		return nil, err
	}

	if err := page.SetRules(rules...); err != nil {
		return nil, err
	}

	return page, nil
}
