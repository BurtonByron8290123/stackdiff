// Package config handles loading and validation of stackdiff CLI configuration.
package config

import (
	"errors"
	"strings"
)

// Config holds all runtime options for a stackdiff comparison run.
type Config struct {
	BaseFile    string
	TargetFile  string
	Format      string
	OutputFile  string
	Overwrite   bool
	FilterTypes []string
	FilterName  string
	NoColor     bool
}

// Validate checks that required fields are present and values are acceptable.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.BaseFile) == "" {
		return errors.New("base state file path is required")
	}
	if strings.TrimSpace(c.TargetFile) == "" {
		return errors.New("target state file path is required")
	}
	validFormats := map[string]bool{
		"text":     true,
		"markdown": true,
		"json":     true,
	}
	if c.Format != "" && !validFormats[strings.ToLower(c.Format)] {
		return errors.New("invalid format: must be one of text, markdown, json")
	}
	return nil
}

// Default returns a Config with sensible defaults applied.
func Default() *Config {
	return &Config{
		Format:    "text",
		Overwrite: false,
		NoColor:   false,
	}
}
