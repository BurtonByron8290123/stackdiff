// Package baseline provides functionality to save and load a drift report
// as a baseline, enabling comparison against future runs.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/stackdiff/internal/differ"
)

// Snapshot represents a saved baseline of drift items.
type Snapshot struct {
	CreatedAt time.Time          `json:"created_at"`
	Items     []differ.DriftItem `json:"items"`
}

// Save writes the given drift items to a JSON baseline file at path.
func Save(path string, items []differ.DriftItem) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Items:     items,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", path, err)
	}
	return nil
}

// Load reads a baseline snapshot from the JSON file at path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("baseline: file not found: %s", path)
		}
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &snap, nil
}

// Subtract returns the drift items in current that are NOT present in the
// baseline snapshot (keyed by resource address + change type).
func Subtract(current []differ.DriftItem, base *Snapshot) []differ.DriftItem {
	type key struct {
		Address    string
		ChangeType string
	}
	seen := make(map[key]struct{}, len(base.Items))
	for _, item := range base.Items {
		seen[key{item.Address, item.ChangeType}] = struct{}{}
	}
	var result []differ.DriftItem
	for _, item := range current {
		if _, ok := seen[key{item.Address, item.ChangeType}]; !ok {
			result = append(result, item)
		}
	}
	return result
}
