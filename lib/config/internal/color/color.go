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

// ChromaNames returns the chroma style names in color.
func ChromaNames(lexer string) string {
	var s Names = styles.Names()
	return s.String(true, lexer)
}

// ChromaNamesMono returns the chroma style names.
func ChromaNamesMono(lexer string) string {
	var s Names = styles.Names()
	return s.String(false, lexer)
}

// CSS returns the element colored using CSS syntax highlights.
func CSS(elm string) string {
	style := viper.GetString(get.Styleh)
	return Elm(elm, "css", style, true)
}

// Elm applies color syntax to an element.
func Elm(elm, lexer, style string, color bool) string {
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

// HTML returns the element colored using HTML syntax highlights.
func HTML(elm string) string {
	style := viper.GetString(get.Styleh)
	return Elm(elm, "html", style, true)
}

// Names of the themes for the HighlightWriter.
type Names []string

// String lists and applies the named themes for the HighlightWriter.
func (n Names) String(theme bool, lexer string) string {
	if lexer == "json" {
		return n.lexerJSON(theme)
	}
	return n.lexorOthers(theme, lexer)
}

func (n Names) maximunWidth() int {
	maxWidth := 0
	for _, ns := range n {
		if l := len(fmt.Sprintf("%s=%q", ns, ns)); l > maxWidth {
			maxWidth = l
		}
	}
	return maxWidth
}

func (n Names) lexorOthers(theme bool, lexer string) string {
	const space = 2
	maxWidth := n.maximunWidth()
	s := make([]string, len(n))
	split := (len(n) / space)
	for i, name := range n {
		if name == "" {
			continue
		}
		var (
			b bytes.Buffer
			t string
		)
		pad := maxWidth - len(fmt.Sprintf("%s=%q", name, name))
		// prints a sequential list of styles
		if split+i >= len(n) {
			break
		}
		b = bytes.Buffer{}
		t = fmt.Sprintf("<%s=%q>%s", name, name, strings.Repeat(" ", pad+space))
		if err := str.HighlightWriter(&b, t, lexer, name, theme); err != nil {
			logs.FatalMark(name, logs.ErrHighlight, err)
		}
		s[i] = fmt.Sprintf("%2d %s", i, b.String())
		if len(n) == 1 {
			break
		}
		if split+i >= len(n) {
			break
		}
		b = bytes.Buffer{}
		t = fmt.Sprintf("<%s=%q>\n", n[split+i], n[split+i])
		if err := str.HighlightWriter(&b, t, lexer, name, theme); err != nil {
			logs.FatalMark(name, logs.ErrHighlight, err)
		}
		s[i] = fmt.Sprintf("%s%2d %s", s[i], split+i, b.String())
	}
	return strings.Join(s, "")
}

func (n Names) lexerJSON(theme bool) string {
	const space = 2
	maxWidth := n.maximunWidth()
	s := make([]string, len(n))
	split := (len(n) / space)
	for i, name := range n {
		if name == "" {
			continue
		}
		var (
			b bytes.Buffer
			t string
		)
		pad := maxWidth - len(fmt.Sprintf("%s=%q", name, name))
		// prints a sequential list of styles
		b = bytes.Buffer{}
		t = fmt.Sprintf("{ %q:%q }%s", name, name, strings.Repeat(" ", pad+space))
		if err := str.HighlightWriter(&b, t, "json", name, theme); err != nil {
			logs.FatalMark(name, logs.ErrHighlight, err)
		}
		s[i] = fmt.Sprintf("%2d %s", i, b.String())
		if split+i >= len(n) {
			break
		}
		b = bytes.Buffer{}
		t = fmt.Sprintf("{ %q:%q }\n", n[split+i], n[split+i])
		if err := str.HighlightWriter(&b, t, "json", name, theme); err != nil {
			logs.FatalMark(name, logs.ErrHighlight, err)
		}
		s[i] = fmt.Sprintf("%s%2d %s", s[i], split+i, b.String())
		continue
	}
	return strings.Join(s, "")
}
