package differ

import (
	"testing"

	"github.com/user/stackdiff/internal/differ/model"
)

func makeAttrDiff(before, after string) model.AttributeDiff {
	return model.AttributeDiff{Before: before, After: after}
}

func TestSuppressResource_MatchesPrefix(t *testing.T) {
	il := NewIgnoreList([]string{"module.ignored"}, nil)
	if !il.SuppressResource("module.ignored.aws_instance.foo") {
		t.Fatal("expected resource to be suppressed")
	}
}

func TestSuppressResource_NoMatch(t *testing.T) {
	il := NewIgnoreList([]string{"module.ignored"}, nil)
	if il.SuppressResource("aws_s3_bucket.data") {
		t.Fatal("expected resource NOT to be suppressed")
	}
}

func TestSuppressResource_EmptyList(t *testing.T) {
	il := NewIgnoreList(nil, nil)
	if il.SuppressResource("anything") {
		t.Fatal("empty list should suppress nothing")
	}
}

func TestFilterAttributes_RemovesMatchedKey(t *testing.T) {
	il := NewIgnoreList(nil, []string{"tags"})
	attrs := map[string]model.AttributeDiff{
		"tags":       makeAttrDiff("a", "b"),
		"tags.Name":  makeAttrDiff("x", "y"),
		"bucket_name": makeAttrDiff("old", "new"),
	}
	out := il.FilterAttributes(attrs)
	if _, ok := out["tags"]; ok {
		t.Error("'tags' should be filtered")
	}
	if _, ok := out["tags.Name"]; ok {
		t.Error("'tags.Name' should be filtered by prefix rule")
	}
	if _, ok := out["bucket_name"]; !ok {
		t.Error("'bucket_name' should be preserved")
	}
}

func TestFilterAttributes_NoRulesPreservesAll(t *testing.T) {
	il := NewIgnoreList(nil, nil)
	attrs := map[string]model.AttributeDiff{"id": makeAttrDiff("1", "2")}
	out := il.FilterAttributes(attrs)
	if len(out) != len(attrs) {
		t.Errorf("expected %d attrs, got %d", len(attrs), len(out))
	}
}

func TestApply_FiltersResourcesAndAttributes(t *testing.T) {
	il := NewIgnoreList([]string{"module.skip"}, []string{"tags"})
	items := []model.DriftItem{
		{Address: "module.skip.aws_instance.foo", ChangeType: model.ChangeModified},
		{
			Address:    "aws_s3_bucket.keep",
			ChangeType: model.ChangeModified,
			Attributes: map[string]model.AttributeDiff{
				"tags":   makeAttrDiff("a", "b"),
				"region": makeAttrDiff("us-east-1", "eu-west-1"),
			},
		},
	}
	out := il.Apply(items)
	if len(out) != 1 {
		t.Fatalf("expected 1 item, got %d", len(out))
	}
	if _, ok := out[0].Attributes["tags"]; ok {
		t.Error("'tags' attribute should have been removed")
	}
	if _, ok := out[0].Attributes["region"]; !ok {
		t.Error("'region' attribute should be preserved")
	}
}
