package documents

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type fieldTagMappingDocument struct {
	ID           string
	FieldID      string
	TagVersionID string
	Priority     int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func toFieldTagMappingDocument(ftm *domain.FieldTagMapping) *fieldTagMappingDocument {
	return &fieldTagMappingDocument{
		ID:           string(ftm.ID),
		FieldID:      string(ftm.FieldID),
		TagVersionID: string(ftm.TagVersionID),
		Priority:     ftm.Priority,
		CreatedAt:    ftm.CreatedAt,
		UpdatedAt:    ftm.UpdatedAt,
	}
}

func toFieldTagMappingDocuments(tags []*domain.FieldTagMapping) []*fieldTagMappingDocument {
	docs := make([]*fieldTagMappingDocument, 0, len(tags))

	for _, mapping := range tags {
		docs = append(docs, toFieldTagMappingDocument(mapping))
	}

	return docs
}

func fromFieldTagMappingDocument(doc *fieldTagMappingDocument) *domain.FieldTagMapping {
	return domain.HydrateFieldTagMapping(
		domain.FieldTagMappingID(doc.ID),
		domain.FieldID(doc.FieldID),
		domain.TagVersionID(doc.TagVersionID),
		doc.Priority,
		doc.CreatedAt,
		doc.UpdatedAt,
	)
}

func fromFieldTagMappingDocuments(docs []*fieldTagMappingDocument) []*domain.FieldTagMapping {
	tags := make([]*domain.FieldTagMapping, 0, len(docs))

	for _, doc := range docs {
		tags = append(tags, fromFieldTagMappingDocument(doc))
	}

	return tags
}
