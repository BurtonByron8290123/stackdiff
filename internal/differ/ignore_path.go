package differ

import (
	"strings"

	"github.com/you/stackdiff/internal/differ/model"
)

// PathIgnoreList holds glob-style attribute path patterns that should be
// excluded from drift comparison (e.g. "tags.*", "metadata.annotations.*").
type PathIgnoreList struct {
	patterns []string
}

// NewPathIgnoreList creates a PathIgnoreList from a slice of path patterns.
// Patterns are matched against dot-separated attribute paths.
// A trailing ".*" wildcard matches any sub-key at that level.
func NewPathIgnoreList(patterns []string) *PathIgnoreList {
	normalized := make([]string, 0, len(patterns))
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p != "" {
			normalized = append(normalized, p)
		}
	}
	return &PathIgnoreList{patterns: normalized}
}

// FilterAttributeDiffs removes attribute diffs whose keys match any of the
// registered path patterns and returns the filtered slice.
func (pl *PathIgnoreList) FilterAttributeDiffs(diffs []model.AttributeDiff) []model.AttributeDiff {
	if len(pl.patterns) == 0 {
		return diffs
	}
	out := diffs[:0:0]
	for _, d := range diffs {
		if !pl.matchesAny(d.Key) {
			out = append(out, d)
		}
	}
	return out
}

// matchesAny returns true if key matches at least one pattern.
func (pl *PathIgnoreList) matchesAny(key string) bool {
	for _, pat := range pl.patterns {
		if matchPath(pat, key) {
			return true
		}
	}
	return false
}

// matchPath matches a single pattern against a key.
// "foo.bar" matches exactly "foo.bar".
// "foo.*" matches "foo.anything" but not "foo.bar.baz".
// "*" matches any single-segment key.
func matchPath(pattern, key string) bool {
	if pattern == key {
		return true
	}
	if strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, ".*")
		if strings.HasPrefix(key, prefix+".") {
			// ensure no further dots after the wildcard segment
			rest := strings.TrimPrefix(key, prefix+".")
			if !strings.Contains(rest, ".") {
				return true
			}
		}
	}
	if pattern == "*" && !strings.Contains(key, ".") {
		return true
	}
	return false
}
