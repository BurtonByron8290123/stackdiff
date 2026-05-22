package truncator_test

import (
	"strings"
	"testing"

	"github.com/user/stackdiff/internal/differ"
	"github.com/user/stackdiff/internal/truncator"
)

func makeItem(addr, changeType string, attrs map[string][2]string) differ.DriftItem {
	return differ.DriftItem{
		Address:    addr,
		Type:       "aws_instance",
		ChangeType: changeType,
		Attributes: attrs,
	}
}

func TestApply_NoTruncationNeeded(t *testing.T) {
	items := []differ.DriftItem{
		makeItem("aws_instance.web", "changed", map[string][2]string{
			"ami": {"ami-abc", "ami-xyz"},
		}),
	}
	opts := truncator.DefaultOptions()
	got := truncator.Apply(items, opts)
	if got[0].Attributes["ami"][0] != "ami-abc" {
		t.Errorf("expected ami-abc, got %s", got[0].Attributes["ami"][0])
	}
}

func TestApply_TruncatesLongValue(t *testing.T) {
	long := strings.Repeat("x", 200)
	items := []differ.DriftItem{
		makeItem("aws_instance.web", "changed", map[string][2]string{
			"user_data": {long, "short"},
		}),
	}
	opts := truncator.Options{MaxLen: 50, Ellipsis: "..."}
	got := truncator.Apply(items, opts)
	v := got[0].Attributes["user_data"][0]
	if len([]rune(v)) >= 200 {
		t.Errorf("value was not truncated: len=%d", len(v))
	}
	if !strings.Contains(v, "...") {
		t.Errorf("expected ellipsis in truncated value, got: %s", v)
	}
}

func TestApply_ZeroMaxLen_Disables(t *testing.T) {
	long := strings.Repeat("y", 300)
	items := []differ.DriftItem{
		makeItem("aws_s3_bucket.data", "changed", map[string][2]string{
			"policy": {long, long},
		}),
	}
	opts := truncator.Options{MaxLen: 0}
	got := truncator.Apply(items, opts)
	if got[0].Attributes["policy"][0] != long {
		t.Error("expected value to be unchanged when MaxLen is 0")
	}
}

func TestApply_NilAttributes(t *testing.T) {
	items := []differ.DriftItem{
		makeItem("aws_instance.web", "added", nil),
	}
	opts := truncator.DefaultOptions()
	got := truncator.Apply(items, opts)
	if got[0].Attributes != nil {
		t.Error("expected nil attributes to remain nil")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	long := strings.Repeat("z", 200)
	original := []differ.DriftItem{
		makeItem("aws_instance.web", "changed", map[string][2]string{
			"tag": {long, "v2"},
		}),
	}
	opts := truncator.Options{MaxLen: 30, Ellipsis: "..."}
	truncator.Apply(original, opts)
	if original[0].Attributes["tag"][0] != long {
		t.Error("Apply must not mutate the original items")
	}
}
