package skills

import (
	"fmt"
	"strings"
)

// CursorCommand renders a skill as a Cursor command file.
//
// Cursor invokes the file at .cursor/commands/<slug>.md as the slash command
// /<slug>. Cursor does not use Claude Code's YAML frontmatter, so we strip it
// and prepend a plain Markdown title + description header. The body (the bash
// `lathe ...` calls and prose guidance) is preserved verbatim -- it reads fine
// as Cursor instructions. We deliberately do not port the interactive-handoff
// runtime model here; trigger + body only (see CLAUDE.md / the install caveats).
func CursorCommand(s Skill) string {
	body := stripFrontmatter(string(s.Raw))

	var b strings.Builder
	fmt.Fprintf(&b, "# /%s\n\n", s.Slug)
	if s.Description != "" {
		fmt.Fprintf(&b, "%s\n\n", s.Description)
	}
	b.WriteString(body)
	if !strings.HasSuffix(body, "\n") {
		b.WriteString("\n")
	}
	return b.String()
}

// CursorFilename maps a skill slug to its Cursor command filename.
func CursorFilename(s Skill) string {
	return s.Slug + ".md"
}

// stripFrontmatter removes a leading "---"-delimited YAML block (if present) and
// returns the remaining body with leading blank lines trimmed.
func stripFrontmatter(s string) string {
	// strings.Split always returns at least one element, so lines[0] is safe.
	lines := strings.Split(s, "\n")
	// The first line must be exactly the opening fence (tolerate a trailing \r).
	if strings.TrimRight(lines[0], "\r") != "---" {
		return s
	}
	// Find the closing fence after the opening one.
	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}
	if end == -1 {
		// Malformed frontmatter: leave the content untouched.
		return s
	}
	body := strings.Join(lines[end+1:], "\n")
	return strings.TrimLeft(body, "\n")
}
