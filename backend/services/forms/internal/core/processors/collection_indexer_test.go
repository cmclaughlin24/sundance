package processors

import (
	"errors"
	"reflect"
	"testing"

	"sundance/backend/services/forms/internal/core/domain"
)

func Test_collectionIndexer_resolvePaths(t *testing.T) {
	etmDynamic := newTestElementTagMapping(t, "elem1", "tv1", 1, false, nil)
	cDynamic0 := candidate{etm: etmDynamic, value: newTestSubmissionValue(t, "elem1", "val0", nil), path: []int{0}}
	cDynamic1 := candidate{etm: etmDynamic, value: newTestSubmissionValue(t, "elem1", "val1", nil), path: []int{1}}
	cDynamicNoPath := candidate{etm: etmDynamic, value: newTestSubmissionValue(t, "elem1", "val", nil)}
	cDynamicDeep := candidate{etm: etmDynamic, value: newTestSubmissionValue(t, "elem1", "val", nil), path: []int{0, 1}}

	etmUsersName := newTestElementTagMapping(t, "users", "tv_app_name", 1, true, "Account Type")
	etmUsersVal := newTestElementTagMapping(t, "users", "tv_app_val", 1, true, "User")
	cUsersName := candidate{etm: etmUsersName, value: newTestSubmissionValue(t, "users", "Alice", nil)}
	cUsersVal := candidate{etm: etmUsersVal, value: newTestSubmissionValue(t, "users", "Alice", nil)}

	etmPrivName := newTestElementTagMapping(t, "privileged", "tv_app_name", 1, true, "Privilege")
	etmPrivVal := newTestElementTagMapping(t, "privileged", "tv_app_val", 1, false, nil)
	cPrivName := candidate{etm: etmPrivName, value: newTestSubmissionValue(t, "privileged", "Admin", nil)}
	cPrivVal := candidate{etm: etmPrivVal, value: newTestSubmissionValue(t, "privileged", "Admin", nil)}

	etmItemName := newTestElementTagMapping(t, "item_name", "tv_item_name", 1, false, nil)
	cItemName0 := candidate{etm: etmItemName, value: newTestSubmissionValue(t, "item_name", "Laptop", nil), path: []int{0}}
	cItemName1 := candidate{etm: etmItemName, value: newTestSubmissionValue(t, "item_name", "Mouse", nil), path: []int{1}}

	etmCurrency := newTestElementTagMapping(t, "item_currency", "tv_item_currency", 10, true, "USD")
	cCurrency := candidate{etm: etmCurrency}

	etmLineSku := newTestElementTagMapping(t, "line_sku", "tv_line_sku", 1, false, nil)
	cLineSku00 := candidate{etm: etmLineSku, value: newTestSubmissionValue(t, "line_sku", "A", nil), path: []int{0, 0}}
	cLineSku01 := candidate{etm: etmLineSku, value: newTestSubmissionValue(t, "line_sku", "B", nil), path: []int{0, 1}}
	cLineSku10 := candidate{etm: etmLineSku, value: newTestSubmissionValue(t, "line_sku", "C", nil), path: []int{1, 0}}

	etmLineUnit := newTestElementTagMapping(t, "line_unit", "tv_line_unit", 1, true, "EA")
	cLineUnit := candidate{etm: etmLineUnit}

	etmAccA := newTestElementTagMapping(t, "acc_a", "tv_acc", 1, false, nil)
	etmAccB := newTestElementTagMapping(t, "acc_b", "tv_acc", 1, false, nil)
	cAccA := candidate{etm: etmAccA, value: newTestSubmissionValue(t, "acc_a", "A", nil), path: []int{0, 0}}
	cAccB := candidate{etm: etmAccB, value: newTestSubmissionValue(t, "acc_b", "B", nil), path: []int{0, 0}}

	type prescanEntry struct {
		keypath    string
		candidates []candidate
	}

	type step struct {
		keypath     string
		candidate   candidate
		expected    [][]int
		expectedErr error
	}

	tests := []struct {
		name    string
		prescan []prescanEntry
		steps   []step
	}{
		{
			name: "should map dynamic candidates in an entity collection to their native path",
			prescan: []prescanEntry{
				{keypath: "items[*].status", candidates: []candidate{cDynamic0, cDynamic1}},
			},
			steps: []step{
				{keypath: "items[*].status", candidate: cDynamic0, expected: [][]int{{0}}},
				{keypath: "items[*].status", candidate: cDynamic1, expected: [][]int{{1}}},
			},
		},
		{
			name: "should assign one row per element when multiple elements share a tag path",
			prescan: []prescanEntry{
				{keypath: "appFormData[*].name", candidates: []candidate{cUsersName, cPrivName}},
				{keypath: "appFormData[*].value", candidates: []candidate{cUsersVal, cPrivVal}},
			},
			steps: []step{
				{keypath: "appFormData[*].name", candidate: cUsersName, expected: [][]int{{0}}},
				{keypath: "appFormData[*].value", candidate: cUsersVal, expected: [][]int{{0}}},
				{keypath: "appFormData[*].name", candidate: cPrivName, expected: [][]int{{1}}},
				{keypath: "appFormData[*].value", candidate: cPrivVal, expected: [][]int{{1}}},
			},
		},
		{
			name: "should fan out static tags in an entity collection across known paths",
			prescan: []prescanEntry{
				{keypath: "items[*].name", candidates: []candidate{cItemName0, cItemName1}},
				{keypath: "items[*].currency", candidates: []candidate{cCurrency}},
			},
			steps: []step{
				{keypath: "items[*].currency", candidate: cCurrency, expected: [][]int{{0}, {1}}},
			},
		},
		{
			name: "should map nested dynamic candidates to their native path",
			prescan: []prescanEntry{
				{keypath: "orders[*].lines[*].sku", candidates: []candidate{cLineSku00, cLineSku01, cLineSku10}},
			},
			steps: []step{
				{keypath: "orders[*].lines[*].sku", candidate: cLineSku01, expected: [][]int{{0, 1}}},
				{keypath: "orders[*].lines[*].sku", candidate: cLineSku10, expected: [][]int{{1, 0}}},
			},
		},
		{
			name: "should fan out nested static tags across known paths at the tag's depth",
			prescan: []prescanEntry{
				{keypath: "orders[*].lines[*].sku", candidates: []candidate{cLineSku10, cLineSku00, cLineSku01}},
				{keypath: "orders[*].lines[*].unit", candidates: []candidate{cLineUnit}},
			},
			steps: []step{
				{keypath: "orders[*].lines[*].unit", candidate: cLineUnit, expected: [][]int{{0, 0}, {0, 1}, {1, 0}}},
			},
		},
		{
			name: "should yield ErrCollectionDepthMismatch when path length differs from tag depth",
			prescan: []prescanEntry{
				{keypath: "items[*].status", candidates: []candidate{cDynamicDeep}},
			},
			steps: []step{
				{keypath: "items[*].status", candidate: cDynamicDeep, expectedErr: ErrCollectionDepthMismatch},
			},
		},
		{
			name: "should yield ErrMissingCollectionIndex for nested tag without a path",
			prescan: []prescanEntry{
				{keypath: "orders[*].lines[*].sku", candidates: []candidate{cDynamicNoPath}},
			},
			steps: []step{
				{keypath: "orders[*].lines[*].sku", candidate: cDynamicNoPath, expectedErr: ErrMissingCollectionIndex},
			},
		},
		{
			name: "should yield ErrUnsupportedCollectionDepth for nested accumulator tags",
			prescan: []prescanEntry{
				{keypath: "acc[*].lines[*].value", candidates: []candidate{cAccA, cAccB}},
			},
			steps: []step{
				{keypath: "acc[*].lines[*].value", candidate: cAccA, expectedErr: ErrUnsupportedCollectionDepth},
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
				got, err := indexer.resolvePaths(s.candidate, newTestTag(s.keypath))
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

func Test_collectionIndexer_trackIndices(t *testing.T) {
	etmDynamic := newTestElementTagMapping(t, "elem1", "tv1", 1, false, nil)
	etmOther := newTestElementTagMapping(t, "elem2", "tv1", 1, false, nil)
	etmStatic := newTestElementTagMapping(t, "elem_static", "tv_static", 1, true, "USD")

	tests := []struct {
		name               string
		keypath            string
		candidates         []candidate
		wantKnownPaths     map[string][][]int
		wantElementsPerTag map[string]map[string][]domain.ElementID
	}{
		{
			name:    "should register known paths for dynamic candidates with matching depth",
			keypath: "items[*].name",
			candidates: []candidate{
				{etm: etmDynamic, path: []int{1}},
				{etm: etmDynamic, path: []int{0}},
				{etm: etmDynamic, path: []int{1}},
			},
			wantKnownPaths: map[string][][]int{
				"items[*]": {{0}, {1}},
			},
			wantElementsPerTag: map[string]map[string][]domain.ElementID{
				"items[*]": {"items[*].name": {"elem1"}},
			},
		},
		{
			name:    "should register every prefix of nested paths",
			keypath: "orders[*].lines[*].sku",
			candidates: []candidate{
				{etm: etmDynamic, path: []int{0, 1}},
				{etm: etmDynamic, path: []int{1, 0}},
			},
			wantKnownPaths: map[string][][]int{
				"orders[*]":          {{0}, {1}},
				"orders[*].lines[*]": {{0, 1}, {1, 0}},
			},
			wantElementsPerTag: map[string]map[string][]domain.ElementID{
				"orders[*]": {"orders[*].lines[*].sku": {"elem1"}},
			},
		},
		{
			name:    "should not register paths for static candidates or mismatched depths",
			keypath: "items[*].name",
			candidates: []candidate{
				{etm: etmStatic, path: []int{0}},
				{etm: etmDynamic, path: []int{0, 1}},
				{etm: etmOther},
			},
			wantKnownPaths: map[string][][]int{},
			wantElementsPerTag: map[string]map[string][]domain.ElementID{
				"items[*]": {"items[*].name": {"elem_static", "elem1", "elem2"}},
			},
		},
		{
			name:               "should ignore tags without a collection segment",
			keypath:            "scalar.keypath",
			candidates:         []candidate{{etm: etmDynamic, path: []int{0}}},
			wantKnownPaths:     map[string][][]int{},
			wantElementsPerTag: map[string]map[string][]domain.ElementID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexer := newCollectionIndexer()
			indexer.trackIndices(newTestTag(tt.keypath), tt.candidates)

			if !reflect.DeepEqual(indexer.knownPaths, tt.wantKnownPaths) {
				t.Errorf("knownPaths: expected %v, got %v", tt.wantKnownPaths, indexer.knownPaths)
			}

			if !reflect.DeepEqual(indexer.elementsPerTag, tt.wantElementsPerTag) {
				t.Errorf("elementsPerTag: expected %v, got %v", tt.wantElementsPerTag, indexer.elementsPerTag)
			}
		})
	}
}

func Test_collectionIndexer_isAccumulator(t *testing.T) {
	etmA := newTestElementTagMapping(t, "elem_a", "tv1", 1, false, nil)
	etmB := newTestElementTagMapping(t, "elem_b", "tv1", 1, false, nil)

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
	etmA := newTestElementTagMapping(t, "elem_a", "tv1", 1, false, nil)
	etmB := newTestElementTagMapping(t, "elem_b", "tv1", 1, false, nil)

	cA := candidate{etm: etmA}
	cB := candidate{etm: etmB}

	cA0 := candidate{etm: etmA, value: newTestSubmissionValue(t, "elem_a", "a0", nil), path: []int{0}}
	cA1 := candidate{etm: etmA, value: newTestSubmissionValue(t, "elem_a", "a1", nil), path: []int{1}}
	cA1Empty := candidate{etm: etmA, path: []int{1}}

	tests := []struct {
		name        string
		segments    []string
		sequence    []candidate
		expected    [][][]int
		expectedErr error
	}{
		{
			name:     "should assign sequential row per element and reuse it for the same element",
			segments: []string{"appFormData[*]"},
			sequence: []candidate{cA, cA, cB, cA, cB},
			expected: [][][]int{{{0}}, {{0}}, {{1}}, {{0}}, {{1}}},
		},
		{
			name:     "should assign a separate row per instance of the same element",
			segments: []string{"appFormData[*]"},
			sequence: []candidate{cA0, cA1, cB, cA0},
			expected: [][][]int{{{0}}, {{1}}, {{2}}, {{0}}},
		},
		{
			name:     "should keep empty answers on their instance row",
			segments: []string{"appFormData[*]"},
			sequence: []candidate{cA0, cA1Empty, cA1},
			expected: [][][]int{{{0}}, {{1}}, {{1}}},
		},
		{
			name:        "should yield ErrUnsupportedCollectionDepth for nested segments",
			segments:    []string{"acc[*]", "acc[*].lines[*]"},
			sequence:    []candidate{cA},
			expectedErr: ErrUnsupportedCollectionDepth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexer := newCollectionIndexer()
			for i, c := range tt.sequence {
				got, err := indexer.accumulatorCollection(tt.segments, c, "test.keypath")
				if tt.expectedErr != nil {
					if !errors.Is(err, tt.expectedErr) {
						t.Fatalf("step %d: expected error %v, got %v", i, tt.expectedErr, err)
					}
					continue
				}

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
	etmDynamic := newTestElementTagMapping(t, "elem1", "tv1", 1, false, nil)
	cDynamic0 := candidate{etm: etmDynamic, value: newTestSubmissionValue(t, "elem1", "v0", nil), path: []int{0}}
	cDynamic1 := candidate{etm: etmDynamic, value: newTestSubmissionValue(t, "elem1", "v1", nil), path: []int{1}}
	cDynamicEmpty := candidate{etm: etmDynamic, path: []int{1}}

	etmStatic := newTestElementTagMapping(t, "elem_static", "tv_static", 10, true, "USD")
	cStatic := candidate{etm: etmStatic}

	etmHetA := newTestElementTagMapping(t, "het_a", "tv_het", 1, false, nil)
	etmHetB := newTestElementTagMapping(t, "het_b", "tv_het", 1, false, nil)
	cHetA := candidate{etm: etmHetA, value: newTestSubmissionValue(t, "het_a", "A", nil)}
	cHetB := candidate{etm: etmHetB, value: newTestSubmissionValue(t, "het_b", "B", nil)}

	type evalStep struct {
		candidate   candidate
		expected    [][]int
		expectedErr error
	}

	tests := []struct {
		name       string
		segments   []string
		knownPaths [][]int
		steps      []evalStep
	}{
		{
			name:     "should return the native path for dynamic candidates",
			segments: []string{"items[*]"},
			steps: []evalStep{
				{candidate: cDynamic0, expected: [][]int{{0}}},
				{candidate: cDynamic1, expected: [][]int{{1}}},
			},
		},
		{
			name:     "should keep empty answers on their instance path",
			segments: []string{"items[*]"},
			steps: []evalStep{
				{candidate: cDynamicEmpty, expected: [][]int{{1}}},
			},
		},
		{
			name:       "should fan out static tag across all known paths",
			segments:   []string{"items[*]"},
			knownPaths: [][]int{{0}, {1}},
			steps: []evalStep{
				{candidate: cStatic, expected: [][]int{{0}, {1}}},
			},
		},
		{
			name:     "should return nil for static tag on empty entity collection",
			segments: []string{"items[*]"},
			steps: []evalStep{
				{candidate: cStatic, expected: nil},
			},
		},
		{
			name:     "should assign sequential row for heterogeneous candidates without a path",
			segments: []string{"items[*]"},
			steps: []evalStep{
				{candidate: cHetA, expected: [][]int{{0}}},
				{candidate: cHetB, expected: [][]int{{1}}},
			},
		},
		{
			name:     "should yield ErrCollectionDepthMismatch when path is shorter than tag depth",
			segments: []string{"orders[*]", "orders[*].lines[*]"},
			steps: []evalStep{
				{candidate: cDynamic0, expectedErr: ErrCollectionDepthMismatch},
			},
		},
		{
			name:     "should yield ErrMissingCollectionIndex for nested tag without a path",
			segments: []string{"orders[*]", "orders[*].lines[*]"},
			steps: []evalStep{
				{candidate: cHetA, expectedErr: ErrMissingCollectionIndex},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexer := newCollectionIndexer()
			for _, p := range tt.knownPaths {
				indexer.addKnownPath(tt.segments[len(tt.segments)-1], p)
			}

			for i, s := range tt.steps {
				got, err := indexer.entityCollection(tt.segments, s.candidate, "test.keypath")
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

func Test_collectionIndexer_addKnownPath(t *testing.T) {
	indexer := newCollectionIndexer()

	path := []int{1}
	indexer.addKnownPath("items[*]", path)
	indexer.addKnownPath("items[*]", []int{0})
	indexer.addKnownPath("items[*]", []int{1})

	path[0] = 99

	want := [][]int{{0}, {1}}
	if got := indexer.knownPaths["items[*]"]; !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func newTestTag(keyPath string) domain.Tag {
	return domain.Tag{KeyPath: keyPath}
}

func newTestElementTagMapping(t *testing.T, elementID domain.ElementID, tagVersionID domain.TagVersionID, priority int, hasStatic bool, staticValue any) *domain.ElementTagMapping {
	t.Helper()

	etm, err := domain.NewElementTagMapping(elementID, tagVersionID, priority, hasStatic, staticValue)
	if err != nil {
		t.Fatalf("NewElementTagMapping(%q, %q): unexpected error: %v", elementID, tagVersionID, err)
	}

	return etm
}

func newTestSubmissionValue(t *testing.T, elementID domain.ElementID, value any, itemIndex *int) *domain.SubmissionValue {
	t.Helper()

	sv, err := domain.NewSubmissionValue(elementID, value, itemIndex)
	if err != nil {
		t.Fatalf("NewSubmissionValue(%q): unexpected error: %v", elementID, err)
	}

	return sv
}
