package domain

import (
	"errors"
	"maps"
	"slices"
	"time"

	"github.com/google/uuid"
)

type VersionStatus string

const (
	VersionStatusDraft   VersionStatus = "draft"
	VersionStatusActive  VersionStatus = "active"
	VersionStatusRetired VersionStatus = "retired"
)

var (
	ErrInvalidVersion      = errors.New("invalid version")
	ErrVersionLocked       = errors.New("version is locked")
	ErrPublishedByRequired = errors.New("publishedBy is required")
	ErrRetiredByRequired   = errors.New("retiredBy is required")
)

type VersionID string

type Version struct {
	ID          VersionID
	FormID      FormID
	Version     int
	Status      VersionStatus
	PublishedBy string
	PublishedAt time.Time
	RetiredBy   string
	RetiredAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	pages       map[int]*Page
}

func NewVersion(formID FormID, version int, status VersionStatus) (*Version, error) {
	v := &Version{
		ID:        VersionID(uuid.NewString()),
		FormID:    formID,
		Version:   version,
		Status:    status,
		pages:     make(map[int]*Page),
		CreatedAt: Now(),
	}

	// TODO: Implement domain specific validation.

	return v, nil
}

func HydrateVersion(
	id VersionID,
	formID FormID,
	version int,
	status VersionStatus,
	publishedBy string,
	publishedAt time.Time,
	retiredBy string,
	retiredAt,
	createdAt,
	updatedAt time.Time,
) *Version {
	return &Version{
		ID:          id,
		FormID:      formID,
		Version:     version,
		Status:      status,
		PublishedBy: publishedBy,
		PublishedAt: publishedAt,
		RetiredBy:   retiredBy,
		RetiredAt:   retiredAt,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		pages:       make(map[int]*Page),
	}
}

func (v *Version) GetPages() map[int]*Page {
	return v.pages
}

func (v *Version) GetPagesSlice() []*Page {
	positions := slices.Sorted(maps.Keys(v.pages))
	pages := make([]*Page, 0, len(v.pages))

	for _, position := range positions {
		pages = append(pages, v.pages[position])
	}

	return pages
}

func (v *Version) SetPages(pages ...*Page) error {
	if v == nil {
		return ErrInvalidVersion
	}

	if v.pages == nil {
		v.pages = make(map[int]*Page)
	}

	for _, page := range pages {
		position := page.GetPosition()
		_, exists := v.pages[position]

		if exists {
			return ErrDuplicatePosition
		}

		v.pages[position] = page
	}

	return nil
}

func (v *Version) UpdatePages(pages ...*Page) error {
	if v == nil {
		return ErrInvalidVersion
	}

	if v.Status != VersionStatusDraft {
		return ErrVersionLocked
	}

	old := v.pages
	v.pages = make(map[int]*Page)

	if err := v.SetPages(pages...); err != nil {
		v.pages = old
		return err
	}

	v.UpdatedAt = Now()

	return nil
}

func (v *Version) Publish(publishedBy string) error {
	if v == nil {
		return ErrInvalidVersion
	}

	if v.Status != VersionStatusDraft {
		return ErrVersionLocked
	}

	if publishedBy == "" {
		return ErrPublishedByRequired
	}

	now := Now()
	v.Status = VersionStatusActive
	v.PublishedBy = publishedBy
	v.PublishedAt = now
	v.UpdatedAt = now

	return nil
}

func (v *Version) Retire(retiredBy string) error {
	if v == nil {
		return ErrInvalidVersion
	}

	if v.Status != VersionStatusActive {
		return ErrVersionLocked
	}

	if retiredBy == "" {
		return ErrRetiredByRequired
	}

	now := Now()
	v.Status = VersionStatusRetired
	v.RetiredBy = retiredBy
	v.RetiredAt = now
	v.UpdatedAt = now

	return nil
}
