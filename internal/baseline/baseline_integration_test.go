package baseline_test

import (
	"path/filepath"
	"testing"

	"github.com/user/stackdiff/internal/baseline"
	"github.com/user/stackdiff/internal/differ"
)

// TestRoundTrip_EmptyItems ensures Save/Load works correctly for an empty
// slice without producing a null JSON array.
func TestRoundTrip_EmptyItems(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")

	if err := baseline.Save(path, []differ.DriftItem{}); err != nil {
		t.Fatalf("Save empty: %v", err)
	}
	snap, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load empty: %v", err)
	}
	if snap.Items == nil {
		// nil is acceptable; length must be 0
	}
	if len(snap.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(snap.Items))
	}
}

// TestSubtract_AllNew verifies that when the baseline is empty every current
// item is returned as new drift.
func TestSubtract_AllNew(t *testing.T) {
	base := &baseline.Snapshot{Items: []differ.DriftItem{}}
	current := makeItems()
	result := baseline.Subtract(current, base)
	if len(result) != len(current) {
		t.Errorf("expected %d items, got %d", len(current), len(result))
	}
}

// TestSubtract_SameAddressDifferentChangeType ensures that address alone is
// not enough to suppress an item — the change type must also match.
func TestSubtract_SameAddressDifferentChangeType(t *testing.T) {
	base := &baseline.Snapshot{
		Items: []differ.DriftItem{
			{Address: "aws_instance.web", ChangeType: "changed"},
		},
	}
	current := []differ.DriftItem{
		{Address: "aws_instance.web", ChangeType: "removed"},
	}
	result := baseline.Subtract(current, base)
	if len(result) != 1 {
		t.Fatalf("expected 1 item (different change type), got %d", len(result))
	}
}
