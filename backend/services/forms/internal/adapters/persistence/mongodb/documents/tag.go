package documents

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type TagDocument struct {
	ID            string    `bson:"_id"`
	TenantID      string    `bson:"tenant_id"`
	Key           string    `bson:"key"`
	DisplayName   string    `bson:"display_name"`
	NodeType      string    `bson:"node_type"`
	PrimitiveType *string   `bson:"primitive_type"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

func ToTagDocument(d *domain.Tag) TagDocument {
	var primitiveType *string
	if d.PrimitiveType != nil {
		pt := string(*d.PrimitiveType)
		primitiveType = &pt
	}

	return TagDocument{
		ID:            string(d.ID),
		TenantID:      d.TenantID,
		Key:           d.Key,
		DisplayName:   d.DisplayName,
		NodeType:      string(d.NodeType),
		PrimitiveType: primitiveType,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}

func FromTagDocument(d TagDocument) *domain.Tag {
	var primitiveType *domain.TagPrimitiveType
	if d.PrimitiveType != nil {
		pt := domain.TagPrimitiveType(*d.PrimitiveType)
		primitiveType = &pt
	}

	return domain.HydrateTag(
		domain.TagID(d.ID),
		d.TenantID,
		d.Key,
		d.DisplayName,
		domain.TagNodeType(d.NodeType),
		primitiveType,
		d.CreatedAt,
		d.UpdatedAt,
	)
}
