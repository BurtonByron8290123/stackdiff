package config

import (
	"flag"
	"testing"
)

func TestBindFlags_Defaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	c := BindFlags(fs)

	if err := fs.Parse([]string{}); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if c.Format != "text" {
		t.Errorf("expected default format 'text', got %q", c.Format)
	}
	if c.Overwrite {
		t.Error("expected Overwrite default false")
	}
	if c.NoColor {
		t.Error("expected NoColor default false")
	}
}

func TestBindFlags_AllArgs(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	c := BindFlags(fs)

	args := []string{
		"-base", "a.tfstate",
		"-target", "b.tfstate",
		"-format", "json",
		"-output", "report.json",
		"-overwrite",
		"-filter-name", "aws_",
		"-no-color",
	}
	if err := fs.Parse(args); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if c.BaseFile != "a.tfstate" {
		t.Errorf("expected base 'a.tfstate', got %q", c.BaseFile)
	}
	if c.TargetFile != "b.tfstate" {
		t.Errorf("expected target 'b.tfstate', got %q", c.TargetFile)
	}
	if c.Format != "json" {
		t.Errorf("expected format 'json', got %q", c.Format)
	}
	if c.OutputFile != "report.json" {
		t.Errorf("expected output 'report.json', got %q", c.OutputFile)
	}
	if !c.Overwrite {
		t.Error("expected Overwrite true")
	}
	if c.FilterName != "aws_" {
		t.Errorf("expected filter-name 'aws_', got %q", c.FilterName)
	}
	if !c.NoColor {
		t.Error("expected NoColor true")
	}
}

func TestBindFlags_ValidateAfterParse(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	c := BindFlags(fs)

	_ = fs.Parse([]string{"-base", "x.tfstate", "-target", "y.tfstate"})

	if err := c.Validate(); err != nil {
		t.Fatalf("expected valid config, got: %v", err)
	}
}
