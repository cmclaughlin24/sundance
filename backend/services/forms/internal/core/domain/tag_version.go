package domain

import (
	"errors"
	"fmt"
	"sundance/backend/pkg/common/validate"
	"time"
)

var (
	ErrInvalidTagVersion       = errors.New("invalid tag version")
	ErrDuplicateTagVersion     = errors.New("duplicate tag version")
	ErrInvalidTagVersionStatus = errors.New("invalid tag version status invariant")
)

type TagVersionID string

type TagStatus string

type TagType string

const (
	TagStatusDraft      TagStatus = "draft"
	TagStatusActive     TagStatus = "active"
	TagStatusDeprecated TagStatus = "deprecated"
	TagStatusRetired    TagStatus = "retired"
)

type TagVersion struct {
	ID           TagVersionID
	TagID        TagID
	Version      int
	Type         TagType
	Status       TagStatus
	CreatedAt    time.Time
	DeprecatedAt time.Time
	PublishedAt  time.Time
	RetiredAt    time.Time
}

func NewTagVersion(tagID TagID, version int, tagType TagType) (*TagVersion, error) {
	ctv := &TagVersion{
		ID:        TagVersionID(NewID()),
		TagID:     tagID,
		Version:   version,
		Status:    TagStatusDraft,
		Type:      tagType,
		CreatedAt: Now(),
	}

	if err := validate.ValidateStruct(ctv); err != nil {
		return nil, err
	}

	return ctv, nil
}

func HydrateTagVersion(
	id TagVersionID,
	tagID TagID,
	version int,
	tagType TagType,
	status TagStatus,
	createdAt,
	deprecatedAt,
	publishedAt,
	retiredAt time.Time,
) *TagVersion {
	return &TagVersion{
		ID:           id,
		TagID:        tagID,
		Version:      version,
		Type:         tagType,
		Status:       status,
		CreatedAt:    createdAt,
		DeprecatedAt: deprecatedAt,
		PublishedAt:  publishedAt,
		RetiredAt:    retiredAt,
	}
}

func (tv *TagVersion) Deprecate() error {
	if tv == nil {
		return ErrInvalidTagVersion
	}

	if tv.Status != TagStatusActive {
		return fmt.Errorf("cannot deprecate tag in status %s: %w", tv.Status, ErrInvalidTagVersionStatus)
	}

	tv.Status = TagStatusDeprecated
	tv.DeprecatedAt = Now()

	return nil
}

func (tv *TagVersion) Publish() error {
	if tv == nil {
		return ErrInvalidTagVersion
	}

	if tv.Status != TagStatusDraft {
		return fmt.Errorf("cannot publish tag in status %s: %w", tv.Status, ErrInvalidTagVersionStatus)
	}

	tv.Status = TagStatusActive
	tv.PublishedAt = Now()

	return nil
}

func (tv *TagVersion) Retire() error {
	if tv == nil {
		return ErrInvalidTagVersion
	}

	if tv.Status != TagStatusDeprecated {
		return fmt.Errorf("cannot retire tag in status %s: %w", tv.Status, ErrInvalidTagVersionStatus)
	}

	tv.Status = TagStatusRetired
	tv.RetiredAt = Now()

	return nil
}
