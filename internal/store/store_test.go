package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devenjarvis/lathe/internal/store"
)

func TestStoreSingleTutorial(t *testing.T) {
	src := t.TempDir()
	if err := os.WriteFile(filepath.Join(src, "index.md"), []byte("# Hello"), 0644); err != nil {
		t.Fatal(err)
	}
	// Override home so we don't pollute the real ~/.lathe
	t.Setenv("HOME", t.TempDir())

	tut, err := store.Store(src, false)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}
	if tut.Series {
		t.Error("Store() Series = true, want false for single tutorial")
	}
	if tut.Status != store.StatusVerified {
		t.Errorf("Store() Status = %q, want %q", tut.Status, store.StatusVerified)
	}
}

func TestStoreSeriesTutorial(t *testing.T) {
	src := t.TempDir()
	for _, name := range []string{"part-01.md", "part-02.md", "part-03.md"} {
		if err := os.WriteFile(filepath.Join(src, name), []byte("# Part"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("HOME", t.TempDir())

	tut, err := store.Store(src, false)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}
	if !tut.Series {
		t.Error("Store() Series = false, want true for series")
	}
	if len(tut.Parts) != 3 {
		t.Errorf("Store() Parts = %v, want 3 parts", tut.Parts)
	}
}

func TestStoreVerifyingStatus(t *testing.T) {
	src := t.TempDir()
	if err := os.WriteFile(filepath.Join(src, "index.md"), []byte("# Hello"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", t.TempDir())
	// withVerify=true would try to spawn claude; we skip that by not passing it.
	// This test uses false — verify the status is "verified" (default) when not verifying.
	tut, err := store.Store(src, false)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}
	if tut.Status != store.StatusVerified {
		t.Errorf("Store() Status = %q, want %q", tut.Status, store.StatusVerified)
	}
}

func TestSlugToTitle(t *testing.T) {
	cases := []struct {
		slug  string
		title string
	}{
		{"digital-synth-zig", "Digital Synth Zig"},
		{"database-from-scratch-go", "Database From Scratch Go"},
		{"hello", "Hello"},
	}
	for _, c := range cases {
		got := store.SlugToTitle(c.slug)
		if got != c.title {
			t.Errorf("SlugToTitle(%q) = %q, want %q", c.slug, got, c.title)
		}
	}
}
