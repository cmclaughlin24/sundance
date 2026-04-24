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

func toVersionDocument(v *domain.Version) *versionDocument {
	return &versionDocument{}
}

func fromVersionDocument(v *versionDocument) *domain.Version {
	return &domain.Version{}
}

type pageDocument struct {
	ID       string             `bson:"_id"`
	Key      string             `bson:"key"`
	Name     string             `bson:"name"`
	Position int                `bson:"position"`
	Sections []*sectionDocument `bson:"sections"`
	Rules    []*ruleDocument    `bson:"rules"`
}

func toPageDocument(v *domain.Page) *pageDocument {
	return &pageDocument{}
}

func fromPageDocument(v *pageDocument) *domain.Page {
	return &domain.Page{}
}

type sectionDocument struct {
	ID       string           `bson:"_id"`
	Key      string           `bson:"key"`
	Name     string           `bson:"name"`
	Position int              `bson:"position"`
	Fields   []*fieldDocument `bson:"fields"`
	Rules    []*ruleDocument  `bson:"rules"`
}

func toSectionDocument(v *domain.Section) *sectionDocument {
	return &sectionDocument{}
}

func fromSectionDocument(v *sectionDocument) *domain.Section {
	return &domain.Section{}
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

func toFieldDocument(v *domain.Field) *fieldDocument {
	return &fieldDocument{}
}

func fromFieldDocument(v *fieldDocument) *domain.Field {
	return &domain.Field{}
}

type ruleDocument struct {
	ID         string `bson:"_id"`
	Type       string `bson:"type"`
	Expression string `bson:"expression"`
}

func toRuleDocument(v *domain.Rule) *ruleDocument {
	return &ruleDocument{}
}

func fromRuleDocument(v *ruleDocument) *domain.Rule {
	return &domain.Rule{}
}

type attributeParser func(bson.Raw) (domain.FieldAttributes, error)

var attributeParserStrategies = strategy.NewStrategies[domain.FieldType, attributeParser]().
	Set(domain.FieldTypeText, func(raw bson.Raw) (domain.FieldAttributes, error) {
		return parseFieldAttributes[domain.TextFieldAttributes](raw)
	})

func unmarshalDataSourceAttributes(fieldType domain.FieldType, raw bson.Raw) (domain.FieldAttributes, error) {
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
