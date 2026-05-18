package formatter

import (
	"fmt"
	"strings"

	"github.com/stackdiff/internal/differ"
)

// Format controls the output format of the drift report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatMarkdown Format = "markdown"
)

// ParseFormat parses a format string into a Format value.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "markdown", "md":
		return FormatMarkdown, nil
	default:
		return "", fmt.Errorf("unknown format %q: must be one of text, json, markdown", s)
	}
}

// Render converts a slice of DriftItems into a formatted string.
func Render(items []differ.DriftItem, format Format) (string, error) {
	switch format {
	case FormatJSON:
		return renderJSON(items)
	case FormatMarkdown:
		return renderMarkdown(items), nil
	default:
		return renderText(items), nil
	}
}
