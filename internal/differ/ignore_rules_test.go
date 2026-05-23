package differ

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTempRulesFile(t *testing.T, rules IgnoreRulesFile) string {
	t.Helper()
	data, err := json.Marshal(rules)
	if err != nil {
		t.Fatalf("marshal rules: %v", err)
	}
	p := filepath.Join(t.TempDir(), "ignore_rules.json")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatalf("write rules file: %v", err)
	}
	return p
}

func TestLoadIgnoreRulesFile_EmptyPath(t *testing.T) {
	rules, err := LoadIgnoreRulesFile("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules.ResourcePrefixes) != 0 || len(rules.AttributeKeys) != 0 {
		t.Errorf("expected empty rules, got %+v", rules)
	}
}

func TestLoadIgnoreRulesFile_Valid(t *testing.T) {
	expected := IgnoreRulesFile{
		ResourcePrefixes: []string{"module.temp", "aws_instance.legacy"},
		AttributeKeys:    []string{"tags", "last_modified"},
	}
	path := writeTempRulesFile(t, expected)

	rules, err := LoadIgnoreRulesFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules.ResourcePrefixes) != 2 {
		t.Errorf("expected 2 resource prefixes, got %d", len(rules.ResourcePrefixes))
	}
	if len(rules.AttributeKeys) != 2 {
		t.Errorf("expected 2 attribute keys, got %d", len(rules.AttributeKeys))
	}
}

func TestLoadIgnoreRulesFile_Missing(t *testing.T) {
	_, err := LoadIgnoreRulesFile("/nonexistent/path/rules.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadIgnoreRulesFile_InvalidJSON(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(p, []byte("not-json{"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, err := LoadIgnoreRulesFile(p)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestIgnoreRulesFile_ToIgnoreList(t *testing.T) {
	rules := IgnoreRulesFile{
		ResourcePrefixes: []string{"module.skip"},
		AttributeKeys:    []string{"created_at"},
	}
	il := rules.ToIgnoreList()
	if il == nil {
		t.Fatal("expected non-nil IgnoreList")
	}
}
