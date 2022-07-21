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
	"github.com/gookit/color"
	"github.com/spf13/viper"
)

// ChromaNames returns the chroma style names in color.
func ChromaNames(w io.Writer, lexer string) {
	var s Names = styles.Names()
	s.String(w, true, lexer)
}

// ChromaNamesMono returns the chroma style names.
func ChromaNamesMono(w io.Writer, lexer string) {
	var s Names = styles.Names()
	s.String(w, false, lexer)
}

// CSS returns the element colored using CSS syntax highlights.
func CSS(w io.Writer, elm string) error {
	style := viper.GetString(get.Styleh)
	return Elm(w, elm, "css", style, color.Enable)
}

// Elm applies color syntax to an element.
func Elm(w io.Writer, elm, lexer, style string, color bool) error {
	if elm == "" {
		return nil
	}
	fmt.Fprintln(w)
	if err := str.HighlightWriter(w, elm, lexer, style, color); err != nil {
		return fmt.Errorf("%w, html %s, %s", logs.ErrHighlight, lexer, err)
	}
	fmt.Fprintln(w)
	return nil
}

// HTML returns the element colored using HTML syntax highlights.
func HTML(w io.Writer, elm string) error {
	style := viper.GetString(get.Styleh)
	return Elm(w, elm, "html", style, color.Enable)
}

// Names of the themes for the HighlightWriter.
type Names []string

// String lists and applies the named themes for the HighlightWriter.
func (n Names) String(w io.Writer, theme bool, lexer string) {
	if lexer == "json" {
		n.lexerJSON(w, theme)
		return
	}
	n.lexorOthers(w, theme, lexer)
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

func (n Names) lexorOthers(w io.Writer, theme bool, lexer string) {
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
	fmt.Fprint(w, strings.Join(s, ""))
}

func (n Names) lexerJSON(w io.Writer, theme bool) {
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
	fmt.Fprint(w, strings.Join(s, ""))
}
