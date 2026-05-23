package differ

import (
	"encoding/json"
	"fmt"
	"os"
)

// IgnoreRulesFile represents a JSON file containing ignore rules for the differ.
type IgnoreRulesFile struct {
	// ResourcePrefixes lists resource address prefixes to suppress entirely.
	ResourcePrefixes []string `json:"resource_prefixes"`
	// AttributeKeys lists attribute keys to suppress from drift output.
	AttributeKeys []string `json:"attribute_keys"`
}

// LoadIgnoreRulesFile reads and parses an ignore rules JSON file from the given path.
// Returns an empty IgnoreRulesFile if the path is empty.
func LoadIgnoreRulesFile(path string) (IgnoreRulesFile, error) {
	if path == "" {
		return IgnoreRulesFile{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return IgnoreRulesFile{}, fmt.Errorf("reading ignore rules file %q: %w", path, err)
	}

	var rules IgnoreRulesFile
	if err := json.Unmarshal(data, &rules); err != nil {
		return IgnoreRulesFile{}, fmt.Errorf("parsing ignore rules file %q: %w", path, err)
	}

	return rules, nil
}

// ToIgnoreList converts the IgnoreRulesFile into an IgnoreList ready for use
// during diffing.
func (r IgnoreRulesFile) ToIgnoreList() *IgnoreList {
	return NewIgnoreList(r.ResourcePrefixes, r.AttributeKeys)
}
