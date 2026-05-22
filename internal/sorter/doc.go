// Package sorter provides deterministic ordering of drift report items.
//
// Items can be sorted by resource name, resource type, or change type.
// The Apply function always returns a new slice, leaving the original unchanged.
//
// Typical usage:
//
//	order, err := sorter.ParseSortOrder(cfg.SortBy)
//	if err != nil {
//		log.Fatal(err)
//	}
//	sorted := sorter.Apply(driftItems, order)
package sorter
