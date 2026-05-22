package sorter_test

import (
	"testing"

	"github.com/yourusername/stackdiff/internal/differ"
	"github.com/yourusername/stackdiff/internal/sorter"
)

func items() []differ.DriftItem {
	return []differ.DriftItem{
		{ResourceName: "zebra", ResourceType: "aws_s3_bucket", ChangeType: "modified"},
		{ResourceName: "alpha", ResourceType: "aws_instance", ChangeType: "added"},
		{ResourceName: "mango", ResourceType: "aws_s3_bucket", ChangeType: "removed"},
		{ResourceName: "beta", ResourceType: "aws_instance", ChangeType: "added"},
	}
}

func TestParseSortOrder_Valid(t *testing.T) {
	for _, tc := range []string{"name", "type", "change", ""} {
		if _, err := sorter.ParseSortOrder(tc); err != nil {
			t.Errorf("expected %q to be valid, got error: %v", tc, err)
		}
	}
}

func TestParseSortOrder_Invalid(t *testing.T) {
	if _, err := sorter.ParseSortOrder("unknown"); err == nil {
		t.Error("expected error for unknown sort order")
	}
}

func TestApply_Empty(t *testing.T) {
	result := sorter.Apply(nil, sorter.SortByName)
	if result != nil {
		t.Errorf("expected nil for empty input, got %v", result)
	}
}

func TestApply_NoOrder(t *testing.T) {
	in := items()
	out := sorter.Apply(in, "")
	if &out[0] != &in[0] {
		t.Error("expected same slice when order is empty")
	}
}

func TestApply_SortByName(t *testing.T) {
	out := sorter.Apply(items(), sorter.SortByName)
	names := []string{out[0].ResourceName, out[1].ResourceName, out[2].ResourceName, out[3].ResourceName}
	expected := []string{"alpha", "beta", "mango", "zebra"}
	for i, want := range expected {
		if names[i] != want {
			t.Errorf("position %d: got %q, want %q", i, names[i], want)
		}
	}
}

func TestApply_SortByType(t *testing.T) {
	out := sorter.Apply(items(), sorter.SortByType)
	// aws_instance < aws_s3_bucket; within each type sorted by name
	if out[0].ResourceType != "aws_instance" || out[1].ResourceType != "aws_instance" {
		t.Errorf("expected first two items to be aws_instance, got %v, %v", out[0].ResourceType, out[1].ResourceType)
	}
	if out[0].ResourceName != "alpha" || out[1].ResourceName != "beta" {
		t.Errorf("expected alpha then beta within aws_instance, got %v, %v", out[0].ResourceName, out[1].ResourceName)
	}
}

func TestApply_SortByChangeType(t *testing.T) {
	out := sorter.Apply(items(), sorter.SortByChangeType)
	// added < modified < removed
	if out[0].ChangeType != "added" || out[1].ChangeType != "added" {
		t.Errorf("expected first two to be added, got %v, %v", out[0].ChangeType, out[1].ChangeType)
	}
	if out[2].ChangeType != "modified" {
		t.Errorf("expected third to be modified, got %v", out[2].ChangeType)
	}
	if out[3].ChangeType != "removed" {
		t.Errorf("expected last to be removed, got %v", out[3].ChangeType)
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	in := items()
	origFirst := in[0].ResourceName
	sorter.Apply(in, sorter.SortByName)
	if in[0].ResourceName != origFirst {
		t.Error("Apply mutated the input slice")
	}
}
