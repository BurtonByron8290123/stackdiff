// Package exporter provides functionality to write drift reports to files or stdout.
package exporter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Destination represents an output target for the drift report.
type Destination struct {
	// Path is the file path to write to. If empty, output goes to stdout.
	Path string
	// Overwrite controls whether an existing file should be overwritten.
	Overwrite bool
}

// Export writes the given content to the configured destination.
// If Path is empty, content is written to w (typically os.Stdout).
func Export(dest Destination, content []byte, w io.Writer) error {
	if dest.Path == "" {
		_, err := w.Write(content)
		return err
	}

	if err := ensureDir(dest.Path); err != nil {
		return fmt.Errorf("exporter: create directory: %w", err)
	}

	if !dest.Overwrite {
		if _, err := os.Stat(dest.Path); err == nil {
			return fmt.Errorf("exporter: file already exists: %s (use --overwrite to replace)", dest.Path)
		}
	}

	if err := os.WriteFile(dest.Path, content, 0o644); err != nil {
		return fmt.Errorf("exporter: write file: %w", err)
	}

	return nil
}

// ensureDir creates the parent directory of path if it does not exist.
func ensureDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}
