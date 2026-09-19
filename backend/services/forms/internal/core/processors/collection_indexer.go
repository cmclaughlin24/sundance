package processors

import (
	"strings"
	"sundance/backend/services/forms/internal/core/domain"
)

type collectionIndexer struct {
	indices map[string]map[domain.ElementID]int
}

func newCollectionIndexer() *collectionIndexer {
	return &collectionIndexer{
		indices: make(map[string]map[domain.ElementID]int),
	}
}

func (idx *collectionIndexer) resolveIndex(c candidate, keypath string) (*int, error) {
	if c.value != nil && c.value.CollectionIndex != nil {
		return c.value.CollectionIndex, nil
	}

	root := getCollectionRoot(keypath)
	if root == "" {
		return nil, ErrMissingCollectionIndex
	}

	if _, ok := idx.indices[root]; !ok {
		idx.indices[root] = make(map[domain.ElementID]int)
	}

	rootIndices := idx.indices[root]
	index, exists := rootIndices[c.etm.ElementID]

	if !exists {
		index = len(rootIndices)
		rootIndices[c.etm.ElementID] = index
	}

	return &index, nil
}

func getCollectionRoot(keyPath string) string {
	i := strings.Index(keyPath, domain.TagCollectionSegment)
	if i == -1 {
		return ""
	}

	return keyPath[:i+len(domain.TagCollectionSegment)]
}
