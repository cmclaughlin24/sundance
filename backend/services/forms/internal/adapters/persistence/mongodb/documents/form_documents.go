package documents

import (
	"time"

	"sundance/backend/pkg/common/stratreg"
	"sundance/backend/services/forms/internal/core/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type FormDocument struct {
	ID          string    `bson:"_id"`
	TenantID    string    `bson:"tenant_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func ToFormDocument(f *domain.Form) *FormDocument {
	return &FormDocument{
		ID:          string(f.ID),
		TenantID:    f.TenantID,
		Name:        f.Name,
		Description: f.Description,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}
}

func FromFormDocument(f *FormDocument) *domain.Form {
	return domain.HydrateForm(
		domain.FormID(f.ID),
		f.TenantID,
		f.Name,
		f.Description,
		f.CreatedAt,
		f.UpdatedAt,
	)
}

type VersionDocument struct {
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
	Pages       []*PageDocument `bson:"pages"`
}

func ToVersionDocument(v *domain.Version) (*VersionDocument, error) {
	pages := v.GetPages()
	pageDocs := make([]*PageDocument, 0, len(pages))

	for _, p := range pages {
		doc, err := ToPageDocument(p)

		if err != nil {
			return nil, err
		}

		pageDocs = append(pageDocs, doc)
	}

	return &VersionDocument{
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
		Pages:       pageDocs,
	}, nil
}

func FromVersionDocument(v *VersionDocument) (*domain.Version, error) {
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
		page, err := FromPageDocument(p)

		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	if err := version.AddPages(pages...); err != nil {
		return nil, err
	}

	return version, nil
}

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

type FieldDocument struct {
	ID         string          `bson:"_id"`
	Key        string          `bson:"key"`
	Name       string          `bson:"name"`
	Type       string          `bson:"type"`
	Attributes bson.Raw        `bson:"attributes"`
	Position   float32         `bson:"position"`
	Rules      []*ruleDocument `bson:"rules"`
}

func ToFieldDocument(f *domain.Field) (*FieldDocument, error) {
	attr, err := bson.Marshal(f.Attributes)

	if err != nil {
		return nil, err
	}

	rules := RulesToDocuments(f.GetRules())

	return &FieldDocument{
		ID:         string(f.ID),
		Key:        f.Key,
		Name:       f.Name,
		Type:       string(f.Type),
		Attributes: attr,
		Position:   f.GetPosition(),
		Rules:      rules,
	}, nil
}

func FromFieldDocument(f *FieldDocument) (*domain.Field, error) {
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

	rules, err := documentsToRules(f.Rules)
	if err != nil {
		return nil, err
	}

	if err := field.SetRules(rules...); err != nil {
		return nil, err
	}

	return field, nil
}

type ruleDocument struct {
	ID          string                    `bson:"_id"`
	Type        string                    `bson:"type"`
	Expressions []*ruleExpressionDocument `bson:"expressions"`
}

func RulesToDocuments(rules map[domain.RuleType]*domain.Rule) []*ruleDocument {
	documents := make([]*ruleDocument, 0, len(rules))
	for _, r := range rules {
		documents = append(documents, toRuleDocument(r))
	}
	return documents
}

func toRuleDocument(r *domain.Rule) *ruleDocument {
	expressions := r.GetExpressions()
	documents := make([]*ruleExpressionDocument, 0, len(expressions))

	for _, e := range expressions {
		documents = append(documents, toRuleExpressionDocument(e))
	}

	return &ruleDocument{
		ID:          string(r.ID),
		Type:        string(r.Type),
		Expressions: documents,
	}
}

func documentsToRules(documents []*ruleDocument) ([]*domain.Rule, error) {
	rules := make([]*domain.Rule, 0, len(documents))
	for _, d := range documents {
		rule, err := fromRuleDocument(d)
		if err != nil {
			return nil, err
		}
		
		rules = append(rules, rule)
	}

	return rules, nil
}

func fromRuleDocument(doc *ruleDocument) (*domain.Rule, error) {
	r := domain.HydrateRule(
		domain.RuleID(doc.ID),
		domain.RuleType(doc.Type),
	)

	expressions := make([]*domain.RuleExpression, 0, len(doc.Expressions))
	for _, e := range doc.Expressions {
		expressions = append(expressions, fromRuleExpressionDocument(e))
	}

	if err := r.AddExpressions(expressions...); err != nil {
		return nil, err
	}

	return r, nil
}

type ruleExpressionDocument struct {
	FieldID          string  `bson:"field_id"`
	Operator         string  `bson:"operator"`
	Value            any     `bson:"value"`
	JoinWithPrevious *string `bson:"join_with_previous"`
	Position         float32 `bson:"position"`
}

func toRuleExpressionDocument(e *domain.RuleExpression) *ruleExpressionDocument {
	return &ruleExpressionDocument{
		FieldID:          string(e.FieldID),
		Operator:         string(e.Operator),
		Value:            e.Value,
		JoinWithPrevious: (*string)(e.JoinWithPrevious),
		Position:         e.GetPosition(),
	}
}

func fromRuleExpressionDocument(e *ruleExpressionDocument) *domain.RuleExpression {
	return domain.HydrateRuleExpression(
		domain.FieldID(e.FieldID),
		domain.ExpressionOperator(e.Operator),
		e.Value,
		(*domain.JoinOperator)(e.JoinWithPrevious),
		e.Position,
	)
}

type attributeParser func(bson.Raw) (domain.FieldAttributes, error)

var attributeParserStrategies = stratreg.New[domain.FieldType, attributeParser]().
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
