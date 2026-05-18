package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/stackdiff/internal/differ"
)

// Write outputs a human-readable drift report to the provided writer.
func Write(w io.Writer, result differ.DiffResult) error {
	if len(result.Added) == 0 && len(result.Removed) == 0 && len(result.Changed) == 0 {
		_, err := fmt.Fprintln(w, "✓ No drift detected. States are identical.")
		return err
	}

	fmt.Fprintln(w, "Drift Report")
	fmt.Fprintln(w, strings.Repeat("=", 40))

	if len(result.Added) > 0 {
		fmt.Fprintf(w, "\n+ Added (%d):\n", len(result.Added))
		for _, r := range result.Added {
			fmt.Fprintf(w, "  + %s (%s)\n", r.Address, r.Type)
		}
	}

	if len(result.Removed) > 0 {
		fmt.Fprintf(w, "\n- Removed (%d):\n", len(result.Removed))
		for _, r := range result.Removed {
			fmt.Fprintf(w, "  - %s (%s)\n", r.Address, r.Type)
		}
	}

	if len(result.Changed) > 0 {
		fmt.Fprintf(w, "\n~ Changed (%d):\n", len(result.Changed))
		for _, c := range result.Changed {
			fmt.Fprintf(w, "  ~ %s (%s):\n", c.Address, c.Type)
			for _, attr := range c.AttributeDiffs {
				fmt.Fprintf(w, "      %s: %q -> %q\n", attr.Key, attr.OldValue, attr.NewValue)
			}
		}
	}

	fmt.Fprintln(w, "\n"+strings.Repeat("=", 40))
	fmt.Fprintf(w, "Summary: +%d added, -%d removed, ~%d changed\n",
		len(result.Added), len(result.Removed), len(result.Changed))

	return nil
}
