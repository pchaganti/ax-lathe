package serve

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

const (
	lightStyle = "github"
	darkStyle  = "monokai"
)

func RenderMarkdown(src []byte) ([]byte, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle(lightStyle),
				highlighting.WithFormatOptions(
					html.WithClasses(true),
				),
			),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func HighlightCSS() (template.CSS, error) {
	formatter := html.New(html.WithClasses(true))
	var out bytes.Buffer

	light := styles.Get(lightStyle)
	if light == nil {
		return "", fmt.Errorf("chroma style %q not found", lightStyle)
	}
	if err := formatter.WriteCSS(&out, light); err != nil {
		return "", err
	}

	var darkBuf bytes.Buffer
	dark := styles.Get(darkStyle)
	if dark == nil {
		return "", fmt.Errorf("chroma style %q not found", darkStyle)
	}
	if err := formatter.WriteCSS(&darkBuf, dark); err != nil {
		return "", err
	}
	out.WriteString(scopeCSS(darkBuf.String(), `[data-theme="dark"]`))

	return template.CSS(out.String()), nil
}

// scopeCSS prefixes every CSS rule in src with the given selector. It assumes
// the chroma WriteCSS layout: one rule per line, each starting with a
// "/* ... */" comment followed by selector and declaration block.
func scopeCSS(src, prefix string) string {
	var b strings.Builder
	for _, line := range strings.Split(src, "\n") {
		if line == "" {
			b.WriteByte('\n')
			continue
		}
		end := strings.LastIndex(line, "*/")
		if end == -1 {
			b.WriteString(line)
			b.WriteByte('\n')
			continue
		}
		b.WriteString(line[:end+2])
		b.WriteByte(' ')
		b.WriteString(prefix)
		b.WriteString(line[end+2:])
		b.WriteByte('\n')
	}
	return b.String()
}
