package str

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/alecthomas/chroma/quick"
	"github.com/alecthomas/chroma/styles"
	"github.com/gookit/color"
	"gopkg.in/yaml.v3"
)

// YamlExample is a YAML example
type YamlExample struct {
	Style struct {
		Name    string `yaml:"name"`
		Count   int    `yaml:"count"`
		Default bool   `yaml:"default"`
	}
}

func (s YamlExample) String(flag string) {
	fmt.Println()
	out, _ := yaml.Marshal(s)
	Highlight(string(out), "yaml", s.Style.Name)
	if flag != "" {
		fmt.Println(color.Secondary.Sprintf("%s=%q", flag, s.Style.Name))
	}
}

// YamlStyles prints out a list of available YAML color styles.
func YamlStyles(cmd string) {
	for i, s := range styles.Names() {
		var styles YamlExample
		styles.Style.Name = s
		styles.Style.Count = i
		if s == "monokai" {
			styles.Style.Default = true
		}
		styles.String(cmd)
	}
}

// Highlight and print the syntax of the source string except when piped to stdout.
func Highlight(source, lexer, style string) (err error) {
	return HighlightWriter(os.Stdout, source, lexer, style)
}

// HighlightWriter writes the highlight syntax of the source string except when piped to stdout.
func HighlightWriter(w io.Writer, source, lexer, style string) (err error) {
	var term = Term()
	// detect piping for text output or ansi for printing
	// source: https://stackoverflow.com/questions/43947363/detect-if-a-command-is-piped-or-not
	fo, err := os.Stdout.Stat()
	if err != nil {
		return err
	}
	// html json noop svg terminal terminal16m terminal256 tokens
	if term == "none" {
		fmt.Println(source)
	} else if (fo.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println(source)
	} else if err := quick.Highlight(w, source, lexer, term, style); err != nil {
		fmt.Println(source)
	}
	return nil
}

// Term determines the terminal type based on the COLORTERM and TERM environment variables.
func Term() (term string) {
	// 9.11.2 The environment variable TERM
	// https://www.gnu.org/software/gettext/manual/html_node/The-TERM-variable.html
	// Terminal Colors
	// https://gist.github.com/XVilka/8346728
	//
	term = "terminal256" // 256 colors (default)
	// first attempt to detect COLORTERM variable
	c := strings.TrimSpace(strings.ToLower(os.Getenv("COLORTERM")))
	switch c {
	case "24bit", "truecolor":
		return "terminal16m"
	}
	// then fallback to the -color suffix in TERM variable values
	t := strings.TrimSpace(strings.ToLower(os.Getenv("TERM")))
	s := strings.Split(t, "-")
	if len(s) > 1 {
		switch s[len(s)-1] {
		case "mono":
			return "none"
		case "color", "16color", "88color":
			return "terminal"
		case "256color":
			return term
		}
	}
	// otherwise do a direct match of the TERM variable value
	switch t {
	case "linux":
		return "none"
	case "konsole", "rxvt", "xterm", "vt100":
		return "terminal"
	}
	// anything else defaults to 256 colors
	return term
}

// UnderlineKeys uses ANSI to underline the first letter of each key.
func UnderlineKeys(keys []string) string {
	if len(keys) == 0 {
		return ""
	}
	sort.Strings(keys)
	for i, key := range keys {
		if utf8.RuneCountInString(key) > 1 {
			r, _ := utf8.DecodeRuneInString(key)
			c, err := underlineChar(string(r))
			if err != nil {
				keys[i] = key
			}
			keys[i] = fmt.Sprintf("%s%s", c, key[utf8.RuneLen(r):])
			if filepath.Ext(key) == ".min" {
				s := strings.Split(keys[i], ".")
				base := strings.Join(s[0:len(s)-1], ".")
				m, _ := underlineChar("m")
				keys[i] = fmt.Sprintf("%s.%sin", base, m)
			}
		} else {
			c, err := underlineChar(key)
			if err != nil {
				keys[i] = key
			} else {
				keys[i] = c
			}
		}
	}
	return strings.Join(keys, ", ")
}

func underlineChar(c string) (s string, err error) {
	if c == "" {
		return s, err
	}
	if !utf8.ValidString(c) {
		return s, errors.New("underlinechar: invalid utf-8 encoded rune")
	}
	var buf bytes.Buffer
	r, _ := utf8.DecodeRuneInString(c)
	t, err := template.New("underline").Parse("{{define \"TEXT\"}}\033[0m\033[4m{{.}}\033[0m{{end}}")
	if err != nil {
		return s, err
	}
	if err = t.ExecuteTemplate(&buf, "TEXT", string(r)); err != nil {
		return s, err
	}
	return buf.String(), nil
}
