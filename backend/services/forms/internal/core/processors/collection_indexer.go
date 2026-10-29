package processors

import (
	"slices"
	"strings"
	"sundance/backend/services/forms/internal/core/domain"
)

type collectionIndexer struct {
	knownIndices   map[string][]int
	elementIndices map[string]map[domain.ElementID]int
}

func newCollectionIndexer() *collectionIndexer {
	return &collectionIndexer{
		knownIndices:   make(map[string][]int),
		elementIndices: make(map[string]map[domain.ElementID]int),
	}
}

func (idx *collectionIndexer) addKnownIndex(root string, index int) {
	if slices.Contains(idx.knownIndices[root], index) {
		return
	}

	idx.knownIndices[root] = append(idx.knownIndices[root], index)
	slices.Sort(idx.knownIndices[root])
}

func (idx *collectionIndexer) trackIndices(keypath string, candidates []candidate) {
	root := getCollectionRoot(keypath)
	if root == "" {
		return
	}

	for _, c := range candidates {
		if c.value != nil && c.value.CollectionIndex != nil {
			idx.addKnownIndex(root, *c.value.CollectionIndex)
			continue
		}

		if c.etm.HasStaticValue || c.value == nil || c.value.Value == nil {
			continue
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

		idx.addKnownIndex(root, index)
	}
}

func (idx *collectionIndexer) resolveIndices(c candidate, keypath string) ([]int, error) {
	if c.value != nil && c.value.CollectionIndex != nil {
		return []int{*c.value.CollectionIndex}, nil
	}

	root := getCollectionRoot(keypath)
	if root == "" {
		return nil, ErrMissingCollectionIndex
	}

	if c.etm.HasStaticValue {
		indices, exists := idx.knownIndices[root]
		if !exists || len(indices) == 0 {
			return nil, nil
		}
		return slices.Clone(indices), nil
	}

	if _, ok := idx.elementIndices[root]; !ok {
		idx.elementIndices[root] = make(map[domain.ElementID]int)
	}

	rootIndices := idx.elementIndices[root]
	index, exists := rootIndices[c.etm.ElementID]
	if !exists {
		index = len(rootIndices)
		rootIndices[c.etm.ElementID] = index
		idx.addKnownIndex(root, index)
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
