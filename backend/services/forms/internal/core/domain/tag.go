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
)

type TagID string

type TagNodeType string

const (
	TagNodeTypePrimitive TagNodeType = "primitive"
	TagNodeTypeObject    TagNodeType = "object"
	collectionSegment    string      = "[*]"
)

type Tag struct {
	ID            TagID
	TenantID      string
	KeyPath       string
	DisplayName   string
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
	return strings.Contains(t.KeyPath, collectionSegment)
}

var isTagValueKind = validate.NewTypeValidator([]TagNodeType{
	TagNodeTypePrimitive,
	TagNodeTypeObject,
})
