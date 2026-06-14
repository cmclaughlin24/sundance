package domain

import (
	"errors"
	"slices"
	"time"

	"sundance/backend/pkg/common/validate"
)

type FormVersionStatus string

const (
	FormVersionStatusDraft   FormVersionStatus = "draft"
	FormVersionStatusActive  FormVersionStatus = "active"
	FormVersionStatusRetired FormVersionStatus = "retired"
)

var (
	ErrInvalidVersion       = errors.New("invalid form version")
	ErrInvalidVersionStatus = errors.New("invalid version status")
	ErrVersionLocked        = errors.New("version is locked")
	ErrDuplicateVersion     = errors.New("duplicate version")
	ErrPublishedByRequired  = errors.New("publishedBy is required")
	ErrRetiredByRequired    = errors.New("retiredBy is required")
)

type FormVersionID string

type FormVersion struct {
	ID          FormVersionID
	FormID      FormID            `validate:"required"`
	Version     int               `validate:"required,min=1"`
	Status      FormVersionStatus `validate:"required"`
	PublishedBy string
	PublishedAt time.Time
	RetiredBy   string
	RetiredAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	pages       PositionElements[*Page]
}

func NewFormVersion(formID FormID, version int, status FormVersionStatus) (*FormVersion, error) {
	if !isValidFormVersionStatus(status) {
		return nil, ErrInvalidVersionStatus
	}

	v := &FormVersion{
		ID:        FormVersionID(NewID()),
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

func HydrateFormVersion(
	id FormVersionID,
	formID FormID,
	version int,
	status FormVersionStatus,
	publishedBy string,
	publishedAt time.Time,
	retiredBy string,
	retiredAt,
	createdAt,
	updatedAt time.Time,
) *FormVersion {
	return &FormVersion{
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

func (v *FormVersion) FlatFields() []*Field {
	var fields []*Field

	for _, page := range v.pages {
		for _, section := range page.GetSections() {
			fields = append(fields, section.GetFields()...)
		}
	}

	return fields
}

func (v *FormVersion) GetPages() []*Page {
	return v.pages
}

func (v *FormVersion) GetPage(pageID PageID) *Page {
	idx := slices.IndexFunc(v.pages, func(p *Page) bool {
		return p.ID == pageID
	})

	if idx == -1 {
		return nil
	}

	return v.pages[idx]
}

func (v *FormVersion) AddPages(pages ...*Page) error {
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

func (v *FormVersion) ReplacePages(pages ...*Page) error {
	if v == nil {
		return ErrInvalidVersion
	}

	if v.Status != FormVersionStatusDraft {
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

func (v *FormVersion) Publish(publishedBy string) error {
	if v == nil {
		return ErrInvalidVersion
	}

	if v.Status != FormVersionStatusDraft {
		return ErrVersionLocked
	}

	if publishedBy == "" {
		return ErrPublishedByRequired
	}

	now := Now()
	v.Status = FormVersionStatusActive
	v.PublishedBy = publishedBy
	v.PublishedAt = now
	v.UpdatedAt = now

	return nil
}

func (v *FormVersion) Retire(retiredBy string) error {
	if v == nil {
		return ErrInvalidVersion
	}

	if v.Status != FormVersionStatusActive {
		return ErrVersionLocked
	}

	if retiredBy == "" {
		return ErrRetiredByRequired
	}

	now := Now()
	v.Status = FormVersionStatusRetired
	v.RetiredBy = retiredBy
	v.RetiredAt = now
	v.UpdatedAt = now

	return nil
}

var isValidFormVersionStatus = validate.NewTypeValidator([]FormVersionStatus{
	FormVersionStatusDraft,
	FormVersionStatusActive,
	FormVersionStatusRetired,
})
