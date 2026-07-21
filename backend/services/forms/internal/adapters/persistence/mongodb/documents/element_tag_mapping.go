package documents

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type elementTagMappingDocument struct {
	ID           string
	ElementID    string
	TagVersionID string
	Priority     int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func toElementTagMappingDocument(etm *domain.ElementTagMapping) *elementTagMappingDocument {
	return &elementTagMappingDocument{
		ID:           string(etm.ID),
		ElementID:    string(etm.ElementID),
		TagVersionID: string(etm.TagVersionID),
		Priority:     etm.Priority,
		CreatedAt:    etm.CreatedAt,
		UpdatedAt:    etm.UpdatedAt,
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
