// Package baseline provides save, load, and subtraction operations for drift
// report baselines.
//
// A baseline is a snapshot of drift items captured at a point in time and
// persisted as a JSON file. On subsequent runs the caller can load the saved
// baseline and call Subtract to surface only the drift that is NEW since the
// baseline was recorded — suppressing noise from known, long-standing drift.
//
// Typical workflow:
//
//	// First run – save a baseline
//	 baseline.Save("drift.baseline.json", items)
//
//	// Later run – show only new drift
//	 snap, _ := baseline.Load("drift.baseline.json")
//	 newItems := baseline.Subtract(currentItems, snap)
package baseline
