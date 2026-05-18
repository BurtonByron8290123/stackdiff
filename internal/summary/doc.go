// Package summary provides utilities for computing and formatting
// aggregated statistics over a set of drift items produced by the
// differ package.
//
// Typical usage:
//
//	items := differ.Compare(stateA, stateB)
//	stats := summary.Compute(items)
//	fmt.Println(stats.Format())
//	if stats.HasDrift() {
//		os.Exit(1)
//	}
package summary
