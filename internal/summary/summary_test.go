package summary_test

import (
	"testing"

	"github.com/stackdiff/internal/differ"
	"github.com/stackdiff/internal/summary"
)

func items(types ...string) []differ.DriftItem {
	out := make([]differ.DriftItem, 0, len(types))
	for i, t := range types {
		out = append(out, differ.DriftItem{
			Address:    fmt.Sprintf("resource.r%d", i),
			ChangeType: t,
		})
	}
	return out
}

func TestCompute_Empty(t *testing.T) {
	s := summary.Compute(nil)
	if s.Total != 0 || s.HasDrift() {
		t.Errorf("expected zero stats, got %+v", s)
	}
}

func TestCompute_Mixed(t *testing.T) {
	list := []differ.DriftItem{
		{Address: "a", ChangeType: "added"},
		{Address: "b", ChangeType: "added"},
		{Address: "c", ChangeType: "removed"},
		{Address: "d", ChangeType: "changed"},
		{Address: "e", ChangeType: "changed"},
		{Address: "f", ChangeType: "changed"},
	}
	s := summary.Compute(list)
	if s.Added != 2 {
		t.Errorf("Added: want 2, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("Removed: want 1, got %d", s.Removed)
	}
	if s.Changed != 3 {
		t.Errorf("Changed: want 3, got %d", s.Changed)
	}
	if s.Total != 6 {
		t.Errorf("Total: want 6, got %d", s.Total)
	}
	if !s.HasDrift() {
		t.Error("expected HasDrift() == true")
	}
}

func TestFormat_NoDrift(t *testing.T) {
	s := summary.Compute(nil)
	got := s.Format()
	want := "No drift detected."
	if got != want {
		t.Errorf("Format(): want %q, got %q", want, got)
	}
}

func TestFormat_WithDrift(t *testing.T) {
	list := []differ.DriftItem{
		{Address: "a", ChangeType: "added"},
		{Address: "b", ChangeType: "removed"},
	}
	s := summary.Compute(list)
	got := s.Format()
	want := "Drift summary: 2 change(s) total — 1 added, 1 removed, 0 changed."
	if got != want {
		t.Errorf("Format(): want %q, got %q", want, got)
	}
}
