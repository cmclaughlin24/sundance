package processors

import (
	"reflect"
	"testing"

	"sundance/backend/services/forms/internal/core/domain"
)

func TestCollectionIndexer_TrackAndResolveIndices(t *testing.T) {
	tagKeyPath := "items[*].status"

	t.Run("explicit collection indices from dynamic candidates", func(t *testing.T) {
		indexer := newCollectionIndexer()
		idx0 := 0
		idx1 := 1

		etm, _ := domain.NewElementTagMapping("elem1", "tv1", 1, false, nil)
		c0 := candidate{etm: etm, value: domain.NewSubmissionValue("elem1", "val0", &idx0)}
		c1 := candidate{etm: etm, value: domain.NewSubmissionValue("elem1", "val1", &idx1)}

		indices0, err := indexer.resolveIndices(c0, tagKeyPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(indices0, []int{0}) {
			t.Fatalf("expected [0], got %v", indices0)
		}

		indices1, err := indexer.resolveIndices(c1, tagKeyPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(indices1, []int{1}) {
			t.Fatalf("expected [1], got %v", indices1)
		}
	})

	t.Run("static tags ignore dynamic collection index from another collection and index per element", func(t *testing.T) {
		indexer := newCollectionIndexer()
		idx0 := 0
		idx1 := 1

		etmUsersName, _ := domain.NewElementTagMapping("users", "tv_app_name", 1, true, "Account Type")
		etmUsersVal, _ := domain.NewElementTagMapping("users", "tv_app_val", 1, true, "User")

		cUsers0 := candidate{etm: etmUsersName, value: domain.NewSubmissionValue("users", "Alice", &idx0)}
		cUsers1 := candidate{etm: etmUsersVal, value: domain.NewSubmissionValue("users", "Bob", &idx1)}

		etmPrivName, _ := domain.NewElementTagMapping("privileged", "tv_app_name", 1, true, "Privilege")
		etmPrivVal, _ := domain.NewElementTagMapping("privileged", "tv_app_val", 1, false, nil)

		cPrivName := candidate{etm: etmPrivName, value: domain.NewSubmissionValue("privileged", "Admin", nil)}
		cPrivVal := candidate{etm: etmPrivVal, value: domain.NewSubmissionValue("privileged", "Admin", nil)}

		idxUsersName, err := indexer.resolveIndices(cUsers0, "appFormData[*].name")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(idxUsersName, []int{0}) {
			t.Fatalf("expected users name index to be [0], got %v", idxUsersName)
		}

		idxUsersVal, err := indexer.resolveIndices(cUsers1, "appFormData[*].value")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(idxUsersVal, []int{0}) {
			t.Fatalf("expected users val index to be [0], got %v", idxUsersVal)
		}

		// Both 'privileged' mappings should resolve to index 1 in appFormData[*]
		idxPrivName, err := indexer.resolveIndices(cPrivName, "appFormData[*].name")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(idxPrivName, []int{1}) {
			t.Fatalf("expected privileged name index to be [1], got %v", idxPrivName)
		}

		idxPrivVal, err := indexer.resolveIndices(cPrivVal, "appFormData[*].value")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(idxPrivVal, []int{1}) {
			t.Fatalf("expected privileged val index to be [1], got %v", idxPrivVal)
		}
	})

	t.Run("heterogeneous dynamic elements without explicit index get sequential indices", func(t *testing.T) {
		indexer := newCollectionIndexer()

		etm1, _ := domain.NewElementTagMapping("elem_a", "tv_a", 1, false, nil)
		etm2, _ := domain.NewElementTagMapping("elem_b", "tv_b", 1, false, nil)

		c1 := candidate{etm: etm1, value: domain.NewSubmissionValue("elem_a", "A", nil)}
		c2 := candidate{etm: etm2, value: domain.NewSubmissionValue("elem_b", "B", nil)}

		indices1, err := indexer.resolveIndices(c1, "records[*].field")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(indices1, []int{0}) {
			t.Fatalf("expected [0], got %v", indices1)
		}

		indices2, err := indexer.resolveIndices(c2, "records[*].field")
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
