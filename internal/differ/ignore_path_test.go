package differ

import (
	"testing"

	"github.com/you/stackdiff/internal/differ/model"
)

func makeAttrDiffs(keys ...string) []model.AttributeDiff {
	out := make([]model.AttributeDiff, len(keys))
	for i, k := range keys {
		out[i] = model.AttributeDiff{Key: k, OldValue: "a", NewValue: "b"}
	}
	return out
}

func TestNewPathIgnoreList_EmptyPatterns(t *testing.T) {
	pl := NewPathIgnoreList([]string{})
	if len(pl.patterns) != 0 {
		t.Fatalf("expected 0 patterns, got %d", len(pl.patterns))
	}
}

func TestNewPathIgnoreList_TrimsWhitespace(t *testing.T) {
	pl := NewPathIgnoreList([]string{"  tags.*  ", "", "  "})
	if len(pl.patterns) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(pl.patterns))
	}
	if pl.patterns[0] != "tags.*" {
		t.Errorf("unexpected pattern: %q", pl.patterns[0])
	}
}

func TestFilterAttributeDiffs_NoPatterns(t *testing.T) {
	pl := NewPathIgnoreList(nil)
	diffs := makeAttrDiffs("name", "tags.env", "size")
	result := pl.FilterAttributeDiffs(diffs)
	if len(result) != 3 {
		t.Errorf("expected 3 diffs unchanged, got %d", len(result))
	}
}

func TestFilterAttributeDiffs_ExactMatch(t *testing.T) {
	pl := NewPathIgnoreList([]string{"metadata.name"})
	diffs := makeAttrDiffs("metadata.name", "metadata.namespace", "size")
	result := pl.FilterAttributeDiffs(diffs)
	if len(result) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(result))
	}
	for _, d := range result {
		if d.Key == "metadata.name" {
			t.Error("metadata.name should have been filtered out")
		}
	}
}

func TestFilterAttributeDiffs_WildcardSuffix(t *testing.T) {
	pl := NewPathIgnoreList([]string{"tags.*"})
	diffs := makeAttrDiffs("tags.env", "tags.owner", "name", "tags.env.sub")
	result := pl.FilterAttributeDiffs(diffs)
	// "tags.env" and "tags.owner" removed; "tags.env.sub" has nested dot so kept
	if len(result) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(result))
	}
	for _, d := range result {
		if d.Key == "tags.env" || d.Key == "tags.owner" {
			t.Errorf("key %q should have been filtered", d.Key)
		}
	}
}

func TestFilterAttributeDiffs_StarWildcard(t *testing.T) {
	pl := NewPathIgnoreList([]string{"*"})
	diffs := makeAttrDiffs("name", "size", "nested.key")
	result := pl.FilterAttributeDiffs(diffs)
	// only top-level single-segment keys removed
	if len(result) != 1 {
		t.Fatalf("expected 1 diff (nested.key), got %d", len(result))
	}
	if result[0].Key != "nested.key" {
		t.Errorf("expected nested.key, got %q", result[0].Key)
	}
}

func TestMatchPath_Exact(t *testing.T) {
	if !matchPath("foo.bar", "foo.bar") {
		t.Error("exact match should return true")
	}
	if matchPath("foo.bar", "foo.baz") {
		t.Error("non-matching exact should return false")
	}
}
