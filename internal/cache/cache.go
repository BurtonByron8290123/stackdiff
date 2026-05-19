// Package cache provides a simple file-based caching layer for parsed
// Terraform state files, keyed by file path and modification time.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Entry holds a cached parsed state along with metadata used for
// invalidation.
type Entry struct {
	FileHash string          `json:"file_hash"`
	Payload  json.RawMessage `json:"payload"`
}

// Cache manages on-disk JSON cache entries stored in a configurable
// directory.
type Cache struct {
	Dir string
}

// New returns a Cache that stores entries under dir.
func New(dir string) *Cache {
	return &Cache{Dir: dir}
}

// key derives a deterministic cache filename from the source path.
func (c *Cache) key(sourcePath string) string {
	h := sha256.Sum256([]byte(sourcePath))
	return filepath.Join(c.Dir, hex.EncodeToString(h[:16])+".json")
}

// Get returns the cached Entry for sourcePath if the file at sourcePath
// still has the same SHA-256 hash as when it was stored. Returns
// (nil, nil) on a miss.
func (c *Cache) Get(sourcePath string) (*Entry, error) {
	cachePath := c.key(sourcePath)
	data, err := os.ReadFile(cachePath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cache read: %w", err)
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, nil // treat corrupt entry as miss
	}
	current, err := fileHash(sourcePath)
	if err != nil {
		return nil, err
	}
	if current != e.FileHash {
		return nil, nil // stale
	}
	return &e, nil
}

// Put stores payload as the cached Entry for sourcePath, computing and
// recording the current file hash.
func (c *Cache) Put(sourcePath string, payload interface{}) error {
	if err := os.MkdirAll(c.Dir, 0o755); err != nil {
		return fmt.Errorf("cache mkdir: %w", err)
	}
	hash, err := fileHash(sourcePath)
	if err != nil {
		return err
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cache marshal: %w", err)
	}
	e := Entry{FileHash: hash, Payload: raw}
	out, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("cache marshal entry: %w", err)
	}
	return os.WriteFile(c.key(sourcePath), out, 0o644)
}

// fileHash returns the hex-encoded SHA-256 of the file at path.
func fileHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cache hash: %w", err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
