package domain

import (
	"errors"
	"slices"
	"time"

	"sundance/backend/pkg/common/validate"
)

type VersionStatus string

const (
	VersionStatusDraft   VersionStatus = "draft"
	VersionStatusActive  VersionStatus = "active"
	VersionStatusRetired VersionStatus = "retired"
)

var (
	ErrInvalidVersion       = errors.New("invalid version")
	ErrInvalidVersionStatus = errors.New("invalid version status")
	ErrVersionLocked        = errors.New("version is locked")
	ErrDuplicateVersion     = errors.New("duplicate version")
	ErrPublishedByRequired  = errors.New("publishedBy is required")
	ErrRetiredByRequired    = errors.New("retiredBy is required")
)

type VersionID string

type Version struct {
	ID          VersionID
	FormID      FormID        `validate:"required"`
	Version     int           `validate:"required,min=1"`
	Status      VersionStatus `validate:"required"`
	PublishedBy string
	PublishedAt time.Time
	RetiredBy   string
	RetiredAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	pages       PositionElements[*Page]
}

func NewVersion(formID FormID, version int, status VersionStatus) (*Version, error) {
	if !isValidVersionStatus(status) {
		return nil, ErrInvalidVersionStatus
	}

	v := &Version{
		ID:        VersionID(NewID()),
		FormID:    formID,
		Version:   version,
		Status:    status,
		pages:     make(PositionElements[*Page], 0),
		CreatedAt: Now(),
	}

	if err := validate.ValidateStruct(v); err != nil {
		return nil, err
	}

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
		pages:       make(PositionElements[*Page], 0),
	}
}

func (v *Version) FlatFields() []*Field {
	var fields []*Field

	for _, page := range v.pages {
		for _, section := range page.GetSections() {
			fields = append(fields, section.GetFields()...)
		}
	}

	return fields
}

func (v *Version) GetPages() []*Page {
	return v.pages
}

func (v *Version) AddPages(pages ...*Page) error {
	if v == nil {
		return ErrInvalidVersion
	}

	cpy := slices.Clone(v.pages)
	cpy = append(cpy, pages...)

	if ok := hasUniqueElements(cpy); !ok {
		return ErrDuplicatePosition
	}

	sortElements(cpy)
	v.pages = cpy

	return nil
}

func (v *Version) ReplacePages(pages ...*Page) error {
	if v == nil {
		return ErrInvalidVersion
	}

	if v.Status != VersionStatusDraft {
		return ErrVersionLocked
	}

	old := v.pages
	v.pages = make(PositionElements[*Page], 0)

	if err := v.AddPages(pages...); err != nil {
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

var isValidVersionStatus = validate.NewTypeValidator([]VersionStatus{
	VersionStatusDraft,
	VersionStatusActive,
	VersionStatusRetired,
})
