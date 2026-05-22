package documents

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type VersionDocument struct {
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
	Pages       []*PageDocument `bson:"pages"`
}

func ToVersionDocument(v *domain.Version) (*VersionDocument, error) {
	pages := v.GetPages()
	pageDocs := make([]*PageDocument, 0, len(pages))

	for _, p := range pages {
		doc, err := ToPageDocument(p)

		if err != nil {
			return nil, err
		}

		pageDocs = append(pageDocs, doc)
	}

	return &VersionDocument{
		ID:          string(v.ID),
		FormID:      string(v.FormID),
		Version:     v.Version,
		Status:      string(v.Status),
		PublishedBy: v.PublishedBy,
		PublishedAt: v.PublishedAt,
		RetiredBy:   v.RetiredBy,
		RetiredAt:   v.RetiredAt,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
		Pages:       pageDocs,
	}, nil
}

func FromVersionDocument(v *VersionDocument) (*domain.Version, error) {
	version := domain.HydrateVersion(
		domain.VersionID(v.ID),
		domain.FormID(v.FormID),
		v.Version,
		domain.VersionStatus(v.Status),
		v.PublishedBy,
		v.PublishedAt,
		v.RetiredBy,
		v.RetiredAt,
		v.CreatedAt,
		v.UpdatedAt,
	)

	pages := make([]*domain.Page, 0, len(v.Pages))
	for _, p := range v.Pages {
		page, err := FromPageDocument(p)

		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	if err := version.AddPages(pages...); err != nil {
		return nil, err
	}

	return version, nil
}
