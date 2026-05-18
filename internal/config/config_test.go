package config

import (
	"testing"
)

func TestValidate_Valid(t *testing.T) {
	c := &Config{
		BaseFile:   "base.tfstate",
		TargetFile: "target.tfstate",
		Format:     "json",
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_MissingBase(t *testing.T) {
	c := &Config{TargetFile: "target.tfstate"}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing base file")
	}
}

func TestValidate_MissingTarget(t *testing.T) {
	c := &Config{BaseFile: "base.tfstate"}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing target file")
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	c := &Config{
		BaseFile:   "base.tfstate",
		TargetFile: "target.tfstate",
		Format:     "yaml",
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestValidate_EmptyFormatAllowed(t *testing.T) {
	c := &Config{
		BaseFile:   "base.tfstate",
		TargetFile: "target.tfstate",
		Format:     "",
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected no error for empty format, got: %v", err)
	}
}

func TestDefault(t *testing.T) {
	c := Default()
	if c.Format != "text" {
		t.Errorf("expected default format 'text', got %q", c.Format)
	}
	if c.Overwrite {
		t.Error("expected Overwrite to default to false")
	}
	if c.NoColor {
		t.Error("expected NoColor to default to false")
	}
}
