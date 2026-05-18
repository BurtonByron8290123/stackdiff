package differ

import (
	"fmt"

	"github.com/stackdiff/internal/parser"
)

// ChangeType represents the kind of drift detected.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

// ResourceDiff describes a single resource-level difference.
type ResourceDiff struct {
	Address    string
	ChangeType ChangeType
	OldValue   map[string]interface{}
	NewValue   map[string]interface{}
}

// DriftReport holds the full comparison result.
type DriftReport struct {
	Diffs []ResourceDiff
}

// HasDrift returns true when at least one difference was found.
func (r *DriftReport) HasDrift() bool {
	return len(r.Diffs) > 0
}

// Compare computes the drift between a baseline and current Terraform state.
func Compare(baseline, current *parser.TerraformState) (*DriftReport, error) {
	if baseline == nil || current == nil {
		return nil, fmt.Errorf("both baseline and current states must be non-nil")
	}

	report := &DriftReport{}

	baselineMap := indexResources(baseline)
	currentMap := indexResources(current)

	// Detect removed and modified resources.
	for addr, baseRes := range baselineMap {
		if curRes, ok := currentMap[addr]; !ok {
			report.Diffs = append(report.Diffs, ResourceDiff{
				Address:    addr,
				ChangeType: Removed,
				OldValue:   baseRes.Values,
			})
		} else if !mapsEqual(baseRes.Values, curRes.Values) {
			report.Diffs = append(report.Diffs, ResourceDiff{
				Address:    addr,
				ChangeType: Modified,
				OldValue:   baseRes.Values,
				NewValue:   curRes.Values,
			})
		}
	}

	// Detect added resources.
	for addr, curRes := range currentMap {
		if _, ok := baselineMap[addr]; !ok {
			report.Diffs = append(report.Diffs, ResourceDiff{
				Address:    addr,
				ChangeType: Added,
				NewValue:   curRes.Values,
			})
		}
	}

	return report, nil
}

func indexResources(state *parser.TerraformState) map[string]*parser.Resource {
	m := make(map[string]*parser.Resource)
	for i := range state.Resources {
		r := &state.Resources[i]
		m[r.Address()] = r
	}
	return m
}

func mapsEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, va := range a {
		vb, ok := b[k]
		if !ok || fmt.Sprintf("%v", va) != fmt.Sprintf("%v", vb) {
			return false
		}
	}
	return true
}
