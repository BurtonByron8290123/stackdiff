// Package sorter provides ordering utilities for drift report items.
package sorter

import (
	"fmt"
	"sort"

	"github.com/yourusername/stackdiff/internal/differ"
)

// SortOrder defines the available sort orders.
type SortOrder string

const (
	SortByName       SortOrder = "name"
	SortByType       SortOrder = "type"
	SortByChangeType SortOrder = "change"
)

// ParseSortOrder parses a string into a SortOrder, returning an error if invalid.
func ParseSortOrder(s string) (SortOrder, error) {
	switch SortOrder(s) {
	case SortByName, SortByType, SortByChangeType, "":
		return SortOrder(s), nil
	default:
		return "", fmt.Errorf("unknown sort order %q: must be one of name, type, change", s)
	}
}

// Apply sorts a copy of items according to the given order and returns it.
// If order is empty the original slice is returned unchanged.
func Apply(items []differ.DriftItem, order SortOrder) []differ.DriftItem {
	if len(items) == 0 || order == "" {
		return items
	}

	out := make([]differ.DriftItem, len(items))
	copy(out, items)

	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		switch order {
		case SortByType:
			if a.ResourceType != b.ResourceType {
				return a.ResourceType < b.ResourceType
			}
			return a.ResourceName < b.ResourceName
		case SortByChangeType:
			if a.ChangeType != b.ChangeType {
				return a.ChangeType < b.ChangeType
			}
			return a.ResourceName < b.ResourceName
		default: // SortByName
			return a.ResourceName < b.ResourceName
		}
	})

	return out
}
