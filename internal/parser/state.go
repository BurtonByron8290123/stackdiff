package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

// TerraformState represents the top-level structure of a .tfstate file.
type TerraformState struct {
	Version          int        `json:"version"`
	TerraformVersion string     `json:"terraform_version"`
	Resources        []Resource `json:"resources"`
}

// Resource represents a single resource block in the state file.
type Resource struct {
	Module    string     `json:"module,omitempty"`
	Mode      string     `json:"mode"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Instances []Instance `json:"instances"`
}

// Instance holds the actual attribute values for a resource instance.
type Instance struct {
	IndexKey   interface{}            `json:"index_key,omitempty"`
	Attributes map[string]interface{} `json:"attributes"`
}

// ResourceKey returns a unique string identifier for the resource.
func (r Resource) ResourceKey() string {
	if r.Module != "" {
		return fmt.Sprintf("%s.%s.%s", r.Module, r.Type, r.Name)
	}
	return fmt.Sprintf("%s.%s", r.Type, r.Name)
}

// ParseStateFile reads and parses a Terraform state file from disk.
func ParseStateFile(path string) (*TerraformState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", path, err)
	}

	var state TerraformState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file %q: %w", path, err)
	}

	return &state, nil
}
