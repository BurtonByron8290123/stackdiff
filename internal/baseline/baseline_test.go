package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/stackdiff/internal/baseline"
	"github.com/user/stackdiff/internal/differ"
)

func makeItems() []differ.DriftItem {
	return []differ.DriftItem{
		{Address: "aws_instance.web", ChangeType: "changed"},
		{Address: "aws_s3_bucket.data", ChangeType: "added"},
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	items := makeItems()
	if err := baseline.Save(path, items); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(snap.Items) != len(items) {
		t.Errorf("expected %d items, got %d", len(items), len(snap.Items))
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if snap.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("CreatedAt is in the future")
	}
}

func TestLoad_Missing(t *testing.T) {
	_, err := baseline.Load("/nonexistent/path/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)

	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("expected unmarshal error")
	}
}

func TestSubtract_NewItems(t *testing.T) {
	base := &baseline.Snapshot{
		Items: []differ.DriftItem{
			{Address: "aws_instance.web", ChangeType: "changed"},
		},
	}
	current := []differ.DriftItem{
		{Address: "aws_instance.web", ChangeType: "changed"},
		{Address: "aws_s3_bucket.data", ChangeType: "added"},
	}
	result := baseline.Subtract(current, base)
	if len(result) != 1 {
		t.Fatalf("expected 1 new item, got %d", len(result))
	}
	if result[0].Address != "aws_s3_bucket.data" {
		t.Errorf("unexpected address: %s", result[0].Address)
	}
}

func TestSubtract_NoDiff(t *testing.T) {
	items := makeItems()
	base := &baseline.Snapshot{Items: items}
	result := baseline.Subtract(items, base)
	if len(result) != 0 {
		t.Errorf("expected 0 new items, got %d", len(result))
	}
}
