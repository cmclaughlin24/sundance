package domain

import (
	"errors"
	"time"

	"github.com/cmclaughlin24/sundance/forms/internal/types"
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
	ID            VersionID
	FormID        FormID
	Version       int
	Status        VersionStatus
	PublishedByID string
	PublishedAt   time.Time
	RetiredByID   string
	RetiredAt     time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Pages         map[int]*Page
}

func NewVersion(id VersionID, formID FormID, version int, status VersionStatus) (*Version, error) {
	v := &Version{
		ID:      id,
		FormID:  formID,
		Version: version,
		Status:  status,
		Pages:   make(map[int]*Page),
	}

	// TODO: Implement domain specific validation.

	return v, nil
}

func (v *Version) SetPages(pages ...*Page) error {
	if v == nil {
		return ErrInvalidVersion
	}

	if v.Pages == nil {
		v.Pages = make(map[int]*Page)
	}

	for _, page := range pages {
		_, exists := v.Pages[page.Position]

		if exists {
			return types.ErrDuplicatePosition
		}

		v.Pages[page.Position] = page
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

	v.Pages = make(map[int]*Page)

	return v.SetPages(pages...)
}

func (v *Version) Publish(publishedBy string, publishedAt time.Time) error {
	if v == nil {
		return ErrInvalidVersion
	}

	if v.Status != VersionStatusDraft {
		return ErrVersionLocked
	}

	if publishedBy == "" {
		return ErrPublishedByRequired
	}

	v.Status = VersionStatusActive
	v.PublishedByID = publishedBy
	v.PublishedAt = publishedAt

	return nil
}

func (v *Version) Retire(retiredBy string, retiredAt time.Time) error {
	if v == nil {
		return ErrInvalidVersion
	}

	if v.Status != VersionStatusActive {
		return ErrVersionLocked
	}

	if retiredBy == "" {
		return ErrRetiredByRequired
	}

	v.Status = VersionStatusRetired
	v.RetiredByID = retiredBy
	v.RetiredAt = retiredAt

	return nil
}
