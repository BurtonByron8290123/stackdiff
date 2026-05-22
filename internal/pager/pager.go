// Package pager provides utilities for paginating drift report output,
// allowing large reports to be split into discrete pages.
package pager

import (
	"fmt"
	"strings"

	"github.com/stackdiff/stackdiff/internal/differ"
)

// Page represents a single page of drift items.
type Page struct {
	Number int
	Items  []differ.DriftItem
	Total  int
}

// Options controls pagination behaviour.
type Options struct {
	// PageSize is the maximum number of items per page. Zero disables pagination.
	PageSize int
}

// Paginate splits items into pages according to opts. If PageSize is zero or
// negative all items are returned as a single page.
func Paginate(items []differ.DriftItem, opts Options) []Page {
	if opts.PageSize <= 0 || len(items) == 0 {
		return []Page{{Number: 1, Items: items, Total: len(items)}}
	}

	var pages []Page
	for i := 0; i < len(items); i += opts.PageSize {
		end := i + opts.PageSize
		if end > len(items) {
			end = len(items)
		}
		pages = append(pages, Page{
			Number: len(pages) + 1,
			Items:  items[i:end],
			Total:  len(items),
		})
	}
	return pages
}

// Summary returns a human-readable description of a page's position.
func Summary(p Page, pageSize int) string {
	if pageSize <= 0 {
		return fmt.Sprintf("showing all %d item(s)", p.Total)
	}
	totalPages := (p.Total + pageSize - 1) / pageSize
	return fmt.Sprintf("page %d of %d (%d item(s))", p.Number, totalPages, p.Total)
}

// Header renders a simple text header for a page.
func Header(p Page, pageSize int) string {
	var b strings.Builder
	b.WriteString("--- ")
	b.WriteString(Summary(p, pageSize))
	b.WriteString(" ---")
	return b.String()
}
