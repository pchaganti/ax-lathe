package skills

import (
	"strings"
	"testing"
)

func TestAllReturnsEverySkillWithMetadata(t *testing.T) {
	all, err := All()
	if err != nil {
		t.Fatalf("All() error: %v", err)
	}
	const want = 5
	if len(all) != want {
		t.Fatalf("All() returned %d skills, want %d", len(all), want)
	}

	wantSlugs := map[string]bool{
		"lathe":        false,
		"lathe-ask":    false,
		"lathe-extend": false,
		"lathe-tag":    false,
		"lathe-verify": false,
	}
	for _, s := range all {
		if _, ok := wantSlugs[s.Slug]; !ok {
			t.Errorf("unexpected slug %q", s.Slug)
			continue
		}
		wantSlugs[s.Slug] = true
		if s.Name == "" {
			t.Errorf("skill %q has empty name", s.Slug)
		}
		if strings.TrimSpace(s.Description) == "" {
			t.Errorf("skill %q has empty description", s.Slug)
		}
		if len(s.Raw) == 0 {
			t.Errorf("skill %q has empty raw bytes", s.Slug)
		}
	}
	for slug, found := range wantSlugs {
		if !found {
			t.Errorf("missing expected skill %q", slug)
		}
	}
}

func TestParseFrontmatter(t *testing.T) {
	raw := []byte("---\nname: lathe\ndescription: Generate tutorials.\n---\n\n# Body\n")
	name, desc := parseFrontmatter(raw)
	if name != "lathe" {
		t.Errorf("name = %q, want lathe", name)
	}
	if desc != "Generate tutorials." {
		t.Errorf("description = %q", desc)
	}
}

func TestParseFrontmatterNoFrontmatter(t *testing.T) {
	name, desc := parseFrontmatter([]byte("# Just a heading\n"))
	if name != "" || desc != "" {
		t.Errorf("expected empty name/desc, got %q / %q", name, desc)
	}
}

func TestParseFrontmatterUnclosedFenceIsIgnored(t *testing.T) {
	// Opening fence but no closing one: must not harvest body lines.
	raw := []byte("---\nname: real\n\n# Body\nname: not-frontmatter\n")
	name, desc := parseFrontmatter(raw)
	if name != "" || desc != "" {
		t.Errorf("unclosed frontmatter should yield empty name/desc, got %q / %q", name, desc)
	}
}
