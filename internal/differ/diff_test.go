package differ_test

import (
	"testing"

	"github.com/stackdiff/internal/differ"
	"github.com/stackdiff/internal/parser"
)

func makeState(resources []parser.Resource) *parser.TerraformState {
	return &parser.TerraformState{Resources: resources}
}

func res(rtype, name string, values map[string]interface{}) parser.Resource {
	return parser.Resource{Type: rtype, Name: name, Values: values}
}

func TestCompare_NoDrift(t *testing.T) {
	v := map[string]interface{}{"size": "t2.micro"}
	base := makeState([]parser.Resource{res("aws_instance", "web", v)})
	cur := makeState([]parser.Resource{res("aws_instance", "web", v)})

	report, err := differ.Compare(base, cur)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.HasDrift() {
		t.Errorf("expected no drift, got %d diff(s)", len(report.Diffs))
	}
}

func TestCompare_Added(t *testing.T) {
	base := makeState([]parser.Resource{})
	cur := makeState([]parser.Resource{res("aws_s3_bucket", "logs", map[string]interface{}{"region": "us-east-1"})})

	report, err := differ.Compare(base, cur)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Diffs) != 1 || report.Diffs[0].ChangeType != differ.Added {
		t.Errorf("expected 1 added diff, got %+v", report.Diffs)
	}
}

func TestCompare_Removed(t *testing.T) {
	base := makeState([]parser.Resource{res("aws_instance", "web", map[string]interface{}{"size": "t2.micro"})})
	cur := makeState([]parser.Resource{})

	report, err := differ.Compare(base, cur)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Diffs) != 1 || report.Diffs[0].ChangeType != differ.Removed {
		t.Errorf("expected 1 removed diff, got %+v", report.Diffs)
	}
}

func TestCompare_Modified(t *testing.T) {
	base := makeState([]parser.Resource{res("aws_instance", "web", map[string]interface{}{"size": "t2.micro"})})
	cur := makeState([]parser.Resource{res("aws_instance", "web", map[string]interface{}{"size": "t2.large"})})

	report, err := differ.Compare(base, cur)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Diffs) != 1 || report.Diffs[0].ChangeType != differ.Modified {
		t.Errorf("expected 1 modified diff, got %+v", report.Diffs)
	}
}

func TestCompare_NilState(t *testing.T) {
	_, err := differ.Compare(nil, nil)
	if err == nil {
		t.Error("expected error for nil states")
	}
}
