package processors

import (
	"strings"
	"sundance/backend/services/forms/internal/core/domain"
)

type collectionIndexer struct {
	elementIndices map[string]map[domain.ElementID]int
}

func newCollectionIndexer() *collectionIndexer {
	return &collectionIndexer{
		elementIndices: make(map[string]map[domain.ElementID]int),
	}
}

func (idx *collectionIndexer) trackIndices(keypath string, candidates []candidate) {
}

func (idx *collectionIndexer) resolveIndices(c candidate, keypath string) ([]int, error) {
	root := getCollectionRoot(keypath)
	if root == "" {
		return nil, ErrMissingCollectionIndex
	}

	if !c.etm.HasStaticValue && c.value != nil && c.value.CollectionIndex != nil {
		return []int{*c.value.CollectionIndex}, nil
	}

	if _, ok := idx.elementIndices[root]; !ok {
		idx.elementIndices[root] = make(map[domain.ElementID]int)
	}

	rootIndices := idx.elementIndices[root]
	index, exists := rootIndices[c.etm.ElementID]
	if !exists {
		index = len(rootIndices)
		rootIndices[c.etm.ElementID] = index
	}

	return []int{index}, nil
}

func getCollectionRoot(keyPath string) string {
	i := strings.Index(keyPath, domain.TagCollectionSegment)
	if i == -1 {
		return ""
	}

	return keyPath[:i+len(domain.TagCollectionSegment)]
}
