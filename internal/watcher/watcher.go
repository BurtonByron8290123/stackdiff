// Package watcher monitors Terraform state files for changes and triggers
// a re-diff when modifications are detected.
package watcher

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// Event is emitted when one or both watched files change.
type Event struct {
	Path    string
	OldHash string
	NewHash string
}

// Watcher polls a set of file paths at a given interval and sends Events
// on the Changes channel whenever a file's content hash changes.
type Watcher struct {
	Paths    []string
	Interval time.Duration
	Changes  chan Event

	hashes map[string]string
}

// New creates a Watcher for the given paths with the specified poll interval.
func New(paths []string, interval time.Duration) *Watcher {
	return &Watcher{
		Paths:    paths,
		Interval: interval,
		Changes:  make(chan Event, len(paths)),
		hashes:   make(map[string]string),
	}
}

// Start begins polling. It blocks until ctx is cancelled.
func (w *Watcher) Start(ctx context.Context) error {
	if err := w.snapshot(); err != nil {
		return fmt.Errorf("watcher: initial snapshot: %w", err)
	}
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.check(); err != nil {
				return fmt.Errorf("watcher: check: %w", err)
			}
		}
	}
}

func (w *Watcher) snapshot() error {
	for _, p := range w.Paths {
		h, err := hashFile(p)
		if err != nil {
			return err
		}
		w.hashes[p] = h
	}
	return nil
}

func (w *Watcher) check() error {
	for _, p := range w.Paths {
		h, err := hashFile(p)
		if err != nil {
			return err
		}
		if prev := w.hashes[p]; prev != h {
			w.Changes <- Event{Path: p, OldHash: prev, NewHash: h}
			w.hashes[p] = h
		}
	}
	return nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("hashFile: %w", err)
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hashFile: %w", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
