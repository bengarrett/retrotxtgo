package color

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/chroma/styles"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

// ColorCSS returns the element colored using CSS syntax highlights.
func ColorCSS(elm string) string {
	style := viper.GetString(get.Styleh)
	return ColorElm(elm, "css", style, true)
}

// ColorHTML returns the element colored using HTML syntax highlights.
func ColorHTML(elm string) string {
	style := viper.GetString(get.Styleh)
	return ColorElm(elm, "html", style, true)
}

// ColorElm applies color syntax to an element.
func ColorElm(elm, lexer, style string, color bool) string {
	if elm == "" {
		return ""
	}
	var b bytes.Buffer
	_ = io.Writer(&b)
	if err := str.HighlightWriter(&b, elm, lexer, style, color); err != nil {
		logs.FatalMark(fmt.Sprint("html ", lexer), logs.ErrHighlight, err)
	}
	return fmt.Sprintf("\n%s\n", b.String())
}

type Names []string

// ChromaNames returns the chroma style names.
func ChromaNames(lexer string) string {
	var s Names = styles.Names()
	return s.String(true, lexer)
}

// String lists and applies the named themes for the HighlightWriter.
func (n Names) String(theme bool, lexer string) string {
	maxWidth := 0
	for _, s := range n {
		if l := len(fmt.Sprintf("%s=%q", s, s)); l > maxWidth {
			maxWidth = l
		}
	}
	if !theme {
		return strings.Join(n, ", ")
	}
	s := make([]string, len(n))
	split := (len(n) / 2)
	const space = 2
	for i, name := range n {
		b, t := bytes.Buffer{}, ""
		pad := maxWidth - len(fmt.Sprintf("%s=%q", name, name))
		// prints a sequential list of styles
		if lexer == "json" {
			b = bytes.Buffer{}
			t = fmt.Sprintf("{ %q:%q }%s", name, name, strings.Repeat(" ", pad+space))
			if err := str.HighlightWriter(&b, t, lexer, name, true); err != nil {
				logs.FatalMark(name, logs.ErrHighlight, err)
			}
			s = append(s, fmt.Sprintf("%2d %s", i, b.String()))
			if split+i >= len(n) {
				break
			}
			b = bytes.Buffer{}
			t = fmt.Sprintf("{ %q:%q }\n", n[split+i], n[split+i])
			if err := str.HighlightWriter(&b, t, lexer, name, true); err != nil {
				logs.FatalMark(name, logs.ErrHighlight, err)
			}
			s = append(s, fmt.Sprintf("%2d %s", split+i, b.String()))
			continue
		}
		b = bytes.Buffer{}
		t = fmt.Sprintf("<%s=%q>%s", name, name, strings.Repeat(" ", pad+space))
		if err := str.HighlightWriter(&b, t, lexer, name, true); err != nil {
			logs.FatalMark(name, logs.ErrHighlight, err)
		}
		s = append(s, fmt.Sprintf("%2d %s", i, b.String()))
		if split+i >= len(n) {
			break
		}
		b = bytes.Buffer{}
		t = fmt.Sprintf("<%s=%q>\n", n[split+i], n[split+i])
		if err := str.HighlightWriter(&b, t, lexer, name, true); err != nil {
			logs.FatalMark(name, logs.ErrHighlight, err)
		}
		s = append(s, fmt.Sprintf("%2d %s", split+i, b.String()))
	}
	return strings.Join(s, "")
}
