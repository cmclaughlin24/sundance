package mongodb

import (
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/strategy"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type formDocument struct {
	ID          string    `bson:"_id"`
	TenantID    string    `bson:"tenant_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func toFormDocument(f *domain.Form) *formDocument {
	return &formDocument{
		ID:          string(f.ID),
		TenantID:    f.TenantID,
		Name:        f.Name,
		Description: f.Description,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}
}

func fromFormDocument(f *formDocument) *domain.Form {
	return domain.HydrateForm(
		domain.FormID(f.ID),
		f.TenantID,
		f.Name,
		f.Description,
		f.CreatedAt,
		f.UpdatedAt,
	)
}

type versionDocument struct {
	ID          string          `bson:"_id"`
	FormID      string          `bson:"form_id"`
	Version     int             `bson:"version"`
	Status      string          `bson:"status"`
	PublishedBy string          `bson:"published_by"`
	PublishedAt time.Time       `bson:"published_at"`
	RetiredBy   string          `bson:"retired_by"`
	RetiredAt   time.Time       `bson:"retired_at"`
	CreatedAt   time.Time       `bson:"created_at"`
	UpdatedAt   time.Time       `bson:"updated_at"`
	Pages       []*pageDocument `bson:"pages"`
}

func toVersionDocument(v *domain.Version) (*versionDocument, error) {
	pages := make([]*pageDocument, 0, len(v.GetPages()))

	for _, p := range v.GetPages() {
		doc, err := toPageDocument(p)

		if err != nil {
			return nil, err
		}

		pages = append(pages, doc)
	}

	return &versionDocument{
		ID:          string(v.ID),
		FormID:      string(v.FormID),
		Version:     v.Version,
		Status:      string(v.Status),
		PublishedBy: v.PublishedBy,
		PublishedAt: v.PublishedAt,
		RetiredBy:   v.RetiredBy,
		RetiredAt:   v.RetiredAt,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
		Pages:       pages,
	}, nil
}

func fromVersionDocument(v *versionDocument) (*domain.Version, error) {
	version := domain.HydrateVersion(
		domain.VersionID(v.ID),
		domain.FormID(v.FormID),
		v.Version,
		domain.VersionStatus(v.Status),
		v.PublishedBy,
		v.PublishedAt,
		v.RetiredBy,
		v.RetiredAt,
		v.CreatedAt,
		v.UpdatedAt,
	)

	pages := make([]*domain.Page, 0, len(v.Pages))
	for _, p := range v.Pages {
		page, err := fromPageDocument(p)

		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	if err := version.SetPages(pages...); err != nil {
		return nil, err
	}

	return version, nil
}

type pageDocument struct {
	ID       string             `bson:"_id"`
	Key      string             `bson:"key"`
	Name     string             `bson:"name"`
	Position int                `bson:"position"`
	Sections []*sectionDocument `bson:"sections"`
	Rules    []*ruleDocument    `bson:"rules"`
}

func toPageDocument(p *domain.Page) (*pageDocument, error) {
	sections := make([]*sectionDocument, 0, len(p.GetSections()))

	for _, s := range p.GetSections() {
		doc, err := toSectionDocument(s)

		if err != nil {
			return nil, err
		}

		sections = append(sections, doc)
	}

	rules := rulesToDocuments(p.GetRules())

	return &pageDocument{
		ID:       string(p.ID),
		Key:      p.Key,
		Name:     p.Name,
		Position: p.GetPosition(),
		Sections: sections,
		Rules:    rules,
	}, nil
}

func fromPageDocument(p *pageDocument) (*domain.Page, error) {
	page := domain.HydratePage(
		domain.PageID(p.ID),
		p.Key,
		p.Name,
		p.Position,
	)

	sections := make([]*domain.Section, 0, len(p.Sections))
	for _, s := range p.Sections {
		section, err := fromSectionDocument(s)

		if err != nil {
			return nil, err
		}

		sections = append(sections, section)
	}

	if err := page.SetSections(sections...); err != nil {
		return nil, err
	}

	rules := documentsToRules(p.Rules)
	if err := page.SetRules(rules...); err != nil {
		return nil, err
	}

	return page, nil
}

type sectionDocument struct {
	ID       string           `bson:"_id"`
	Key      string           `bson:"key"`
	Name     string           `bson:"name"`
	Position int              `bson:"position"`
	Fields   []*fieldDocument `bson:"fields"`
	Rules    []*ruleDocument  `bson:"rules"`
}

func toSectionDocument(s *domain.Section) (*sectionDocument, error) {
	fields := make([]*fieldDocument, 0, len(s.GetFields()))

	for _, f := range s.GetFields() {
		doc, err := toFieldDocument(f)

		if err != nil {
			return nil, err
		}

		fields = append(fields, doc)
	}

	rules := rulesToDocuments(s.GetRules())

	return &sectionDocument{
		ID:       string(s.ID),
		Key:      s.Key,
		Name:     s.Name,
		Position: s.GetPosition(),
		Fields:   fields,
		Rules:    rules,
	}, nil
}

func fromSectionDocument(s *sectionDocument) (*domain.Section, error) {
	section := domain.HydrateSection(
		domain.SectionID(s.ID),
		s.Key,
		s.Name,
		s.Position,
	)

	fields := make([]*domain.Field, 0, len(s.Fields))
	for _, f := range s.Fields {
		field, err := fromFieldDocument(f)

		if err != nil {
			return nil, err
		}

		fields = append(fields, field)
	}

	if err := section.SetFields(fields...); err != nil {
		return nil, err
	}

	rules := documentsToRules(s.Rules)
	if err := section.SetRules(rules...); err != nil {
		return nil, err
	}

	return section, nil
}

type fieldDocument struct {
	ID         string          `bson:"_id"`
	Key        string          `bson:"key"`
	Name       string          `bson:"name"`
	Type       string          `bson:"type"`
	Attributes bson.Raw        `bson:"attributes"`
	Position   int             `bson:"position"`
	Rules      []*ruleDocument `bson:"rules"`
}

func toFieldDocument(f *domain.Field) (*fieldDocument, error) {
	attr, err := bson.Marshal(f.Attributes)

	if err != nil {
		return nil, err
	}

	rules := rulesToDocuments(f.GetRules())

	return &fieldDocument{
		ID:         string(f.ID),
		Key:        f.Key,
		Name:       f.Name,
		Type:       string(f.Type),
		Attributes: attr,
		Position:   f.GetPosition(),
		Rules:      rules,
	}, nil
}

func fromFieldDocument(f *fieldDocument) (*domain.Field, error) {
	fieldType := domain.FieldType(f.Type)
	attr, err := unmarshalFieldAttributes(fieldType, f.Attributes)

	if err != nil {
		return nil, err
	}

	field := domain.HydrateField(
		domain.FieldID(f.ID),
		f.Key,
		f.Name,
		fieldType,
		attr,
		f.Position,
	)

	rules := documentsToRules(f.Rules)
	if err := field.SetRules(rules...); err != nil {
		return nil, err
	}

	return field, nil
}

type ruleDocument struct {
	ID         string `bson:"_id"`
	Type       string `bson:"type"`
	Expression string `bson:"expression"`
}

func rulesToDocuments(rules map[domain.RuleType]*domain.Rule) []*ruleDocument {
	documents := make([]*ruleDocument, 0, len(rules))
	for _, r := range rules {
		documents = append(documents, toRuleDocument(r))
	}
	return documents
}

func toRuleDocument(r *domain.Rule) *ruleDocument {
	return &ruleDocument{
		ID:         string(r.ID),
		Type:       string(r.Type),
		Expression: r.Expression,
	}
}

func documentsToRules(documents []*ruleDocument) []*domain.Rule {
	rules := make([]*domain.Rule, 0, len(documents))
	for _, d := range documents {
		rules = append(rules, fromRuleDocument(d))
	}
	return rules
}

func fromRuleDocument(r *ruleDocument) *domain.Rule {
	return domain.HydrateRule(
		domain.RuleID(r.ID),
		domain.RuleType(r.Type),
		r.Expression,
	)
}

type attributeParser func(bson.Raw) (domain.FieldAttributes, error)

var attributeParserStrategies = strategy.NewStrategies[domain.FieldType, attributeParser]().
	Set(domain.FieldTypeText, func(raw bson.Raw) (domain.FieldAttributes, error) {
		return parseFieldAttributes[domain.TextFieldAttributes](raw)
	}).
	Set(domain.FieldTypeNumber, func(raw bson.Raw) (domain.FieldAttributes, error) {
		return parseFieldAttributes[domain.NumberFieldAttributes](raw)
	}).
	Set(domain.FieldTypeSelect, func(raw bson.Raw) (domain.FieldAttributes, error) {
		return parseFieldAttributes[domain.SelectFieldAttributes](raw)
	}).
	Set(domain.FieldTypeCheckbox, func(raw bson.Raw) (domain.FieldAttributes, error) {
		return parseFieldAttributes[domain.CheckboxFieldAttributes](raw)
	}).
	Set(domain.FieldTypeDate, func(raw bson.Raw) (domain.FieldAttributes, error) {
		return parseFieldAttributes[domain.DateFieldAttributes](raw)
	})

func unmarshalFieldAttributes(fieldType domain.FieldType, raw bson.Raw) (domain.FieldAttributes, error) {
	strategy, err := attributeParserStrategies.Get(fieldType)

	if err != nil {
		return nil, err
	}

	return strategy(raw)
}

func parseFieldAttributes[T domain.FieldAttributes](raw bson.Raw) (domain.FieldAttributes, error) {
	var attr T

	if err := bson.Unmarshal(raw, &attr); err != nil {
		return nil, err
	}

	return attr, nil
}
