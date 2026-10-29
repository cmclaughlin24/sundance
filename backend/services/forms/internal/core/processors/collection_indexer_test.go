package processors

import (
	"reflect"
	"testing"

	"sundance/backend/services/forms/internal/core/domain"
)

func TestCollectionIndexer_TrackAndResolveIndices(t *testing.T) {
	tagKeyPath := "items[*].status"
	root := "items[*]"

	t.Run("explicit collection indices from candidates", func(t *testing.T) {
		indexer := newCollectionIndexer()
		idx0 := 0
		idx1 := 1

		etm, _ := domain.NewElementTagMapping("elem1", "tv1", 1, false, nil)
		candidates := []candidate{
			{etm: etm, value: domain.NewSubmissionValue("elem1", "val0", &idx0)},
			{etm: etm, value: domain.NewSubmissionValue("elem1", "val1", &idx1)},
		}

		indexer.trackIndices(tagKeyPath, candidates)

		if !reflect.DeepEqual(indexer.knownIndices[root], []int{0, 1}) {
			t.Fatalf("expected knownIndices to be [0, 1], got %v", indexer.knownIndices[root])
		}

		// Explicit index candidate resolves to its explicit index
		indices, err := indexer.resolveIndices(candidates[0], tagKeyPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(indices, []int{0}) {
			t.Fatalf("expected [0], got %v", indices)
		}
	})

	t.Run("static tag fans out across all known collection indices", func(t *testing.T) {
		indexer := newCollectionIndexer()
		idx0 := 0
		idx1 := 1

		etmDynamic, _ := domain.NewElementTagMapping("elem_name", "tv_name", 1, false, nil)
		dynamicCandidates := []candidate{
			{etm: etmDynamic, value: domain.NewSubmissionValue("elem_name", "item 1", &idx0)},
			{etm: etmDynamic, value: domain.NewSubmissionValue("elem_name", "item 2", &idx1)},
		}

		indexer.trackIndices("items[*].name", dynamicCandidates)

		// Static tag mapping for items[*].status with no submitted value
		etmStatic, _ := domain.NewElementTagMapping("elem_status", "tv_status", 10, true, "ACTIVE")
		staticCandidate := candidate{
			etm:   etmStatic,
			value: nil,
		}

		indices, err := indexer.resolveIndices(staticCandidate, "items[*].status")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(indices, []int{0, 1}) {
			t.Fatalf("expected static tag to fan out to [0, 1], got %v", indices)
		}
	})

	t.Run("static tag on empty collection resolves to nil", func(t *testing.T) {
		indexer := newCollectionIndexer()

		etmStatic, _ := domain.NewElementTagMapping("elem_status", "tv_status", 10, true, "ACTIVE")
		staticCandidate := candidate{
			etm:   etmStatic,
			value: nil,
		}

		indices, err := indexer.resolveIndices(staticCandidate, "items[*].status")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if indices != nil {
			t.Fatalf("expected nil indices for empty collection, got %v", indices)
		}
	})

	t.Run("heterogeneous dynamic elements without explicit index get sequential indices", func(t *testing.T) {
		indexer := newCollectionIndexer()

		etm1, _ := domain.NewElementTagMapping("elem_a", "tv_a", 1, false, nil)
		etm2, _ := domain.NewElementTagMapping("elem_b", "tv_b", 1, false, nil)

		candidates := []candidate{
			{etm: etm1, value: domain.NewSubmissionValue("elem_a", "A", nil)},
			{etm: etm2, value: domain.NewSubmissionValue("elem_b", "B", nil)},
		}

		indexer.trackIndices("records[*].field", candidates)

		if !reflect.DeepEqual(indexer.knownIndices["records[*]"], []int{0, 1}) {
			t.Fatalf("expected knownIndices [0, 1], got %v", indexer.knownIndices["records[*]"])
		}

		indices1, err := indexer.resolveIndices(candidates[0], "records[*].field")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(indices1, []int{0}) {
			t.Fatalf("expected [0], got %v", indices1)
		}

		indices2, err := indexer.resolveIndices(candidates[1], "records[*].field")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(indices2, []int{1}) {
			t.Fatalf("expected [1], got %v", indices2)
		}
	})

	t.Run("invalid tag path without collection segment returns ErrMissingCollectionIndex", func(t *testing.T) {
		indexer := newCollectionIndexer()
		etm, _ := domain.NewElementTagMapping("elem1", "tv1", 1, false, nil)
		c := candidate{etm: etm, value: domain.NewSubmissionValue("elem1", "val", nil)}

		_, err := indexer.resolveIndices(c, "scalar.path")
		if err != ErrMissingCollectionIndex {
			t.Fatalf("expected ErrMissingCollectionIndex, got %v", err)
		}
	})
}
