// Package truncator provides utilities for truncating long attribute values
// in drift reports to keep output readable.
package truncator

import (
	"fmt"
	"strings"

	"github.com/user/stackdiff/internal/differ"
)

const defaultMaxLen = 120

// Options controls truncation behaviour.
type Options struct {
	// MaxLen is the maximum character length for any single attribute value.
	// Values longer than this are truncated and suffixed with an ellipsis.
	// A value of 0 disables truncation.
	MaxLen int
	// Ellipsis is appended to truncated values. Defaults to "...".
	Ellipsis string
}

// DefaultOptions returns sensible truncation defaults.
func DefaultOptions() Options {
	return Options{
		MaxLen:   defaultMaxLen,
		Ellipsis: "...",
	}
}

// Apply truncates attribute values inside each DriftItem according to opts.
// It returns a new slice; the originals are not mutated.
func Apply(items []differ.DriftItem, opts Options) []differ.DriftItem {
	if opts.MaxLen <= 0 {
		return items
	}
	ellipsis := opts.Ellipsis
	if ellipsis == "" {
		ellipsis = "..."
	}

	out := make([]differ.DriftItem, len(items))
	for i, item := range items {
		copy := item
		copy.Attributes = truncateMap(item.Attributes, opts.MaxLen, ellipsis)
		out[i] = copy
	}
	return out
}

func truncateMap(attrs map[string][2]string, maxLen int, ellipsis string) map[string][2]string {
	if attrs == nil {
		return nil
	}
	result := make(map[string][2]string, len(attrs))
	for k, pair := range attrs {
		result[k] = [2]string{
			truncateString(pair[0], maxLen, ellipsis),
			truncateString(pair[1], maxLen, ellipsis),
		}
	}
	return result
}

func truncateString(s string, maxLen int, ellipsis string) string {
	if len(s) <= maxLen {
		return s
	}
	cutAt := maxLen - len(ellipsis)
	if cutAt < 0 {
		cutAt = 0
	}
	// Avoid splitting inside a multi-byte rune.
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	truncated := string(runes[:cutAt])
	// Append a hint showing original length.
	return fmt.Sprintf("%s%s[%d chars]", truncated, ellipsis, len(strings.TrimSpace(s)))
}
