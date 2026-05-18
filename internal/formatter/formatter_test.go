package formatter

import (
	"strings"
	"testing"

	"github.com/stackdiff/internal/differ"
)

func makeItems() []differ.DriftItem {
	return []differ.DriftItem{
		{
			Address:    "aws_instance.web",
			ChangeType: "modified",
			ChangedAttributes: []differ.AttributeDiff{
				{Key: "instance_type", OldValue: "t2.micro", NewValue: "t3.small"},
			},
		},
		{
			Address:    "aws_s3_bucket.logs",
			ChangeType: "added",
		},
	}
}

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct{ input string; want Format }{
		{"text", FormatText},
		{"", FormatText},
		{"json", FormatJSON},
		{"markdown", FormatMarkdown},
		{"md", FormatMarkdown},
		{"JSON", FormatJSON},
	}
	for _, c := range cases {
		got, err := ParseFormat(c.input)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", c.input, err)
		}
		if got != c.want {
			t.Errorf("ParseFormat(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("xml")
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestRender_Text(t *testing.T) {
	out, err := Render(makeItems(), FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "aws_instance.web") {
		t.Error("text output missing resource address")
	}
	if !strings.Contains(out, "instance_type") {
		t.Error("text output missing changed attribute")
	}
	if !strings.Contains(out, "2 change(s)") {
		t.Error("text output missing summary count")
	}
}

func TestRender_Markdown(t *testing.T) {
	out, err := Render(makeItems(), FormatMarkdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "## Drift Report") {
		t.Error("markdown output missing heading")
	}
	if !strings.Contains(out, "`aws_instance.web`") {
		t.Error("markdown output missing resource address")
	}
}

func TestRender_JSON(t *testing.T) {
	out, err := Render(makeItems(), FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"total": 2`) {
		t.Error("json output missing total field")
	}
	if !strings.Contains(out, `"aws_instance.web"`) {
		t.Error("json output missing resource address")
	}
}

func TestRender_NoDrift_Text(t *testing.T) {
	out, err := Render(nil, FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", out)
	}
}
