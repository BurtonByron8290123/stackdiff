package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stackdiff/internal/differ"
)

func TestWrite_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, differ.DiffResult{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWrite_Added(t *testing.T) {
	var buf bytes.Buffer
	result := differ.DiffResult{
		Added: []differ.ResourceSummary{
			{Address: "aws_instance.web", Type: "aws_instance"},
		},
	}
	if err := Write(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ Added (1)") {
		t.Errorf("expected added section, got: %s", out)
	}
	if !strings.Contains(out, "aws_instance.web") {
		t.Errorf("expected resource address in output, got: %s", out)
	}
}

func TestWrite_Removed(t *testing.T) {
	var buf bytes.Buffer
	result := differ.DiffResult{
		Removed: []differ.ResourceSummary{
			{Address: "aws_s3_bucket.data", Type: "aws_s3_bucket"},
		},
	}
	if err := Write(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "- Removed (1)") {
		t.Errorf("expected removed section, got: %s", out)
	}
}

func TestWrite_Summary(t *testing.T) {
	var buf bytes.Buffer
	result := differ.DiffResult{
		Added:   []differ.ResourceSummary{{Address: "a", Type: "t"}},
		Removed: []differ.ResourceSummary{{Address: "b", Type: "t"}},
	}
	if err := Write(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Summary: +1 added, -1 removed, ~0 changed") {
		t.Errorf("expected summary line, got: %s", out)
	}
}
