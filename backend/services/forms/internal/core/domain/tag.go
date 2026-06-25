package domain

import (
	"errors"
	"sundance/backend/pkg/common/validate"
	"time"
)

var (
	ErrInvalidCanonicalTag = errors.New("invalid tag")
	ErrInvalidTagValueKind = errors.New("invalid tag value kind")
)

type TagID string

type TagValueKind string

const (
	TagValueKindPrimitive TagValueKind = "primitive"
	TagValueKindObject    TagValueKind = "object"
)

type Tag struct {
	ID           TagID
	TenantID     string
	Key          string
	DisplayName  string
	ValueKind    TagValueKind
	IsCollection bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewTag(tenantID, key, displayName string, valueKind TagValueKind, isCollection bool) (*Tag, error) {
	if !isTagValueKind(valueKind) {
		return nil, ErrInvalidTagValueKind
	}

	ct := &Tag{
		ID:           TagID(NewID()),
		TenantID:     tenantID,
		Key:          key,
		DisplayName:  displayName,
		ValueKind:    valueKind,
		IsCollection: isCollection,
		CreatedAt:    Now(),
	}

	if err := validate.ValidateStruct(ct); err != nil {
		return nil, err
	}

	return ct, nil
}

func HydrateTag(
	id TagID,
	tenantID,
	key,
	displayName string,
	valueKind TagValueKind,
	isCollection bool,
	createdAt,
	updatedAt time.Time,
) *Tag {
	return &Tag{
		ID:           id,
		TenantID:     tenantID,
		Key:          key,
		DisplayName:  displayName,
		ValueKind:    valueKind,
		IsCollection: isCollection,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

func (t *Tag) Update(displayName string) error {
	if t == nil {
		return ErrInvalidCanonicalTag
	}

	cpy := *t
	cpy.DisplayName = displayName

	if err := validate.ValidateStruct(cpy); err != nil {
		return err
	}

	*t = cpy
	t.UpdatedAt = Now()

	return nil
}

var isTagValueKind = validate.NewTypeValidator([]TagValueKind{
	TagValueKindPrimitive,
	TagValueKindObject,
})
