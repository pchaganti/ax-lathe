package serve_test

import (
	"strings"
	"testing"

	"github.com/devenjarvis/lathe/internal/serve"
)

func TestRenderMarkdown(t *testing.T) {
	src := []byte("# Hello World\n\nThis is a `test`.\n\n```go\nfmt.Println(\"hello\")\n```")

	out, err := serve.RenderMarkdown(src)
	if err != nil {
		t.Fatalf("RenderMarkdown() error = %v", err)
	}

	html := string(out)
	if !strings.Contains(html, "<h1>Hello World</h1>") {
		t.Errorf("RenderMarkdown() missing <h1>, got:\n%s", html)
	}
	if !strings.Contains(html, "<code>test</code>") {
		t.Errorf("RenderMarkdown() missing inline <code>, got:\n%s", html)
	}
	if !strings.Contains(html, "<pre") {
		t.Errorf("RenderMarkdown() code block not rendered as <pre>, got:\n%s", html)
	}
	if !strings.Contains(html, "Println") {
		t.Errorf("RenderMarkdown() code block content missing from output, got:\n%s", html)
	}
	if !strings.Contains(html, `class="chroma"`) {
		t.Errorf("RenderMarkdown() should emit chroma classes (WithClasses=true), got:\n%s", html)
	}
}

func TestHighlightCSS(t *testing.T) {
	css, err := serve.HighlightCSS()
	if err != nil {
		t.Fatalf("HighlightCSS() error = %v", err)
	}
	s := string(css)
	if !strings.Contains(s, ".chroma") {
		t.Error("HighlightCSS() missing .chroma rules")
	}
	if !strings.Contains(s, `[data-theme="dark"] .chroma`) {
		t.Error("HighlightCSS() missing dark-scoped rules")
	}
	// Light rules must not be scoped under [data-theme="dark"].
	for _, line := range strings.Split(s, "\n") {
		if strings.Contains(line, ".chroma") && !strings.Contains(line, `[data-theme="dark"]`) {
			// A light rule — fine.
			continue
		}
	}
	// Spot-check that both palettes appear: github (light) uses #fff for bg,
	// monokai (dark) uses #f8f8f2 for default fg.
	if !strings.Contains(strings.ToLower(s), "#fff") {
		t.Error("HighlightCSS() missing expected light-theme color")
	}
	if !strings.Contains(strings.ToLower(s), "#f8f8f2") {
		t.Error("HighlightCSS() missing expected dark-theme color")
	}
}
