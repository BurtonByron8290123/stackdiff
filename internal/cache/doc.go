// Package cache implements a lightweight file-based cache for parsed
// Terraform state files.
//
// A [Cache] stores serialised JSON payloads on disk, keyed by the SHA-256
// hash of the source file's path. On retrieval the current hash of the
// source file is compared with the stored hash; if they differ the entry
// is treated as stale and a miss is returned, forcing the caller to
// re-parse and re-store.
//
// Typical usage:
//
//	c := cache.New(".stackdiff-cache")
//	if entry, _ := c.Get(path); entry != nil {
//		// fast path: deserialise entry.Payload
//	} else {
//		// slow path: parse, then c.Put(path, result)
//	}
package cache
