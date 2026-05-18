package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

// TerraformState represents the top-level structure of a Terraform state file.
type TerraformState struct {
	Version   int        `json:"version"`
	TFVersion string     `json:"terraform_version"`
	Resources []Resource `json:"resources"`
}

// Resource represents a single managed resource in the state.
type Resource struct {
	Module string                 `json:"module,omitempty"`
	Mode   string                 `json:"mode"`
	Type   string                 `json:"type"`
	Name   string                 `json:"name"`
	Values map[string]interface{} `json:"values"`
}

// Address returns a canonical identifier for the resource.
func (r *Resource) Address() string {
	if r.Module != "" {
		return fmt.Sprintf("%s.%s.%s", r.Module, r.Type, r.Name)
	}
	return fmt.Sprintf("%s.%s", r.Type, r.Name)
}

// ParseStateFile reads and unmarshals a Terraform state file from disk.
func ParseStateFile(path string) (*TerraformState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading state file %q: %w", path, err)
	}

	var state TerraformState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parsing state file %q: %w", path, err)
	}

	return &state, nil
}
