package documents

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type TagDocument struct {
	ID           string    `bson:"_id"`
	TenantID     string    `bson:"tenant_id"`
	Key          string    `bson:"key"`
	DisplayName  string    `bson:"display_name"`
	ValueKind    string    `bson:"value_kind"`
	IsCollection bool      `bson:"is_collection"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

func ToTagDocument(d *domain.Tag) TagDocument {
	return TagDocument{
		ID:           string(d.ID),
		TenantID:     d.TenantID,
		Key:          d.Key,
		DisplayName:  d.DisplayName,
		ValueKind:    string(d.ValueKind),
		IsCollection: d.IsCollection,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

func FromTagDocument(d TagDocument) *domain.Tag {
	return domain.HydrateTag(
		domain.TagID(d.ID),
		d.TenantID,
		d.Key,
		d.DisplayName,
		domain.TagValueKind(d.ValueKind),
		d.IsCollection,
		d.CreatedAt,
		d.UpdatedAt,
	)
}
