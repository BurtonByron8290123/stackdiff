package summary

import (
	"fmt"

	"github.com/stackdiff/internal/differ"
)

// Stats holds aggregated counts for a drift report.
type Stats struct {
	Added    int
	Removed  int
	Changed  int
	Total    int
}

// Compute calculates drift statistics from a slice of DriftItems.
func Compute(items []differ.DriftItem) Stats {
	s := Stats{}
	for _, item := range items {
		switch item.ChangeType {
		case "added":
			s.Added++
		case "removed":
			s.Removed++
		case "changed":
			s.Changed++
		}
	}
	s.Total = s.Added + s.Removed + s.Changed
	return s
}

// Format returns a human-readable one-line summary string.
func (s Stats) Format() string {
	if s.Total == 0 {
		return "No drift detected."
	}
	return fmt.Sprintf(
		"Drift summary: %d change(s) total — %d added, %d removed, %d changed.",
		s.Total, s.Added, s.Removed, s.Changed,
	)
}

// HasDrift returns true when any changes were detected.
func (s Stats) HasDrift() bool {
	return s.Total > 0
}
