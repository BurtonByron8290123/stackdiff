package parser_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/example/stackdiff/internal/parser"
)

func writeStateFile(t *testing.T, resources []map[string]interface{}) string {
	t.Helper()
	state := map[string]interface{}{
		"version": 4,
		"resources": resources,
	}
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.CreateTemp(t.TempDir(), "state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.Write(data)
	_ = f.Close()
	return f.Name()
}

func TestCachedParser_ParseAndHit(t *testing.T) {
	path := writeStateFile(t, []map[string]interface{}{
		{"type": "aws_s3_bucket", "name": "my_bucket", "mode": "managed",
			"instances": []map[string]interface{}{
				{"attributes": map[string]interface{}{"id": "bucket-1"}},
			}},
	})

	cp := parser.NewCachedParser(t.TempDir())

	s1, err := cp.Parse(path)
	if err != nil {
		t.Fatalf("first parse: %v", err)
	}
	if len(s1.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(s1.Resources))
	}

	// Second call should return cached result.
	s2, err := cp.Parse(path)
	if err != nil {
		t.Fatalf("second parse: %v", err)
	}
	if s2.Resources[0].Name != s1.Resources[0].Name {
		t.Errorf("cache hit returned different resource name")
	}
}

func TestCachedParser_StaleAfterModification(t *testing.T) {
	path := writeStateFile(t, []map[string]interface{}{
		{"type": "aws_instance", "name": "web", "mode": "managed",
			"instances": []map[string]interface{}{
				{"attributes": map[string]interface{}{"id": "i-abc"}},
			}},
	})

	cp := parser.NewCachedParser(t.TempDir())
	if _, err := cp.Parse(path); err != nil {
		t.Fatal(err)
	}

	// Overwrite the file with a different resource.
	newState := map[string]interface{}{
		"version": 4,
		"resources": []map[string]interface{}{
			{"type": "aws_instance", "name": "web", "mode": "managed",
				"instances": []map[string]interface{}{
					{"attributes": map[string]interface{}{"id": "i-xyz"}},
				}},
		},
	}
	data, _ := json.Marshal(newState)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := cp.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if s.Resources[0].Instances[0].Attributes["id"] != "i-xyz" {
		t.Errorf("expected fresh parse after file change")
	}
}
