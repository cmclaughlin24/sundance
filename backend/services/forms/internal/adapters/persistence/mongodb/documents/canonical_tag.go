package documents

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type CanonicalTagDocument struct {
	ID          string    `bson:"_id"`
	TenantID    string    `bson:"tenant_id"`
	Key         string    `bson:"key"`
	DisplayName string    `bson:"display_name"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `wson:"updated_at"`
}

func ToCanonicalTagDocument(d *domain.CanonicalTag) CanonicalTagDocument {
	return CanonicalTagDocument{
		ID:          string(d.ID),
		TenantID:    d.TenantID,
		Key:         d.Key,
		DisplayName: d.DisplayName,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func FromCanonicalTagDocument(d CanonicalTagDocument) *domain.CanonicalTag {
	return domain.HydrateCanonicalTag(
		domain.CanonicalTagID(d.ID),
		d.TenantID,
		d.Key,
		d.DisplayName,
		d.CreatedAt,
		d.UpdatedAt,
	)
}
