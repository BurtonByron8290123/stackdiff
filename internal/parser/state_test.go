package parser_test

import (
	"os"
	"testing"

	"github.com/stackdiff/stackdiff/internal/parser"
)

const sampleState = `{
  "version": 4,
  "terraform_version": "1.5.0",
  "resources": [
    {
      "mode": "managed",
      "type": "aws_instance",
      "name": "web",
      "instances": [
        {
          "attributes": {
            "id": "i-0abc123",
            "instance_type": "t3.micro",
            "ami": "ami-0deadbeef"
          }
        }
      ]
    }
  ]
}`

func writeTempState(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.tfstate")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParseStateFile_Valid(t *testing.T) {
	path := writeTempState(t, sampleState)
	state, err := parser.ParseStateFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Version != 4 {
		t.Errorf("expected version 4, got %d", state.Version)
	}
	if len(state.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(state.Resources))
	}
	res := state.Resources[0]
	if res.ResourceKey() != "aws_instance.web" {
		t.Errorf("unexpected resource key: %s", res.ResourceKey())
	}
}

func TestParseStateFile_Missing(t *testing.T) {
	_, err := parser.ParseStateFile("/nonexistent/path.tfstate")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseStateFile_InvalidJSON(t *testing.T) {
	path := writeTempState(t, `{not valid json`)
	_, err := parser.ParseStateFile(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
