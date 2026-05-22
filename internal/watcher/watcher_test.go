package watcher_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/your-org/stackdiff/internal/watcher"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "state-*.tfstate")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestWatcher_NoChangeProducesNoEvent(t *testing.T) {
	path := writeTempFile(t, `{"version":4}`)
	w := watcher.New([]string{path}, 20*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	go w.Start(ctx) //nolint:errcheck

	select {
	case ev := <-w.Changes:
		t.Errorf("unexpected event for unchanged file: %+v", ev)
	case <-ctx.Done():
		// expected: no events
	}
}

func TestWatcher_DetectsChange(t *testing.T) {
	path := writeTempFile(t, `{"version":4}`)
	w := watcher.New([]string{path}, 20*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go w.Start(ctx) //nolint:errcheck

	// Give the watcher time to take the initial snapshot.
	time.Sleep(30 * time.Millisecond)

	if err := os.WriteFile(path, []byte(`{"version":5}`), 0o644); err != nil {
		t.Fatalf("rewrite file: %v", err)
	}

	select {
	case ev := <-w.Changes:
		if ev.Path != path {
			t.Errorf("got path %q, want %q", ev.Path, path)
		}
		if ev.OldHash == ev.NewHash {
			t.Errorf("expected hashes to differ, both are %q", ev.OldHash)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcher_MultipleFiles(t *testing.T) {
	p1 := writeTempFile(t, `{"version":4}`)
	p2 := writeTempFile(t, `{"version":4}`)
	w := watcher.New([]string{p1, p2}, 20*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go w.Start(ctx) //nolint:errcheck
	time.Sleep(30 * time.Millisecond)

	if err := os.WriteFile(p2, []byte(`{"version":99}`), 0o644); err != nil {
		t.Fatalf("rewrite p2: %v", err)
	}

	select {
	case ev := <-w.Changes:
		if ev.Path != p2 {
			t.Errorf("expected event for p2, got %q", ev.Path)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for change event")
	}
}
