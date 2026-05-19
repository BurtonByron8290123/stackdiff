package parser

import (
	"encoding/json"
	"fmt"

	"github.com/example/stackdiff/internal/cache"
)

// CachedParser wraps ParseStateFile with a [cache.Cache] so that
// unchanged state files are not re-parsed on every run.
type CachedParser struct {
	c *cache.Cache
}

// NewCachedParser returns a CachedParser backed by the given cache
// directory. Pass an empty string to disable caching (the zero-value
// Cache will always miss).
func NewCachedParser(cacheDir string) *CachedParser {
	return &CachedParser{c: cache.New(cacheDir)}
}

// Parse returns a TerraformState for path, reading from the cache when
// the file is unchanged since the last successful parse.
func (cp *CachedParser) Parse(path string) (*TerraformState, error) {
	entry, err := cp.c.Get(path)
	if err != nil {
		// Non-fatal: fall through to a fresh parse.
		fmt.Printf("cache lookup warning: %v\n", err)
	}

	if entry != nil {
		var s TerraformState
		if jsonErr := json.Unmarshal(entry.Payload, &s); jsonErr == nil {
			return &s, nil
		}
		// Corrupt cached payload — fall through.
	}

	s, parseErr := ParseStateFile(path)
	if parseErr != nil {
		return nil, parseErr
	}

	if putErr := cp.c.Put(path, s); putErr != nil {
		// Non-fatal: the caller still gets a valid result.
		fmt.Printf("cache store warning: %v\n", putErr)
	}
	return s, nil
}
