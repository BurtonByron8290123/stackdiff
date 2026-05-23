package differ

import (
	"testing"
)

func TestSuppressResource_MatchesPrefix(t *testing.T) {
	il := NewIgnoreList([]IgnoreRule{
		{AddressPrefix: "module.legacy"},
	})

	if !il.SuppressResource("module.legacy.aws_s3_bucket.data") {
		t.Error("expected resource to be suppressed")
	}
}

func TestSuppressResource_NoMatch(t *testing.T) {
	il := NewIgnoreList([]IgnoreRule{
		{AddressPrefix: "module.legacy"},
	})

	if il.SuppressResource("aws_instance.web") {
		t.Error("expected resource NOT to be suppressed")
	}
}

func TestSuppressResource_EmptyList(t *testing.T) {
	il := NewIgnoreList(nil)
	if il.SuppressResource("module.anything") {
		t.Error("empty list should never suppress")
	}
}

func TestFilterAttributes_RemovesMatchedKey(t *testing.T) {
	il := NewIgnoreList([]IgnoreRule{
		{AttributeKey: "tags.LastModified"},
	})

	attrs := map[string]string{
		"id":               "abc",
		"tags.LastModified": "2024-01-01",
		"name":             "foo",
	}

	out := il.FilterAttributes(attrs)

	if _, ok := out["tags.LastModified"]; ok {
		t.Error("expected 'tags.LastModified' to be removed")
	}
	if out["id"] != "abc" || out["name"] != "foo" {
		t.Error("expected non-ignored attributes to be preserved")
	}
}

func TestFilterAttributes_NoRulesPreservesAll(t *testing.T) {
	il := NewIgnoreList(nil)

	attrs := map[string]string{"a": "1", "b": "2"}
	out := il.FilterAttributes(attrs)

	if len(out) != 2 {
		t.Errorf("expected 2 attributes, got %d", len(out))
	}
}

func TestFilterAttributes_NilInput(t *testing.T) {
	il := NewIgnoreList([]IgnoreRule{{AttributeKey: "x"}})
	out := il.FilterAttributes(nil)
	if out != nil {
		t.Error("expected nil output for nil input")
	}
}

func TestFilterAttributes_MultipleRules(t *testing.T) {
	il := NewIgnoreList([]IgnoreRule{
		{AttributeKey: "last_updated"},
		{AttributeKey: "checksum"},
	})

	attrs := map[string]string{
		"last_updated": "ts",
		"checksum":     "abc123",
		"region":       "us-east-1",
	}

	out := il.FilterAttributes(attrs)
	if len(out) != 1 || out["region"] != "us-east-1" {
		t.Errorf("unexpected output: %v", out)
	}
}
