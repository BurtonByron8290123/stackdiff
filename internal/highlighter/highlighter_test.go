package highlighter_test

import (
	"strings"
	"testing"

	"github.com/your-org/stackdiff/internal/highlighter"
)

func TestLabel_Disabled(t *testing.T) {
	h := highlighter.New(false)
	for _, ct := range []string{"added", "removed", "changed", "unknown"} {
		got := h.Label(ct)
		if got != ct {
			t.Errorf("disabled: Label(%q) = %q, want %q", ct, got, ct)
		}
	}
}

func TestLabel_Enabled_ContainsEscapeCode(t *testing.T) {
	h := highlighter.New(true)
	cases := []struct {
		changeType string
		wantCode   string
	}{
		{"added", "\033[32m"},
		{"removed", "\033[31m"},
		{"changed", "\033[33m"},
	}
	for _, tc := range cases {
		got := h.Label(tc.changeType)
		if !strings.Contains(got, tc.wantCode) {
			t.Errorf("Label(%q) missing escape code %q; got %q", tc.changeType, tc.wantCode, got)
		}
		if !strings.Contains(got, "\033[0m") {
			t.Errorf("Label(%q) missing reset code; got %q", tc.changeType, got)
		}
		if !strings.Contains(got, tc.changeType) {
			t.Errorf("Label(%q) missing original text; got %q", tc.changeType, got)
		}
	}
}

func TestLabel_Enabled_UnknownType(t *testing.T) {
	h := highlighter.New(true)
	got := h.Label("noop")
	if got != "noop" {
		t.Errorf("Label(\"noop\") = %q, want %q", got, "noop")
	}
}

func TestValue_Disabled(t *testing.T) {
	h := highlighter.New(false)
	got := h.Value("added", "somevalue")
	if got != "somevalue" {
		t.Errorf("disabled: Value() = %q, want %q", got, "somevalue")
	}
}

func TestValue_Enabled(t *testing.T) {
	h := highlighter.New(true)
	got := h.Value("removed", "old-val")
	if !strings.Contains(got, "\033[31m") {
		t.Errorf("Value(removed) missing red code; got %q", got)
	}
	if !strings.Contains(got, "old-val") {
		t.Errorf("Value(removed) missing original value; got %q", got)
	}
}

func TestValue_Enabled_UnknownType(t *testing.T) {
	h := highlighter.New(true)
	got := h.Value("noop", "x")
	if got != "x" {
		t.Errorf("Value(noop) = %q, want plain %q", got, "x")
	}
}
