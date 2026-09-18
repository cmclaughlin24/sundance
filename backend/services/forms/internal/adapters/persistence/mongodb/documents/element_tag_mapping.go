package documents

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type elementTagMappingDocument struct {
	ID             string    `bson:"_id"`
	ElementID      string    `bson:"element_id"`
	TagVersionID   string    `bson:"tag_version_id"`
	Priority       int       `bson:"priority"`
	HasStaticValue bool      `bson:"has_static_value"`
	StaticValue    any       `bson:"static_value"`
	CreatedAt      time.Time `bson:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at"`
}

func toElementTagMappingDocument(etm *domain.ElementTagMapping) *elementTagMappingDocument {
	return &elementTagMappingDocument{
		ID:             string(etm.ID),
		ElementID:      string(etm.ElementID),
		TagVersionID:   string(etm.TagVersionID),
		Priority:       etm.Priority,
		HasStaticValue: etm.HasStaticValue,
		StaticValue:    etm.StaticValue,
		CreatedAt:      etm.CreatedAt,
		UpdatedAt:      etm.UpdatedAt,
	}
}

func toElementTagMappingDocuments(tags []*domain.ElementTagMapping) []*elementTagMappingDocument {
	docs := make([]*elementTagMappingDocument, 0, len(tags))

	for _, mapping := range tags {
		docs = append(docs, toElementTagMappingDocument(mapping))
	}

	return docs
}

func fromElementTagMappingDocument(doc *elementTagMappingDocument) *domain.ElementTagMapping {
	return domain.HydrateElementTagMapping(
		domain.ElementTagMappingID(doc.ID),
		domain.ElementID(doc.ElementID),
		domain.TagVersionID(doc.TagVersionID),
		doc.Priority,
		doc.HasStaticValue,
		doc.StaticValue,
		doc.CreatedAt,
		doc.UpdatedAt,
	)
}

func fromElementTagMappingDocuments(docs []*elementTagMappingDocument) []*domain.ElementTagMapping {
	tags := make([]*domain.ElementTagMapping, 0, len(docs))

	for _, doc := range docs {
		tags = append(tags, fromElementTagMappingDocument(doc))
	}

	return tags
}
