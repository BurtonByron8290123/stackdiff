// Package pipeline wires together the full stackdiff processing pipeline:
// parse → diff → filter → sort → paginate → format → export.
//
// Callers construct a Pipeline from a Config and call Run to produce output.
package pipeline

import (
	"fmt"

	"github.com/yourorg/stackdiff/internal/baseline"
	"github.com/yourorg/stackdiff/internal/cache"
	"github.com/yourorg/stackdiff/internal/config"
	"github.com/yourorg/stackdiff/internal/differ"
	"github.com/yourorg/stackdiff/internal/exporter"
	"github.com/yourorg/stackdiff/internal/filter"
	"github.com/yourorg/stackdiff/internal/formatter"
	"github.com/yourorg/stackdiff/internal/pager"
	"github.com/yourorg/stackdiff/internal/parser"
	"github.com/yourorg/stackdiff/internal/sorter"
	"github.com/yourorg/stackdiff/internal/summary"
)

// Pipeline holds all dependencies needed to execute a drift comparison.
type Pipeline struct {
	cfg *config.Config
}

// New constructs a Pipeline from the given configuration.
// It returns an error if the configuration is invalid.
func New(cfg *config.Config) (*Pipeline, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &Pipeline{cfg: cfg}, nil
}

// Run executes the full pipeline and writes the result according to cfg.
// Steps:
//  1. Parse both state files (with optional caching).
//  2. Diff the two states.
//  3. Optionally subtract a saved baseline.
//  4. Filter by resource type, name prefix, or change kind.
//  5. Sort the drift items.
//  6. Paginate if a page size is configured.
//  7. Render to the requested format.
//  8. Export to stdout or a file.
func (p *Pipeline) Run() error {
	cfg := p.cfg

	// --- 1. Parse ---
	var cp *parser.CachedParser
	if cfg.CacheDir != "" {
		c, err := cache.New(cfg.CacheDir)
		if err != nil {
			return fmt.Errorf("cache init: %w", err)
		}
		cp = parser.NewCachedParser(c)
	}

	parseFile := func(path string) (*parser.StateFile, error) {
		if cp != nil {
			return cp.Parse(path)
		}
		return parser.ParseStateFile(path)
	}

	base, err := parseFile(cfg.BaseFile)
	if err != nil {
		return fmt.Errorf("parsing base state %q: %w", cfg.BaseFile, err)
	}
	target, err := parseFile(cfg.TargetFile)
	if err != nil {
		return fmt.Errorf("parsing target state %q: %w", cfg.TargetFile, err)
	}

	// --- 2. Diff ---
	items := differ.Compare(base, target)

	// --- 3. Baseline subtraction ---
	if cfg.BaselineFile != "" {
		bl, err := baseline.Load(cfg.BaselineFile)
		if err != nil {
			return fmt.Errorf("loading baseline %q: %w", cfg.BaselineFile, err)
		}
		items = baseline.Subtract(items, bl)
	}

	// --- 4. Filter ---
	items = filter.Apply(items, filter.Options{
		Types:       cfg.FilterTypes,
		NamePrefix:  cfg.FilterNamePrefix,
		ChangeTypes: cfg.FilterChangeTypes,
	})

	// --- 5. Sort ---
	order, err := sorter.ParseSortOrder(cfg.SortOrder)
	if err != nil {
		return fmt.Errorf("invalid sort order: %w", err)
	}
	items = sorter.Apply(items, order)

	// --- 6. Paginate ---
	var page []differ.DriftItem
	if cfg.PageSize > 0 {
		pages := pager.Paginate(items, cfg.PageSize)
		if cfg.Page < 1 || cfg.Page > len(pages) {
			page = []differ.DriftItem{}
		} else {
			page = pages[cfg.Page-1]
		}
	} else {
		page = items
	}

	// --- 7. Render ---
	fmt2, err := formatter.ParseFormat(cfg.Format)
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}
	sum := summary.Compute(items)
	output, err := formatter.Render(fmt2, page, sum)
	if err != nil {
		return fmt.Errorf("rendering output: %w", err)
	}

	// --- 8. Export ---
	return exporter.Export(output, exporter.Options{
		OutputFile: cfg.OutputFile,
		Overwrite:  cfg.Overwrite,
	})
}
