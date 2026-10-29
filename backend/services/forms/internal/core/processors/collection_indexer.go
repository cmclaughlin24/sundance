package processors

import (
	"slices"
	"strings"
	"sundance/backend/services/forms/internal/core/domain"
)

type collectionIndexer struct {
	knownExplicitIndices map[string][]int
	elementsPerTag       map[string]map[string][]domain.ElementID
	elementIndices       map[string]map[domain.ElementID]int
}

func newCollectionIndexer() *collectionIndexer {
	return &collectionIndexer{
		knownExplicitIndices: make(map[string][]int),
		elementsPerTag:       make(map[string]map[string][]domain.ElementID),
		elementIndices:       make(map[string]map[domain.ElementID]int),
	}
}

func (idx *collectionIndexer) addExplicitIndex(root string, index int) {
	if slices.Contains(idx.knownExplicitIndices[root], index) {
		return
	}
	idx.knownExplicitIndices[root] = append(idx.knownExplicitIndices[root], index)
	slices.Sort(idx.knownExplicitIndices[root])
}

func (idx *collectionIndexer) trackIndices(keypath string, candidates []candidate) {
	root := getCollectionRoot(keypath)
	if root == "" {
		return
	}

	if _, ok := idx.elementsPerTag[root]; !ok {
		idx.elementsPerTag[root] = make(map[string][]domain.ElementID)
	}

	for _, c := range candidates {
		if !c.etm.HasStaticValue && c.value != nil && c.value.CollectionIndex != nil {
			idx.addExplicitIndex(root, *c.value.CollectionIndex)
		}

		elemList := idx.elementsPerTag[root][keypath]
		if !slices.Contains(elemList, c.etm.ElementID) {
			idx.elementsPerTag[root][keypath] = append(elemList, c.etm.ElementID)
		}
	}
}

func (idx *collectionIndexer) resolveIndices(c candidate, keypath string) ([]int, error) {
	root := getCollectionRoot(keypath)
	if root == "" {
		return nil, ErrMissingCollectionIndex
	}

	if idx.isAccumulator(root) {
		return idx.accumulatorCollection(root, c)
	}

	return idx.entityCollection(root, c)
}

func (idx *collectionIndexer) isAccumulator(root string) bool {
	for _, elemList := range idx.elementsPerTag[root] {
		if len(elemList) > 1 {
			return true
		}
	}
	return false
}

func (idx *collectionIndexer) accumulatorCollection(root string, c candidate) ([]int, error) {
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

func (idx *collectionIndexer) entityCollection(root string, c candidate) ([]int, error) {
	// a) Dynamic values with an explicit collection index belong to that native row.
	if !c.etm.HasStaticValue && c.value != nil && c.value.CollectionIndex != nil {
		return []int{*c.value.CollectionIndex}, nil
	}

	// b) Static tags in an entity collection fan out across all active explicit rows.
	if c.etm.HasStaticValue {
		indices, exists := idx.knownExplicitIndices[root]
		if !exists || len(indices) == 0 {
			return nil, nil // Empty collection -> 0 facts (no phantom rows)
		}
		return slices.Clone(indices), nil
	}

	// c) Fallback for dynamic heterogeneous elements without an explicit index
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
