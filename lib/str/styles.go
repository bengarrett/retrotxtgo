package str

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/alecthomas/chroma/quick"
	"github.com/alecthomas/chroma/styles"
	"github.com/gookit/color"
)

var ErrRune = errors.New("invalid encoded rune")

type terminal int

const (
	// TermMono no colour.
	TermMono terminal = iota
	// Term16 ANSI standard 16 colour.
	Term16
	// Term88 XTerm 88 colour.
	Term88
	// Term256 XTerm 256 colour.
	Term256
	// Term16M ANSI high-colour.
	Term16M
)

// Formatter takes a terminal value and returns the quick.Highlight formatter value.
func (t terminal) Formatter() string {
	switch t {
	case TermMono:
		return "none"
	case Term16, Term88:
		return "terminal"
	case Term256:
		return "terminal256"
	case Term16M:
		return "terminal16m"
	}
	return ""
}

// Border wraps text around a single line border.
func Border(text string) *bytes.Buffer {
	const split = 2
	maxLen, scanner := 0, bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		l := utf8.RuneCountInString(scanner.Text())
		if l > maxLen {
			maxLen = l
		}
	}
	maxLen += split
	scanner = bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanLines)
	var b bytes.Buffer
	fmt.Fprintln(&b, ("┌" + strings.Repeat("─", maxLen) + "┐"))
	for scanner.Scan() {
		l := utf8.RuneCountInString(scanner.Text())
		lp := ((maxLen - l) / 2)
		rp := lp
		// if lp/rp are X.5 decimal values, add 1 right padd to account for the uneven split
		if float32((maxLen-l)/split) != float32(maxLen-l)/split {
			rp++
		}
		fmt.Fprintf(&b, "│%s%s%s│\n", strings.Repeat(" ", lp), scanner.Text(), strings.Repeat(" ", rp))
	}
	fmt.Fprintln(&b, "└"+strings.Repeat("─", maxLen)+"┘")
	return &b
}

// Center align text to the width.
func Center(text string, width int) string {
	const split = 2
	w := (width - len(text)) / split
	if w > 0 {
		return strings.Repeat("\u0020", w) + text
	}
	return text
}

// Highlight and print the syntax of the source string except when piped to stdout.
func Highlight(source, lexer, style string, ansi bool) (err error) {
	return HighlightWriter(os.Stdout, source, lexer, style, ansi)
}

// HighlightWriter writes the highlight syntax of the source string except when piped to stdout.
func HighlightWriter(w io.Writer, source, lexer, style string, ansi bool) (err error) {
	var term = Term()
	// detect piping for text output or ansi for printing
	// source: https://stackoverflow.com/questions/43947363/detect-if-a-command-is-piped-or-not
	fo, err := os.Stdout.Stat()
	if err != nil {
		return fmt.Errorf("highlight writer stdout error: %w", err)
	}
	if term == "none" {
		// user disabled color output, but it doesn't disable ANSI output
		fmt.Fprintln(w, source)
	} else if !ansi && (fo.Mode()&os.ModeCharDevice) == 0 {
		// disable colour when piping or running unit tests
		fmt.Fprintln(w, source)
	} else if err := quick.Highlight(w, source, lexer, term, style); err != nil {
		return fmt.Errorf("highlight writer: %w", err)
	}
	return nil
}

// NumberizeKeys uses ANSI to underline and prefix a sequential number in front of each key.
func NumberizeKeys(keys ...string) string {
	if len(keys) == 0 {
		return ""
	}
	var s = make([]string, len(keys))
	sort.Strings(keys)
	for i, key := range keys {
		n, err := UnderlineChar(strconv.Itoa(i))
		if err != nil {
			log.Fatal(err)
		}
		s[i] = fmt.Sprintf("%s)\u00a0%s", n, key)
	}
	return strings.Join(s, "\n")
}

// Term determines the terminal type based on the COLORTERM and TERM environment variables.
func Term() string {
	// 9.11.2 The environment variable TERM
	// https://www.gnu.org/software/gettext/manual/html_node/The-TERM-variable.html
	// Terminal Colors
	// https://gist.github.com/XVilka/8346728
	var t terminal = -1
	// first attempt to detect COLORTERM variable
	c := strings.TrimSpace(strings.ToLower(os.Getenv("COLORTERM")))
	switch c {
	case "24bit", "truecolor":
		t = Term16M
		return t.Formatter()
	}
	// then fallback to the -color suffix in TERM variable values
	env := strings.TrimSpace(strings.ToLower(os.Getenv("TERM")))
	s := strings.Split(env, "-")
	if len(s) > 1 {
		switch s[len(s)-1] {
		case "mono":
			t = TermMono
			return t.Formatter()
		case "color", "16color", "88color":
			t = Term16
			return t.Formatter()
		case "256color":
			t = Term256
			return t.Formatter()
		}
	}
	// otherwise do a direct match of the TERM variable value
	switch env {
	case "linux":
		t = TermMono
		return t.Formatter()
	case "konsole", "rxvt", "xterm", "vt100":
		t = Term16
		return t.Formatter()
	}
	t = Term256
	return t.Formatter()
}

// UnderlineChar uses ANSI to underline the first character of a string.
func UnderlineChar(c string) (s string, err error) {
	if c == "" {
		return "", nil
	}
	if !utf8.ValidString(c) {
		return s, fmt.Errorf("underlinechar %q: %w", c, ErrRune)
	}
	var buf bytes.Buffer
	r, _ := utf8.DecodeRuneInString(c)
	t, err := template.New("underline").Parse("{{define \"TEXT\"}}\033[0m\033[4m{{.}}\033[0m{{end}}")
	if err != nil {
		return s, fmt.Errorf("underlinechar new template: %w", err)
	}
	if err = t.ExecuteTemplate(&buf, "TEXT", string(r)); err != nil {
		return s, fmt.Errorf("underlinechar execute template: %w", err)
	}
	return buf.String(), nil
}

// UnderlineKeys uses ANSI to underline the first letter of each key.
func UnderlineKeys(keys ...string) string {
	if len(keys) == 0 {
		return ""
	}
	sort.Strings(keys)
	for i, key := range keys {
		if utf8.RuneCountInString(key) > 1 {
			r, _ := utf8.DecodeRuneInString(key)
			c, err := UnderlineChar(string(r))
			if err != nil {
				keys[i] = key
			}
			keys[i] = fmt.Sprintf("%s%s", c, key[utf8.RuneLen(r):])
			if filepath.Ext(key) == ".min" {
				s := strings.Split(keys[i], ".")
				base := strings.Join(s[0:len(s)-1], ".")
				m, err := UnderlineChar("m")
				if err != nil {
					// must use standard log package
					log.Fatal("underline keys", keys, err)
				}
				keys[i] = fmt.Sprintf("%s.%sin", base, m)
			}
		} else {
			c, err := UnderlineChar(key)
			if err != nil {
				keys[i] = key
			} else {
				keys[i] = c
			}
		}
	}
	return strings.Join(keys, ", ")
}

// JSONExample is used for previewing color themes.
type JSONExample struct {
	Style struct {
		Name    string `json:"name"`
		Count   int    `json:"count"`
		Default bool   `json:"default"`
	}
}

func (s JSONExample) String(flag string) {
	fmt.Println()
	// config is stored as YAML but printed as JSON
	out, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatalln(fmt.Errorf("json example marshal indent: %w", err))
	}
	if flag != "" {
		fmt.Println("\n" + color.Secondary.Sprintf("%s=%q", flag, s.Style.Name))
	}
	if err := Highlight(string(out), "json", s.Style.Name, true); err != nil {
		// cannot use the logs package as it causes an import cycle error
		log.Fatalln(fmt.Errorf("json example highlight syntax error: %w", err))
	}
}

// JSONStyles prints out a list of available YAML color styles.
func JSONStyles(cmd string) {
	for i, s := range styles.Names() {
		var example JSONExample
		example.Style.Name, example.Style.Count = s, i
		if s == "dracula" {
			example.Style.Default = true
		}
		example.String(cmd)
	}
	fmt.Println()
}

// Valid checks style name validity.
func Valid(style string) bool {
	for _, n := range styles.Names() {
		if n == style {
			return true
		}
	}
	return false
}
