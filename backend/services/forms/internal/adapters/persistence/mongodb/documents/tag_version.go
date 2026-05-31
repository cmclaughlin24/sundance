package documents

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type TagVersionDocument struct {
	ID           string    `bson:"_id"`
	TagID        string    `bson:"tag_id"`
	Version      int       `bson:"version"`
	Type         string    `bson:"type"`
	Status       string    `bson:"status"`
	CreatedAt    time.Time `bson:"created_at"`
	DeprecatedAt time.Time `bson:"deprecated_at"`
	PublishedAt  time.Time `bson:"published_at,omitempty"`
	RetiredAt    time.Time `bson:"retired_at,omitempty"`
}

func ToTagVersionDocument(tv *domain.TagVersion) TagVersionDocument {
	return TagVersionDocument{
		ID:           string(tv.ID),
		TagID:        string(tv.TagID),
		Version:      tv.Version,
		Type:         string(tv.Type),
		Status:       string(tv.Status),
		CreatedAt:    tv.CreatedAt,
		DeprecatedAt: tv.DeprecatedAt,
		PublishedAt:  tv.PublishedAt,
		RetiredAt:    tv.RetiredAt,
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
		d.DeprecatedAt,
		d.PublishedAt,
		d.RetiredAt,
	)
}
