package exporter_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/stackdiff/internal/exporter"
)

func TestExport_Stdout(t *testing.T) {
	var buf bytes.Buffer
	content := []byte("drift report content")

	err := exporter.Export(exporter.Destination{}, content, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != string(content) {
		t.Errorf("expected %q, got %q", content, buf.String())
	}
}

func TestExport_ToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.txt")
	content := []byte("file report")

	err := exporter.Export(exporter.Destination{Path: path, Overwrite: false}, content, os.Stdout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("expected %q, got %q", content, got)
	}
}

func TestExport_OverwriteFalse_Conflicts(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.txt")
	_ = os.WriteFile(path, []byte("existing"), 0o644)

	err := exporter.Export(exporter.Destination{Path: path, Overwrite: false}, []byte("new"), os.Stdout)
	if err == nil {
		t.Fatal("expected error when file exists and Overwrite is false")
	}
}

func TestExport_OverwriteTrue_Replaces(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.txt")
	_ = os.WriteFile(path, []byte("old content"), 0o644)

	newContent := []byte("new content")
	err := exporter.Export(exporter.Destination{Path: path, Overwrite: true}, newContent, os.Stdout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := os.ReadFile(path)
	if string(got) != string(newContent) {
		t.Errorf("expected %q, got %q", newContent, got)
	}
}

func TestExport_CreatesNestedDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "report.md")

	err := exporter.Export(exporter.Destination{Path: path, Overwrite: false}, []byte("md"), os.Stdout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist at %s", path)
	}
}
