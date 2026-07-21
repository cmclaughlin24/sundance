package documents

import (
	"sundance/backend/pkg/common/stratreg"
	"sundance/backend/services/forms/internal/core/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ElementDocument struct {
	ID         string                       `bson:"_id"`
	Key        string                       `bson:"key"`
	Name       string                       `bson:"name"`
	Type       string                       `bson:"type"`
	Attributes bson.Raw                     `bson:"attributes"`
	Position   float32                      `bson:"position"`
	Tags       []*elementTagMappingDocument `bson:"tags"`
	Rules      []*ruleDocument              `bson:"rules"`
}

func ToElementDocument(e *domain.Element) (*ElementDocument, error) {
	attr, err := bson.Marshal(e.Attributes)

	if err != nil {
		return nil, err
	}

	rules := RulesToDocuments(e.GetRules())
	tags := toElementTagMappingDocuments(e.GetTags())

	return &ElementDocument{
		ID:         string(e.ID),
		Key:        e.Key,
		Name:       e.Name,
		Type:       string(e.Type),
		Attributes: attr,
		Position:   e.GetPosition(),
		Tags:       tags,
		Rules:      rules,
	}, nil
}

func FromElementDocument(e *ElementDocument) (*domain.Element, error) {
	elementType := domain.ElementType(e.Type)
	attr, err := unmarshalElementAttributes(elementType, e.Attributes)

	if err != nil {
		return nil, err
	}

	element := domain.HydrateElement(
		domain.ElementID(e.ID),
		e.Key,
		e.Name,
		elementType,
		attr,
		e.Position,
		fromElementTagMappingDocuments(e.Tags),
	)

	rules, err := documentsToRules(e.Rules)
	if err != nil {
		return nil, err
	}

	if err := element.SetRules(rules...); err != nil {
		return nil, err
	}

	return element, nil
}

type attributeParser func(bson.Raw) (domain.ElementAttributes, error)

var attributeParserStrategies = stratreg.New[domain.ElementType, attributeParser]().
	Set(domain.ElementTypeText, func(raw bson.Raw) (domain.ElementAttributes, error) {
		return parseElementAttributes[*domain.TextElementAttributes](raw)
	}).
	Set(domain.ElementTypeNumber, func(raw bson.Raw) (domain.ElementAttributes, error) {
		return parseElementAttributes[*domain.NumberElementAttributes](raw)
	}).
	Set(domain.ElementTypeSelect, func(raw bson.Raw) (domain.ElementAttributes, error) {
		return parseElementAttributes[*domain.SelectElementAttributes](raw)
	}).
	Set(domain.ElementTypeCheckbox, func(raw bson.Raw) (domain.ElementAttributes, error) {
		return parseElementAttributes[*domain.CheckboxElementAttributes](raw)
	}).
	Set(domain.ElementTypeDate, func(raw bson.Raw) (domain.ElementAttributes, error) {
		return parseElementAttributes[*domain.DateElementAttributes](raw)
	})

func unmarshalElementAttributes(elementType domain.ElementType, raw bson.Raw) (domain.ElementAttributes, error) {
	strategy, err := attributeParserStrategies.Get(elementType)

	if err != nil {
		return nil, err
	}

	return strategy(raw)
}

func parseElementAttributes[T domain.ElementAttributes](raw bson.Raw) (domain.ElementAttributes, error) {
	var attr T

	if err := bson.Unmarshal(raw, &attr); err != nil {
		return nil, err
	}

	return attr, nil
}
