package filter

import (
	"testing"

	"github.com/user/stackdiff/internal/differ"
)

func driftItem(rtype, name, change string) differ.DriftItem {
	return differ.DriftItem{
		ResourceType: rtype,
		ResourceName: name,
		ChangeType:   change,
	}
}

func TestApply_NoFilter(t *testing.T) {
	items := []differ.DriftItem{
		driftItem("aws_s3_bucket", "my-bucket", "added"),
		driftItem("aws_instance", "web", "removed"),
	}
	result := Apply(items, Options{})
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
}

func TestApply_FilterByType(t *testing.T) {
	items := []differ.DriftItem{
		driftItem("aws_s3_bucket", "my-bucket", "added"),
		driftItem("aws_instance", "web", "removed"),
	}
	result := Apply(items, Options{ResourceType: "aws_s3_bucket"})
	if len(result) != 1 || result[0].ResourceType != "aws_s3_bucket" {
		t.Errorf("expected 1 aws_s3_bucket item, got %v", result)
	}
}

func TestApply_FilterByNamePrefix(t *testing.T) {
	items := []differ.DriftItem{
		driftItem("aws_instance", "web-prod", "changed"),
		driftItem("aws_instance", "db-prod", "changed"),
	}
	result := Apply(items, Options{NamePrefix: "web"})
	if len(result) != 1 || result[0].ResourceName != "web-prod" {
		t.Errorf("expected 1 item with prefix 'web', got %v", result)
	}
}

func TestApply_FilterByChangeType(t *testing.T) {
	items := []differ.DriftItem{
		driftItem("aws_s3_bucket", "bucket-a", "added"),
		driftItem("aws_instance", "web", "removed"),
		driftItem("aws_instance", "db", "changed"),
	}
	result := Apply(items, Options{ChangeTypes: []string{"added", "removed"}})
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
}

func TestApply_NoMatches(t *testing.T) {
	items := []differ.DriftItem{
		driftItem("aws_s3_bucket", "bucket-a", "added"),
	}
	result := Apply(items, Options{ResourceType: "aws_lambda_function"})
	if len(result) != 0 {
		t.Errorf("expected 0 items, got %d", len(result))
	}
}
