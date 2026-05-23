// Package highlighter provides ANSI colour helpers for drift report output.
// It wraps change-type labels in terminal escape codes when colour is enabled.
package highlighter

import "fmt"

// Palette holds the ANSI escape codes used for each change type.
type Palette struct {
	Added   string
	Removed string
	Changed string
	Reset   string
}

// DefaultPalette returns a Palette with standard ANSI colour codes.
func DefaultPalette() Palette {
	return Palette{
		Added:   "\033[32m", // green
		Removed: "\033[31m", // red
		Changed: "\033[33m", // yellow
		Reset:   "\033[0m",
	}
}

// Highlighter applies colour to change-type labels.
type Highlighter struct {
	enabled bool
	palette Palette
}

// New creates a Highlighter. When enabled is false all methods return plain
// strings without escape codes.
func New(enabled bool) *Highlighter {
	return &Highlighter{enabled: enabled, palette: DefaultPalette()}
}

// Label wraps a change-type string ("added", "removed", "changed") in the
// appropriate ANSI colour. Unknown types are returned unchanged.
func (h *Highlighter) Label(changeType string) string {
	if !h.enabled {
		return changeType
	}
	switch changeType {
	case "added":
		return fmt.Sprintf("%s%s%s", h.palette.Added, changeType, h.palette.Reset)
	case "removed":
		return fmt.Sprintf("%s%s%s", h.palette.Removed, changeType, h.palette.Reset)
	case "changed":
		return fmt.Sprintf("%s%s%s", h.palette.Changed, changeType, h.palette.Reset)
	default:
		return changeType
	}
}

// Value wraps an arbitrary string value in the colour associated with the
// given change type, then resets. Useful for colouring attribute values.
func (h *Highlighter) Value(changeType, value string) string {
	if !h.enabled {
		return value
	}
	var code string
	switch changeType {
	case "added":
		code = h.palette.Added
	case "removed":
		code = h.palette.Removed
	case "changed":
		code = h.palette.Changed
	default:
		return value
	}
	return fmt.Sprintf("%s%s%s", code, value, h.palette.Reset)
}
