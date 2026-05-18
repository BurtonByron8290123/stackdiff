package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stackdiff/internal/differ"
)

func renderText(items []differ.DriftItem) string {
	if len(items) == 0 {
		return "No drift detected.\n"
	}
	var sb strings.Builder
	for _, item := range items {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", item.ChangeType, item.Address))
		for _, attr := range item.ChangedAttributes {
			sb.WriteString(fmt.Sprintf("  ~ %s: %q -> %q\n", attr.Key, attr.OldValue, attr.NewValue))
		}
	}
	sb.WriteString(fmt.Sprintf("\nSummary: %d change(s) detected.\n", len(items)))
	return sb.String()
}

func renderMarkdown(items []differ.DriftItem) string {
	if len(items) == 0 {
		return "## Drift Report\n\n_No drift detected._\n"
	}
	var sb strings.Builder
	sb.WriteString("## Drift Report\n\n")
	sb.WriteString("| Change | Address | Attributes |\n")
	sb.WriteString("|--------|---------|------------|\n")
	for _, item := range items {
		attrs := "-"
		if len(item.ChangedAttributes) > 0 {
			parts := make([]string, 0, len(item.ChangedAttributes))
			for _, a := range item.ChangedAttributes {
				parts = append(parts, fmt.Sprintf("`%s`", a.Key))
			}
			attrs = strings.Join(parts, ", ")
		}
		sb.WriteString(fmt.Sprintf("| %s | `%s` | %s |\n", item.ChangeType, item.Address, attrs))
	}
	sb.WriteString(fmt.Sprintf("\n**Total changes:** %d\n", len(items)))
	return sb.String()
}

func renderJSON(items []differ.DriftItem) (string, error) {
	type jsonReport struct {
		Changes []differ.DriftItem `json:"changes"`
		Total   int                `json:"total"`
	}
	report := jsonReport{Changes: items, Total: len(items)}
	if items == nil {
		report.Changes = []differ.DriftItem{}
	}
	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON report: %w", err)
	}
	return string(b) + "\n", nil
}
