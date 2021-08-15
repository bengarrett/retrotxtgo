package str

import (
	"bufio"
	"bytes"
	"encoding/json"
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

type terminal int

const (
	// TermMono no color.
	TermMono terminal = iota
	// Term16 ANSI standard 16 color.
	Term16
	// Term88 XTerm 88 color.
	Term88
	// Term256 XTerm 256 color.
	Term256
	// Term16M ANSI high-color.
	Term16M
	none = "none"
	term = "terminal"
)

func (t terminal) String() string {
	return [...]string{none, term, term, "terminal256", "terminal16m"}[t]
}

// JSONExample is used to preview theme colors.
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
func Center(width int, text string) string {
	const split, space = 2, "\u0020"
	w := (width - len(text)) / split
	if w > 0 {
		return strings.Repeat(space, w) + text
	}
	return text
}

// GetEnv gets and formats the value of the environment variable named in the key.
func GetEnv(key string) string {
	return strings.TrimSpace(strings.ToLower(os.Getenv(key)))
}

// Highlight and print the syntax of the source string except when piped to stdout.
func Highlight(source, lexer, style string, ansi bool) (err error) {
	return HighlightWriter(os.Stdout, source, lexer, style, ansi)
}

// HighlightWriter writes the highlight syntax of the source string except when piped to stdout.
func HighlightWriter(w io.Writer, source, lexer, style string, ansi bool) (err error) {
	term := Term(GetEnv("COLORTERM"), GetEnv("TERM"))
	// detect piping for text output or ansi for printing
	// source: https://stackoverflow.com/questions/43947363/detect-if-a-command-is-piped-or-not
	fo, err := os.Stdout.Stat()
	if err != nil {
		return fmt.Errorf("highlight writer stdout error: %w", err)
	}
	if term == none {
		// user disabled color output, but it doesn't disable ANSI output
		fmt.Fprintln(w, source)
		return nil
	}
	if !ansi && (fo.Mode()&os.ModeCharDevice) == 0 {
		// disable color when piping or running unit tests
		fmt.Fprintln(w, source)
		return nil
	}
	if err := quick.Highlight(w, source, lexer, term, style); err != nil {
		return fmt.Errorf("highlight writer: %w", err)
	}
	return nil
}

// HR returns a horizontal ruler or line break.
func HR(width int) string {
	const horizontalBar = "\u2500"
	return fmt.Sprintf(" %s", Cb(strings.Repeat(horizontalBar, width)))
}

func HRPadded(width int) string {
	const horizontalBar = "\u2500"
	return fmt.Sprintf(" \n%s\n", Cb(strings.Repeat(horizontalBar, width)))
}

// NumberizeKeys uses ANSI to underline and prefix a sequential number in front of each key.
func NumberizeKeys(keys ...string) string {
	if len(keys) == 0 {
		return ""
	}
	const nbsp = "\u00A0"
	var s = make([]string, len(keys))
	sort.Strings(keys)
	for i, key := range keys {
		if i == 0 {
			s[i] = fmt.Sprintf("  Use %s for%s%s", Example(strconv.Itoa(i)), nbsp, key)
			continue
		}
		s[i] = fmt.Sprintf("      %s for%s%s", Example(strconv.Itoa(i)), nbsp, key)
	}
	return strings.Join(s, "\n")
}

// Term determines the terminal type based on the COLORTERM and TERM environment variables.
func Term(colorEnv, env string) string {
	// 9.11.2 The environment variable TERM
	// https://www.gnu.org/software/gettext/manual/html_node/The-TERM-variable.html
	// Terminal Colors
	// https://gist.github.com/XVilka/8346728

	// first, attempt to detect a COLORTERM variable
	switch colorEnv {
	case "24bit", "truecolor":
		return Term16M.String()
	}
	// then fallback to the -color suffix in TERM variable values
	s := strings.Split(env, "-")
	if len(s) > 1 {
		switch s[len(s)-1] {
		case "mono":
			return TermMono.String()
		case "color", "16color", "88color":
			return Term16.String()
		case "256color":
			return Term256.String()
		}
	}
	// otherwise do a direct match of the TERM variable value
	switch env {
	case "linux":
		return TermMono.String()
	case "konsole", "rxvt", "xterm", "vt100":
		return Term16.String()
	}
	return Term256.String()
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
			continue
		}
		c, err := UnderlineChar(key)
		if err != nil {
			keys[i] = key
			continue
		}
		keys[i] = c
	}
	return strings.Join(keys, ", ")
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
