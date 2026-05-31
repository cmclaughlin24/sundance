package documents

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type TagVersionDocument struct {
	ID        string    `bson:"_id"`
	TagID     string    `bson:"tag_id"`
	Version   int       `bson:"version"`
	Type      string    `bson:"type"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
	RetiredAt time.Time `bson:"retired_at"`
}

func ToTagVersionDocument(v *domain.TagVersion) TagVersionDocument {
	return TagVersionDocument{
		ID:        string(v.ID),
		TagID:     string(v.TagID),
		Version:   v.Version,
		Type:      string(v.Type),
		Status:    string(v.Status),
		CreatedAt: v.CreatedAt,
		RetiredAt: v.RetiredAt,
	}
}

func FromTagVersionDocument(d TagVersionDocument) *domain.TagVersion {
	return domain.HydrateTagVersion(
		domain.TagVersionID(d.ID),
		domain.TagID(d.TagID),
		d.Version,
		domain.TagType(d.Type),
		domain.TagStatus(d.Status),
		d.CreatedAt,
		d.RetiredAt,
	)
}
