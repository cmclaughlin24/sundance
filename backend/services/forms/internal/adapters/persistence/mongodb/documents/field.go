package documents

import (
	"sundance/backend/pkg/common/stratreg"
	"sundance/backend/services/forms/internal/core/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FieldDocument struct {
	ID         string                     `bson:"_id"`
	Key        string                     `bson:"key"`
	Name       string                     `bson:"name"`
	Type       string                     `bson:"type"`
	Attributes bson.Raw                   `bson:"attributes"`
	Position   float32                    `bson:"position"`
	Tags       []*fieldTagMappingDocument `bson:"tags"`
	Rules      []*ruleDocument            `bson:"rules"`
}

func ToFieldDocument(f *domain.Field) (*FieldDocument, error) {
	attr, err := bson.Marshal(f.Attributes)

	if err != nil {
		return nil, err
	}

	rules := RulesToDocuments(f.GetRules())
	tags := toFieldTagMappingDocuments(f.GetTags())

	return &FieldDocument{
		ID:         string(f.ID),
		Key:        f.Key,
		Name:       f.Name,
		Type:       string(f.Type),
		Attributes: attr,
		Position:   f.GetPosition(),
		Tags:       tags,
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
		fromFieldTagMappingDocuments(f.Tags),
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

type attributeParser func(bson.Raw) (domain.FieldAttributes, error)

var attributeParserStrategies = stratreg.New[domain.FieldType, attributeParser]().
	Set(domain.FieldTypeText, func(raw bson.Raw) (domain.FieldAttributes, error) {
		return parseFieldAttributes[*domain.TextFieldAttributes](raw)
	}).
	Set(domain.FieldTypeNumber, func(raw bson.Raw) (domain.FieldAttributes, error) {
		return parseFieldAttributes[*domain.NumberFieldAttributes](raw)
	}).
	Set(domain.FieldTypeSelect, func(raw bson.Raw) (domain.FieldAttributes, error) {
		return parseFieldAttributes[*domain.SelectFieldAttributes](raw)
	}).
	Set(domain.FieldTypeCheckbox, func(raw bson.Raw) (domain.FieldAttributes, error) {
		return parseFieldAttributes[*domain.CheckboxFieldAttributes](raw)
	}).
	Set(domain.FieldTypeDate, func(raw bson.Raw) (domain.FieldAttributes, error) {
		return parseFieldAttributes[*domain.DateFieldAttributes](raw)
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
