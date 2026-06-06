package skills

import (
	"strings"
	"testing"
)

func TestCursorCommandStripsFrontmatterKeepsBody(t *testing.T) {
	s := Skill{
		Slug:        "lathe",
		Name:        "lathe",
		Description: "Generate hands-on technical tutorials.",
		Raw: []byte("---\nname: lathe\ndescription: Generate.\n---\n\n" +
			"# Lathe\n\nRun `lathe store --tag foo`.\n"),
	}
	out := CursorCommand(s)

	if strings.Contains(out, "---") {
		t.Errorf("Cursor output should not contain YAML frontmatter fences:\n%s", out)
	}
	if !strings.HasPrefix(out, "# /lathe\n") {
		t.Errorf("expected a /lathe title header, got:\n%s", out)
	}
	if !strings.Contains(out, "Generate hands-on technical tutorials.") {
		t.Errorf("expected the description header, got:\n%s", out)
	}
	if !strings.Contains(out, "lathe store --tag foo") {
		t.Errorf("body bash invocation was dropped:\n%s", out)
	}
}

func TestCursorCommandRealSkillsPreserveInvocations(t *testing.T) {
	all, err := All()
	if err != nil {
		t.Fatalf("All() error: %v", err)
	}
	byslug := map[string]Skill{}
	for _, s := range all {
		byslug[s.Slug] = s
	}

	// Spot-check that key CLI handoffs survive translation.
	checks := map[string]string{
		"lathe":        "lathe store",
		"lathe-verify": "lathe verify-result",
	}
	for slug, want := range checks {
		s, ok := byslug[slug]
		if !ok {
			t.Fatalf("skill %q not found", slug)
		}
		out := CursorCommand(s)
		if strings.HasPrefix(out, "---") {
			t.Errorf("%s: frontmatter not stripped", slug)
		}
		if !strings.Contains(out, want) {
			t.Errorf("%s: expected %q in Cursor output", slug, want)
		}
	}
}

func TestCursorFilename(t *testing.T) {
	if got := CursorFilename(Skill{Slug: "lathe-ask"}); got != "lathe-ask.md" {
		t.Errorf("CursorFilename = %q, want lathe-ask.md", got)
	}
}
