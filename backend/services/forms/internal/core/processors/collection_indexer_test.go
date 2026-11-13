package processors

import (
	"errors"
	"reflect"
	"testing"

	"sundance/backend/services/forms/internal/core/domain"
)

func Test_collectionIndexer_ResolveIndices(t *testing.T) {
	idx0 := 0
	idx1 := 1

	etmDynamic, _ := domain.NewElementTagMapping("elem1", "tv1", 1, false, nil)
	cDynamic0 := candidate{etm: etmDynamic, value: newTestSubmissionValue(t, "elem1", "val0", &idx0, nil)}
	cDynamic1 := candidate{etm: etmDynamic, value: newTestSubmissionValue(t, "elem1", "val1", &idx1, nil)}

	etmUsersName, _ := domain.NewElementTagMapping("users", "tv_app_name", 1, true, "Account Type")
	etmUsersVal, _ := domain.NewElementTagMapping("users", "tv_app_val", 1, true, "User")
	cUsersName := candidate{etm: etmUsersName, value: newTestSubmissionValue(t, "users", "Alice", &idx0, nil)}
	cUsersVal := candidate{etm: etmUsersVal, value: newTestSubmissionValue(t, "users", "Bob", &idx1, nil)}

	etmPrivName, _ := domain.NewElementTagMapping("privileged", "tv_app_name", 1, true, "Privilege")
	etmPrivVal, _ := domain.NewElementTagMapping("privileged", "tv_app_val", 1, false, nil)
	cPrivName := candidate{etm: etmPrivName, value: newTestSubmissionValue(t, "privileged", "Admin", nil, nil)}
	cPrivVal := candidate{etm: etmPrivVal, value: newTestSubmissionValue(t, "privileged", "Admin", nil, nil)}

	etmItemName, _ := domain.NewElementTagMapping("item_name", "tv_item_name", 1, false, nil)
	cItemName0 := candidate{etm: etmItemName, value: newTestSubmissionValue(t, "item_name", "Laptop", &idx0, nil)}
	cItemName1 := candidate{etm: etmItemName, value: newTestSubmissionValue(t, "item_name", "Mouse", &idx1, nil)}

	etmCurrency, _ := domain.NewElementTagMapping("item_currency", "tv_item_currency", 10, true, "USD")
	cCurrency := candidate{etm: etmCurrency, value: nil}

	type prescanEntry struct {
		keypath    string
		candidates []candidate
	}

	type step struct {
		keypath     string
		candidate   candidate
		expected    []int
		expectedErr error
	}

	tests := []struct {
		name    string
		prescan []prescanEntry
		steps   []step
	}{
		{
			name: "should delegate to entityCollection for entity collection with explicit indices",
			prescan: []prescanEntry{
				{keypath: "items[*].status", candidates: []candidate{cDynamic0, cDynamic1}},
			},
			steps: []step{
				{keypath: "items[*].status", candidate: cDynamic0, expected: []int{0}},
				{keypath: "items[*].status", candidate: cDynamic1, expected: []int{1}},
			},
		},
		{
			name: "should delegate to accumulatorCollection when multiple elements share a tag path",
			prescan: []prescanEntry{
				{keypath: "appFormData[*].name", candidates: []candidate{cUsersName, cPrivName}},
				{keypath: "appFormData[*].value", candidates: []candidate{cUsersVal, cPrivVal}},
			},
			steps: []step{
				{keypath: "appFormData[*].name", candidate: cUsersName, expected: []int{0}},
				{keypath: "appFormData[*].value", candidate: cUsersVal, expected: []int{0}},
				{keypath: "appFormData[*].name", candidate: cPrivName, expected: []int{1}},
				{keypath: "appFormData[*].value", candidate: cPrivVal, expected: []int{1}},
			},
		},
		{
			name: "should fan out static tags in entity collection across explicit rows",
			prescan: []prescanEntry{
				{keypath: "items[*].name", candidates: []candidate{cItemName0, cItemName1}},
				{keypath: "items[*].currency", candidates: []candidate{cCurrency}},
			},
			steps: []step{
				{keypath: "items[*].currency", candidate: cCurrency, expected: []int{0, 1}},
			},
		},
		{
			name:    "should yield ErrMissingCollectionIndex when keypath has no collection segment",
			prescan: nil,
			steps: []step{
				{keypath: "scalar.keypath", candidate: cDynamic0, expectedErr: ErrMissingCollectionIndex},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexer := newCollectionIndexer()
			for _, p := range tt.prescan {
				indexer.trackIndices(newTestTag(p.keypath), p.candidates)
			}

			for i, s := range tt.steps {
				got, err := indexer.resolveIndices(s.candidate, newTestTag(s.keypath))
				if s.expectedErr != nil {
					if !errors.Is(err, s.expectedErr) {
						t.Fatalf("step %d: expected error %v, got %v", i, s.expectedErr, err)
					}
					continue
				}

				if err != nil {
					t.Fatalf("step %d: unexpected error: %v", i, err)
				}

				if !reflect.DeepEqual(got, s.expected) {
					t.Fatalf("step %d: expected %v, got %v", i, s.expected, got)
				}
			}
		})
	}
}

func Test_collectionIndexer_isAccumulator(t *testing.T) {
	etmA, _ := domain.NewElementTagMapping("elem_a", "tv1", 1, false, nil)
	etmB, _ := domain.NewElementTagMapping("elem_b", "tv1", 1, false, nil)

	tests := []struct {
		name    string
		prescan map[string]map[string][]domain.ElementID
		root    string
		want    bool
	}{
		{
			name: "should return true if any tag has multiple contributing element IDs",
			prescan: map[string]map[string][]domain.ElementID{
				"appFormData[*]": {
					"appFormData[*].name":  {etmA.ElementID, etmB.ElementID},
					"appFormData[*].value": {etmA.ElementID},
				},
			},
			root: "appFormData[*]",
			want: true,
		},
		{
			name: "should return false if all tags have at most one contributing element ID",
			prescan: map[string]map[string][]domain.ElementID{
				"items[*]": {
					"items[*].name":  {etmA.ElementID},
					"items[*].price": {etmB.ElementID},
				},
			},
			root: "items[*]",
			want: false,
		},
		{
			name:    "should return false if root does not exist",
			prescan: nil,
			root:    "unknown[*]",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexer := newCollectionIndexer()
			if tt.prescan != nil {
				indexer.elementsPerTag = tt.prescan
			}

			got := indexer.isAccumulator(tt.root)
			if got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func Test_collectionIndexer_accumulatorCollection(t *testing.T) {
	etmA, _ := domain.NewElementTagMapping("elem_a", "tv1", 1, false, nil)
	etmB, _ := domain.NewElementTagMapping("elem_b", "tv1", 1, false, nil)

	cA := candidate{etm: etmA}
	cB := candidate{etm: etmB}

	tests := []struct {
		name     string
		root     string
		sequence []candidate
		expected [][]int
	}{
		{
			name:     "should assign sequential index per element and reuse for same element",
			root:     "appFormData[*]",
			sequence: []candidate{cA, cA, cB, cA, cB},
			expected: [][]int{{0}, {0}, {1}, {0}, {1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexer := newCollectionIndexer()
			for i, c := range tt.sequence {
				got, err := indexer.accumulatorCollection(tt.root, c)
				if err != nil {
					t.Fatalf("step %d: unexpected error: %v", i, err)
				}
				if !reflect.DeepEqual(got, tt.expected[i]) {
					t.Fatalf("step %d: expected %v, got %v", i, tt.expected[i], got)
				}
			}
		})
	}
}

func Test_collectionIndexer_entityCollection(t *testing.T) {
	idx0 := 0
	idx1 := 1

	etmDynamic, _ := domain.NewElementTagMapping("elem1", "tv1", 1, false, nil)
	cDynamic0 := candidate{etm: etmDynamic, value: newTestSubmissionValue(t, "elem1", "v0", &idx0, nil)}
	cDynamic1 := candidate{etm: etmDynamic, value: newTestSubmissionValue(t, "elem1", "v1", &idx1, nil)}

	etmStatic, _ := domain.NewElementTagMapping("elem_static", "tv_static", 10, true, "USD")
	cStatic := candidate{etm: etmStatic, value: nil}

	etmHetA, _ := domain.NewElementTagMapping("het_a", "tv_het", 1, false, nil)
	etmHetB, _ := domain.NewElementTagMapping("het_b", "tv_het", 1, false, nil)
	cHetA := candidate{etm: etmHetA, value: newTestSubmissionValue(t, "het_a", "A", nil, nil)}
	cHetB := candidate{etm: etmHetB, value: newTestSubmissionValue(t, "het_b", "B", nil, nil)}

	type evalStep struct {
		candidate candidate
		expected  []int
	}

	tests := []struct {
		name         string
		root         string
		explicitRows []int
		steps        []evalStep
	}{
		{
			name: "should return explicit collection index for dynamic candidates",
			root: "items[*]",
			steps: []evalStep{
				{candidate: cDynamic0, expected: []int{0}},
				{candidate: cDynamic1, expected: []int{1}},
			},
		},
		{
			name:         "should fan out static tag across all known explicit rows",
			root:         "items[*]",
			explicitRows: []int{0, 1},
			steps: []evalStep{
				{candidate: cStatic, expected: []int{0, 1}},
			},
		},
		{
			name:         "should return nil for static tag on empty entity collection",
			root:         "items[*]",
			explicitRows: nil,
			steps: []evalStep{
				{candidate: cStatic, expected: nil},
			},
		},
		{
			name: "should assign sequential index for heterogeneous candidates without explicit index",
			root: "items[*]",
			steps: []evalStep{
				{candidate: cHetA, expected: []int{0}},
				{candidate: cHetB, expected: []int{1}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexer := newCollectionIndexer()
			for _, r := range tt.explicitRows {
				indexer.addExplicitIndex(tt.root, r)
			}

			for i, s := range tt.steps {
				got, err := indexer.entityCollection(tt.root, s.candidate)
				if err != nil {
					t.Fatalf("step %d: unexpected error: %v", i, err)
				}

				if !reflect.DeepEqual(got, s.expected) {
					t.Fatalf("step %d: expected %v, got %v", i, s.expected, got)
				}
			}
		})
	}
}

func newTestTag(keyPath string) domain.Tag {
	return domain.Tag{KeyPath: keyPath}
}

func newTestSubmissionValue(t *testing.T, elementID domain.ElementID, value any, collectionIndex, instanceIndex *int) *domain.SubmissionValue {
	t.Helper()

	sv, err := domain.NewSubmissionValue(elementID, value, collectionIndex, instanceIndex)
	if err != nil {
		t.Fatalf("NewSubmissionValue(%q): unexpected error: %v", elementID, err)
	}

	return sv
}
