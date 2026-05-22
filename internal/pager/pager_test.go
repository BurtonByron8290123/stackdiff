package pager_test

import (
	"testing"

	"github.com/stackdiff/stackdiff/internal/differ"
	"github.com/stackdiff/stackdiff/internal/pager"
)

func makeItems(n int) []differ.DriftItem {
	items := make([]differ.DriftItem, n)
	for i := range items {
		items[i] = differ.DriftItem{
			Address:    fmt.Sprintf("aws_instance.web_%d", i),
			ChangeType: "modified",
		}
	}
	return items
}

import "fmt"

func TestPaginate_ZeroPageSize(t *testing.T) {
	items := makeItems(10)
	pages := pager.Paginate(items, pager.Options{PageSize: 0})
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(pages))
	}
	if len(pages[0].Items) != 10 {
		t.Errorf("expected 10 items on page 1, got %d", len(pages[0].Items))
	}
}

func TestPaginate_ExactMultiple(t *testing.T) {
	items := makeItems(6)
	pages := pager.Paginate(items, pager.Options{PageSize: 2})
	if len(pages) != 3 {
		t.Fatalf("expected 3 pages, got %d", len(pages))
	}
	for i, p := range pages {
		if len(p.Items) != 2 {
			t.Errorf("page %d: expected 2 items, got %d", i+1, len(p.Items))
		}
		if p.Number != i+1 {
			t.Errorf("page %d: unexpected number %d", i+1, p.Number)
		}
		if p.Total != 6 {
			t.Errorf("page %d: expected Total=6, got %d", i+1, p.Total)
		}
	}
}

func TestPaginate_Remainder(t *testing.T) {
	items := makeItems(7)
	pages := pager.Paginate(items, pager.Options{PageSize: 3})
	if len(pages) != 3 {
		t.Fatalf("expected 3 pages, got %d", len(pages))
	}
	if len(pages[2].Items) != 1 {
		t.Errorf("last page: expected 1 item, got %d", len(pages[2].Items))
	}
}

func TestPaginate_Empty(t *testing.T) {
	pages := pager.Paginate(nil, pager.Options{PageSize: 5})
	if len(pages) != 1 {
		t.Fatalf("expected 1 page for empty input, got %d", len(pages))
	}
	if len(pages[0].Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(pages[0].Items))
	}
}

func TestSummary_NoPagination(t *testing.T) {
	p := pager.Page{Number: 1, Total: 5}
	got := pager.Summary(p, 0)
	expected := "showing all 5 item(s)"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestSummary_WithPagination(t *testing.T) {
	p := pager.Page{Number: 2, Total: 10}
	got := pager.Summary(p, 3)
	expected := "page 2 of 4 (10 item(s))"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestHeader(t *testing.T) {
	p := pager.Page{Number: 1, Total: 4}
	got := pager.Header(p, 2)
	expected := "--- page 1 of 2 (4 item(s)) ---"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}
