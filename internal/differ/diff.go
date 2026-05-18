package differ

// ChangeType constants describe how a resource has changed between states.
const (
	ChangeAdded   = "added"
	ChangeRemoved = "removed"
	ChangeModified = "modified"
	ChangeUnchanged = "unchanged"
)

// DriftItem represents a single resource that has drifted between two states.
type DriftItem struct {
	ResourceType string
	ResourceName string
	ChangeType   string
	Diffs        map[string][2]interface{} // key -> [oldVal, newVal]
}

// DriftReport is the result of comparing two state files.
type DriftReport struct {
	Items   []DriftItem
	Summary Summary
}

// Summary holds aggregate counts of drift changes.
type Summary struct {
	Added    int
	Removed  int
	Modified int
	Total    int
}

// State represents a parsed Terraform state file (minimal structure).
type State struct {
	Resources []Resource
}

// Resource represents a single resource entry within a Terraform state.
type Resource struct {
	Type       string
	Name       string
	Attributes map[string]interface{}
}

// Compare produces a DriftReport by comparing an old state to a new state.
func Compare(old, new State) DriftReport {
	var report DriftReport

	oldIdx := indexResources(old.Resources)
	newIdx := indexResources(new.Resources)

	for key, oldRes := range oldIdx {
		newRes, exists := newIdx[key]
		if !exists {
			report.Items = append(report.Items, DriftItem{
				ResourceType: oldRes.Type,
				ResourceName: oldRes.Name,
				ChangeType:   ChangeRemoved,
			})
			report.Summary.Removed++
			continue
		}
		if diffs := diffAttributes(oldRes.Attributes, newRes.Attributes); len(diffs) > 0 {
			report.Items = append(report.Items, DriftItem{
				ResourceType: oldRes.Type,
				ResourceName: oldRes.Name,
				ChangeType:   ChangeModified,
				Diffs:        diffs,
			})
			report.Summary.Modified++
		}
	}

	for key, newRes := range newIdx {
		if _, exists := oldIdx[key]; !exists {
			report.Items = append(report.Items, DriftItem{
				ResourceType: newRes.Type,
				ResourceName: newRes.Name,
				ChangeType:   ChangeAdded,
			})
			report.Summary.Added++
		}
	}

	report.Summary.Total = report.Summary.Added + report.Summary.Removed + report.Summary.Modified
	return report
}

func indexResources(resources []Resource) map[string]Resource {
	idx := make(map[string]Resource, len(resources))
	for _, r := range resources {
		idx[r.Type+"."+r.Name] = r
	}
	return idx
}

func diffAttributes(old, new map[string]interface{}) map[string][2]interface{} {
	diffs := make(map[string][2]interface{})
	for k, oldVal := range old {
		newVal, exists := new[k]
		if !exists || newVal != oldVal {
			diffs[k] = [2]interface{}{oldVal, newVal}
		}
	}
	for k, newVal := range new {
		if _, exists := old[k]; !exists {
			diffs[k] = [2]interface{}{nil, newVal}
		}
	}
	return diffs
}

func mapsEqual(a, b map[string]interface{}) bool {
	return len(diffAttributes(a, b)) == 0
}
