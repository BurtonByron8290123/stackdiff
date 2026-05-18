package filter

import (
	"strings"

	"github.com/user/stackdiff/internal/differ"
)

// Options holds filtering criteria for drift results.
type Options struct {
	ResourceType string
	ChangeTypes  []string
	NamePrefix   string
}

// Apply filters a slice of DriftItems based on the provided Options.
// If an option field is empty/nil, that filter is not applied.
func Apply(items []differ.DriftItem, opts Options) []differ.DriftItem {
	var result []differ.DriftItem

	for _, item := range items {
		if opts.ResourceType != "" && !strings.EqualFold(item.ResourceType, opts.ResourceType) {
			continue
		}

		if opts.NamePrefix != "" && !strings.HasPrefix(item.ResourceName, opts.NamePrefix) {
			continue
		}

		if len(opts.ChangeTypes) > 0 && !containsChangeType(opts.ChangeTypes, item.ChangeType) {
			continue
		}

		result = append(result, item)
	}

	return result
}

// containsChangeType checks whether a given change type exists in the list.
func containsChangeType(types []string, ct string) bool {
	for _, t := range types {
		if strings.EqualFold(t, ct) {
			return true
		}
	}
	return false
}
