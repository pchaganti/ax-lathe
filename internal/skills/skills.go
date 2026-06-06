// Package skills bundles the Lathe Claude Code skills into the binary so
// `lathe skills install` can write them out with no repo clone.
//
// The embedded copies under data/ are generated from the human-edited source at
// .claude/skills by `mage skills`; a `mage skillsCheck` parity gate (wired into
// `mage check`) keeps the two from drifting. go:embed cannot reach paths that
// begin with "." (so it can't read .claude/skills directly) -- that is the whole
// reason the data/ mirror exists.
package skills

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed data
var dataFS embed.FS

// Skill is one bundled skill: its slug (directory name), the raw SKILL.md bytes
// (Claude Code frontmatter included), and the name/description parsed from the
// YAML frontmatter.
type Skill struct {
	Slug        string
	Name        string
	Description string
	Raw         []byte
}

// All returns every bundled skill, sorted by slug. It returns an error only if
// the embedded tree is malformed (which would be a build-time bug).
func All() ([]Skill, error) {
	entries, err := fs.ReadDir(dataFS, "data")
	if err != nil {
		return nil, fmt.Errorf("read embedded skills: %w", err)
	}
	var out []Skill
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		slug := e.Name()
		raw, err := dataFS.ReadFile("data/" + slug + "/SKILL.md")
		if err != nil {
			return nil, fmt.Errorf("read embedded skill %q: %w", slug, err)
		}
		name, desc := parseFrontmatter(raw)
		if name == "" {
			name = slug
		}
		out = append(out, Skill{Slug: slug, Name: name, Description: desc, Raw: raw})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Slug < out[j].Slug })
	return out, nil
}

// parseFrontmatter pulls name: and description: out of a leading YAML
// frontmatter block delimited by "---" lines. The Lathe skills only use those
// two scalar keys, so a tiny line scanner beats pulling in a YAML dependency.
//
// If there is no well-formed frontmatter block (no leading fence, or an
// unclosed one), it returns empty strings rather than harvesting key-looking
// lines from the body -- this mirrors stripFrontmatter in cursor.go.
func parseFrontmatter(raw []byte) (name, description string) {
	// strings.Split always returns at least one element, so lines[0] is safe.
	lines := strings.Split(string(raw), "\n")
	// The first line must be exactly the opening fence (tolerate a trailing \r).
	if strings.TrimRight(lines[0], "\r") != "---" {
		return "", ""
	}
	closed := false
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "---" {
			closed = true
			break
		}
		if v, ok := frontmatterValue(line, "name"); ok {
			name = v
		} else if v, ok := frontmatterValue(line, "description"); ok {
			description = v
		}
	}
	if !closed {
		return "", ""
	}
	return name, description
}

// frontmatterValue returns the value for "key:" on a frontmatter line, trimming
// whitespace and surrounding quotes.
func frontmatterValue(line, key string) (string, bool) {
	prefix := key + ":"
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	v := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	v = strings.Trim(v, `"'`)
	return v, true
}
