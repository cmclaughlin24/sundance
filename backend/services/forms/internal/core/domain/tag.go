package domain

import (
	"errors"
	"fmt"
	"strings"
	"sundance/backend/pkg/common/validate"
	"time"
)

var (
	ErrInvalidTag                  = errors.New("invalid tag")
	ErrInvalidTagNodeType          = fmt.Errorf("%w; node type", ErrInvalidTag)
	ErrNodeTypeObjectPrimitiveType = fmt.Errorf("%w; node type object cannot have a primitive type", ErrInvalidTag)
	ErrCollectionCount             = fmt.Errorf("%w; tag cannot contain more than %d collection segments", ErrInvalidTag, maxCollectionDepth)
	ErrTagHasVersions              = fmt.Errorf("%w; tag has at least one version", ErrInvalidTag)
)

type TagID string

type TagNodeType string

const (
	TagNodeTypePrimitive TagNodeType = "primitive"
	TagNodeTypeObject    TagNodeType = "object"
	TagCollectionSegment string      = "[*]"
	pathSeparator        string      = "."
	maxCollectionDepth               = 2
)

type Tag struct {
	ID            TagID
	TenantID      string `validate:"required"`
	KeyPath       string `validate:"required"`
	DisplayName   string `validate:"required"`
	NodeType      TagNodeType
	PrimitiveType *TagPrimitiveType
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewTag(tenantID, keyPath, displayName string, nodeType TagNodeType, primitiveType *TagPrimitiveType) (*Tag, error) {
	if !isTagValueKind(nodeType) {
		return nil, ErrInvalidTagNodeType
	}

	if nodeType == TagNodeTypeObject && primitiveType != nil {
		return nil, ErrNodeTypeObjectPrimitiveType
	}

	ct := &Tag{
		ID:            TagID(NewID()),
		TenantID:      tenantID,
		KeyPath:       keyPath,
		DisplayName:   displayName,
		NodeType:      nodeType,
		PrimitiveType: primitiveType,
		CreatedAt:     Now(),
	}

	if ct.CollectionDepth() > maxCollectionDepth {
		return nil, ErrCollectionCount
	}

	if err := validate.ValidateStruct(ct); err != nil {
		return nil, err
	}

	return ct, nil
}

func HydrateTag(
	id TagID,
	tenantID,
	keyPath,
	displayName string,
	nodeType TagNodeType,
	primitiveType *TagPrimitiveType,
	createdAt,
	updatedAt time.Time,
) *Tag {
	return &Tag{
		ID:            id,
		TenantID:      tenantID,
		KeyPath:       keyPath,
		DisplayName:   displayName,
		NodeType:      nodeType,
		PrimitiveType: primitiveType,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (t *Tag) Update(displayName string) error {
	if t == nil {
		return ErrInvalidTag
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

func (t *Tag) HasCollectionAncestor() bool {
	return strings.Contains(t.KeyPath, TagCollectionSegment)
}

func (t *Tag) CollectionSegments() []string {
	segments := make([]string, 0, maxCollectionDepth)
	offset := 0

	for {
		i := strings.Index(t.KeyPath[offset:], TagCollectionSegment)
		if i == -1 {
			return segments
		}

		end := offset + i + len(TagCollectionSegment)
		segments = append(segments, t.KeyPath[:end])
		offset = end
	}
}

func (t *Tag) CollectionDepth() int {
	return strings.Count(t.KeyPath, TagCollectionSegment)
}

func (t *Tag) CollectionRoot() string {
	segments := t.CollectionSegments()
	if len(segments) == 0 {
		return ""
	}
	return segments[0]
}

var isTagValueKind = validate.NewTypeValidator([]TagNodeType{
	TagNodeTypePrimitive,
	TagNodeTypeObject,
})
