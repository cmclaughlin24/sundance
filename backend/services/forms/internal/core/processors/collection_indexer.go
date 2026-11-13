package processors

import (
	"errors"
	"fmt"
	"slices"

	"sundance/backend/services/forms/internal/core/domain"
)

var (
	ErrCollectionDepthMismatch    = errors.New("collection depth mismatch")
	ErrUnsupportedCollectionDepth = errors.New("unsupported collection depth")
)

type rowKey struct {
	elementId    domain.ElementID
	sectionIndex int // -1 when not in a repeatable section
}

type collectionIndexer struct {
	knownPaths     map[string][][]int
	elementsPerTag map[string]map[string][]domain.ElementID
	rows           map[string]map[rowKey]int
}

func newCollectionIndexer() *collectionIndexer {
	return &collectionIndexer{
		knownPaths:     make(map[string][][]int),
		elementsPerTag: make(map[string]map[string][]domain.ElementID),
		rows:           make(map[string]map[rowKey]int),
	}
}

func (idx *collectionIndexer) trackIndices(t domain.Tag, candidates []candidate) {
	segments := t.CollectionSegments()
	if len(segments) == 0 {
		return
	}

	root := segments[0]
	if idx.elementsPerTag[root] == nil {
		idx.elementsPerTag[root] = make(map[string][]domain.ElementID)
	}

	for _, c := range candidates {
		if path := c.value.IndexPath(); !c.etm.HasStaticValue && len(path) == len(segments) {
			for depth, segment := range segments {
				idx.addKnownPath(segment, path[:depth+1])
			}
		}

		if elements := idx.elementsPerTag[root][t.KeyPath]; !slices.Contains(elements, c.etm.ElementID) {
			idx.elementsPerTag[root][t.KeyPath] = append(elements, c.etm.ElementID)
		}
	}
}

func (idx *collectionIndexer) resolvePaths(c candidate, t domain.Tag) ([][]int, error) {
	segments := t.CollectionSegments()
	if len(segments) == 0 {
		return nil, ErrMissingCollectionIndex
	}

	if idx.isAccumulator(segments[0]) {
		return idx.accumulatorCollection(segments, c, t.KeyPath)
	}

	return idx.entityCollection(segments, c, t.KeyPath)
}

func (idx *collectionIndexer) accumulatorCollection(segments []string, c candidate, keyPath string) ([][]int, error) {
	if len(segments) > 1 {
		return nil, fmt.Errorf("%w; accumulator tag=%s", ErrUnsupportedCollectionDepth, keyPath)
	}

	return [][]int{{idx.assignRow(segments[0], c)}}, nil
}

func (idx *collectionIndexer) entityCollection(segments []string, c candidate, keyPath string) ([][]int, error) {
	depth := len(segments)

	// a) Static values fan out across every known path at the tag's depth.
	if c.etm.HasStaticValue {
		return slices.Clone(idx.knownPaths[segments[depth-1]]), nil
	}

	// b) Dynamic values with explicit indices map to their native path.
	switch path := c.value.IndexPath(); {
	case len(path) == depth:
		return [][]int{path}, nil
	case len(path) > 0:
		return nil, fmt.Errorf("%w; tag=%s depth=%d indices=%d", ErrCollectionDepthMismatch, keyPath, depth, len(path))
	case depth > 1:
		return nil, fmt.Errorf("%w; tag=%s", ErrMissingCollectionIndex, keyPath)
	}

	// c) Fallback for dynamic values without indices: one row per element.
	return [][]int{{idx.assignRow(segments[0], c)}}, nil
}

func (idx *collectionIndexer) isAccumulator(root string) bool {
	for _, elements := range idx.elementsPerTag[root] {
		if len(elements) > 1 {
			return true
		}
	}
	return false
}

func (idx *collectionIndexer) assignRow(root string, c candidate) int {
	if idx.rows[root] == nil {
		idx.rows[root] = make(map[rowKey]int)
	}

	key := rowKey{elementId: c.etm.ElementID, sectionIndex: -1}
	if c.value != nil && c.value.SectionIndex != nil {
		key.sectionIndex = *c.value.SectionIndex
	}

	row, ok := idx.rows[root][key]
	if !ok {
		row = len(idx.rows[root])
		idx.rows[root][key] = row
	}
	return row
}

func (idx *collectionIndexer) addKnownPath(segment string, path []int) {
	if slices.ContainsFunc(idx.knownPaths[segment], func(p []int) bool { return slices.Equal(p, path) }) {
		return
	}

	idx.knownPaths[segment] = append(idx.knownPaths[segment], slices.Clone(path))
	slices.SortFunc(idx.knownPaths[segment], slices.Compare[[]int])
}
