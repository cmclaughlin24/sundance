package services

import (
	"errors"
	"fmt"
	"sundance/backend/pkg/common"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports/commands"
)

var (
	ErrDuplicateFormKey     = errors.New("duplicate form key")
	ErrInvalidExpressionKey = errors.New("invalid expression key")
)

type ruleUpdater interface {
	ReplaceRules(...*domain.Rule) error
}

type formDefinitionMapper struct {
	pageData       []commands.PageData
	formKeys       map[string]int
	expressionKeys map[string]bool
}

func newFormDefinitionMapper(pageData []commands.PageData) *formDefinitionMapper {
	return &formDefinitionMapper{
		pageData:       pageData,
		formKeys:       make(map[string]int),
		expressionKeys: make(map[string]bool),
	}
}

func (m *formDefinitionMapper) createFormDefinition() ([]*domain.Page, error) {
	pages := make([]*domain.Page, 0, len(m.pageData))
	for _, p := range m.pageData {
		page, err := m.createPage(p)
		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	if err := m.validate(); err != nil {
		return nil, err
	}

	return pages, nil
}

func (m *formDefinitionMapper) updateFormDefinition(version *domain.FormVersion) ([]*domain.Page, error) {
	pages := make([]*domain.Page, 0, len(m.pageData))
	for _, p := range m.pageData {
		if p.ID != nil {
			page := version.GetPage(domain.PageID(*p.ID))

			if page == nil {
				return nil, fmt.Errorf("%w; page id=%s", common.ErrNotFound, *p.ID)
			}

			if err := m.updatePage(p, page); err != nil {
				return nil, err
			}

			pages = append(pages, page)
		} else {
			page, err := m.createPage(p)
			if err != nil {
				return nil, err
			}

			pages = append(pages, page)
		}
	}

	if err := m.validate(); err != nil {
		return nil, err
	}

	return pages, nil
}

func (m *formDefinitionMapper) createPage(p commands.PageData) (*domain.Page, error) {
	page, err := domain.NewPage(p.Name, p.Key, p.Position)
	if err != nil {
		return nil, err
	}

	sections := make([]*domain.Section, 0, len(p.SectionsData))
	for _, s := range p.SectionsData {
		section, err := m.createSection(s)
		if err != nil {
			return nil, err
		}

		sections = append(sections, section)
	}

	if err := page.AddSections(sections...); err != nil {
		return nil, err
	}

	rules, err := m.createRules(p.Rules)
	if err != nil {
		return nil, err
	}

	if err := page.SetRules(rules...); err != nil {
		return nil, err
	}

	m.trackFormKey(page.Key)
	return page, nil
}

func (m *formDefinitionMapper) updatePage(p commands.PageData, page *domain.Page) error {
	if err := page.Update(p.Key, p.Name, p.Position); err != nil {
		return err
	}

	sections := make([]*domain.Section, 0, len(p.SectionsData))
	for _, s := range p.SectionsData {
		if s.ID != nil {
			section := page.GetSection(domain.SectionID(*s.ID))

			if section == nil {
				return fmt.Errorf("%w; section id=%s", common.ErrNotFound, *s.ID)
			}

			if err := m.updateSection(s, section); err != nil {
				return err
			}

			sections = append(sections, section)
		} else {
			section, err := m.createSection(s)
			if err != nil {
				return err
			}

			sections = append(sections, section)
		}
	}

	if err := page.ReplaceSections(sections...); err != nil {
		return err
	}

	if err := m.updateRules(p.Rules, page); err != nil {
		return err
	}

	m.trackFormKey(page.Key)
	return nil
}

func (m *formDefinitionMapper) createSection(s commands.SectionData) (*domain.Section, error) {
	section, err := domain.NewSection(s.Key, s.Name, s.Position)
	if err != nil {
		return nil, err
	}

	elements := make([]*domain.Element, 0, len(s.ElementsData))
	for _, f := range s.ElementsData {
		element, err := m.createElement(f)
		if err != nil {
			return nil, err
		}

		elements = append(elements, element)
	}

	if err := section.AddElements(elements...); err != nil {
		return nil, err
	}

	rules, err := m.createRules(s.Rules)
	if err != nil {
		return nil, err
	}

	if err := section.SetRules(rules...); err != nil {
		return nil, err
	}

	m.trackFormKey(section.Key)
	return section, nil
}

func (m *formDefinitionMapper) updateSection(s commands.SectionData, section *domain.Section) error {
	if err := section.Update(s.Key, s.Name, s.Position); err != nil {
		return err
	}

	elements := make([]*domain.Element, 0, len(s.ElementsData))
	for _, f := range s.ElementsData {
		if f.ID != nil {
			element := section.GetElement(domain.ElementID(*f.ID))

			if element == nil {
				return fmt.Errorf("%w; element id=%s", common.ErrNotFound, *f.ID)
			}

			if err := m.updateElement(f, element); err != nil {
				return err
			}

			elements = append(elements, element)
		} else {
			element, err := m.createElement(f)
			if err != nil {
				return err
			}

			elements = append(elements, element)
		}
	}

	if err := section.ReplaceElements(elements...); err != nil {
		return err
	}

	if err := m.updateRules(s.Rules, section); err != nil {
		return err
	}

	m.trackFormKey(section.Key)
	return nil
}

func (m *formDefinitionMapper) createElement(f commands.ElementData) (*domain.Element, error) {
	element, err := domain.NewElement(f.Key, f.Name, domain.ElementType(f.Type), f.Attributes, f.Position)
	if err != nil {
		return nil, err
	}

	for _, t := range f.Tags {
		err := element.AddTags(domain.ElementTagMappingConfig{
			TagVersionID: domain.TagVersionID(t.TagVersionID),
			Priority:     t.Priority,
		})

		if err != nil {
			return nil, err
		}
	}

	rules, err := m.createRules(f.Rules)
	if err != nil {
		return nil, err
	}

	if err := element.SetRules(rules...); err != nil {
		return nil, err
	}

	m.trackFormKey(element.Key)
	m.trackExpressionKeys(element.Attributes.GetReferencedKeys()...)
	return element, nil
}

func (m *formDefinitionMapper) updateElement(f commands.ElementData, element *domain.Element) error {
	if err := element.Update(f.Key, f.Name, domain.ElementType(f.Type), f.Attributes, f.Position); err != nil {
		return err
	}

	tags := make([]domain.ElementTagMappingConfig, 0, len(f.Tags))
	for _, etm := range f.Tags {
		tags = append(tags, domain.ElementTagMappingConfig{
			TagVersionID: domain.TagVersionID(etm.TagVersionID),
			Priority:     etm.Priority,
		})
	}

	if err := element.ReplaceTags(tags...); err != nil {
		return err
	}

	if err := m.updateRules(f.Rules, element); err != nil {
		return err
	}

	m.trackFormKey(element.Key)
	m.trackExpressionKeys(element.Attributes.GetReferencedKeys()...)
	return nil
}

func (m *formDefinitionMapper) createRules(ruleData []commands.RuleData) ([]*domain.Rule, error) {
	rules := make([]*domain.Rule, 0, len(ruleData))

	for _, r := range ruleData {
		rule, err := m.createRule(r)
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

func (m *formDefinitionMapper) updateRules(ruleData []commands.RuleData, updater ruleUpdater) error {
	rules := make([]*domain.Rule, 0, len(ruleData))

	for _, r := range ruleData {
		rule, err := m.createRule(r)
		if err != nil {
			return err
		}

		rules = append(rules, rule)
	}

	if err := updater.ReplaceRules(rules...); err != nil {
		return err
	}

	return nil
}

func (m *formDefinitionMapper) createRule(r commands.RuleData) (*domain.Rule, error) {
	rule, err := domain.NewRule(domain.RuleType(r.Type))
	if err != nil {
		return nil, err
	}

	for _, re := range r.Expressions {
		expression, err := domain.NewRuleExpression(
			re.FieldKey,
			domain.ExprOperator(re.Operator),
			re.Value,
			(*domain.JoinOperator)(re.JoinWithPrevious),
			re.Position,
		)
		if err != nil {
			return nil, err
		}

		if err := rule.AddExpressions(expression); err != nil {
			return nil, err
		}

		m.trackExpressionKeys(expression.FieldKey)
	}

	return rule, nil
}

func (m *formDefinitionMapper) trackFormKey(key string) {
	count, ok := m.formKeys[key]

	if !ok {
		m.formKeys[key] = 1
	} else {
		m.formKeys[key] = count + 1
	}
}

func (m *formDefinitionMapper) trackExpressionKeys(keys ...string) {
	for _, key := range keys {
		m.expressionKeys[key] = true
	}
}

func (m *formDefinitionMapper) validate() error {
	for key, value := range m.formKeys {
		if value > 1 {
			return fmt.Errorf("%w; key=%s", ErrDuplicateFormKey, key)
		}
	}

	for key := range m.expressionKeys {
		if _, ok := m.formKeys[key]; !ok {
			return fmt.Errorf("%w; key=%s used in expression but no element with such key exists", ErrInvalidExpressionKey, key)
		}
	}

	return nil
}
