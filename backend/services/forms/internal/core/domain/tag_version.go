package domain

import (
	"errors"
	"sundance/backend/pkg/common/validate"
	"time"
)

var (
	ErrDuplicateTagVersion = errors.New("duplicate tag version")
)

type TagVersionID string

type TagType string

type TagStatus string

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
	deprectatedAt,
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
		DeprecatedAt: deprectatedAt,
		PublishedAt:  publishedAt,
		RetiredAt:    retiredAt,
	}
}

func (tv *TagVersion) Deprecate() {
	tv.Status = TagStatusRetired
	tv.DeprecatedAt = Now()
}

func (tv *TagVersion) Publish() {
	tv.Status = TagStatusActive
	tv.PublishedAt = Now()
}

func (tv *TagVersion) Retire() {
	tv.Status = TagStatusRetired
	tv.RetiredAt = Now()
}
