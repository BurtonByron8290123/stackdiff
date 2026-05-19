package cache_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/stackdiff/internal/cache"
)

type payload struct {
	Value string `json:"value"`
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "src-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestGet_Miss_NonExistent(t *testing.T) {
	c := cache.New(t.TempDir())
	e, err := c.Get("/no/such/file.json")
	if err == nil && e != nil {
		t.Fatal("expected miss, got hit")
	}
}

func TestPutAndGet_Hit(t *testing.T) {
	src := writeTempFile(t, `{"hello":"world"}`)
	c := cache.New(t.TempDir())

	p := payload{Value: "terraform"}
	if err := c.Put(src, p); err != nil {
		t.Fatalf("Put: %v", err)
	}

	e, err := c.Get(src)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if e == nil {
		t.Fatal("expected cache hit, got miss")
	}

	var got payload
	if err := json.Unmarshal(e.Payload, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Value != "terraform" {
		t.Errorf("payload mismatch: got %q", got.Value)
	}
}

func TestGet_Stale_AfterFileChange(t *testing.T) {
	src := writeTempFile(t, `{"v":1}`)
	c := cache.New(t.TempDir())

	if err := c.Put(src, payload{Value: "old"}); err != nil {
		t.Fatal(err)
	}

	// Modify the source file so the hash changes.
	if err := os.WriteFile(src, []byte(`{"v":2}`), 0o644); err != nil {
		t.Fatal(err)
	}

	e, err := c.Get(src)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if e != nil {
		t.Fatal("expected stale miss, got hit")
	}
}

func TestPut_CreatesDir(t *testing.T) {
	src := writeTempFile(t, `{}`)
	dir := filepath.Join(t.TempDir(), "nested", "cache")
	c := cache.New(dir)

	if err := c.Put(src, payload{Value: "x"}); err != nil {
		t.Fatalf("Put: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("cache dir not created: %v", err)
	}
}
