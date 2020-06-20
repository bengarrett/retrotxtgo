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

// TestMode disables piping detection which conflicts with go test
var TestMode = false

// Border wraps text around a single line border.
func Border(text string) *bytes.Buffer {
	maxLen := 0
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		l := utf8.RuneCountInString(scanner.Text())
		if l > maxLen {
			maxLen = l
		}
	}
	maxLen = maxLen + 2
	scanner = bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanLines)
	var b bytes.Buffer
	fmt.Fprintln(&b, ("┌" + strings.Repeat("─", maxLen) + "┐"))
	for scanner.Scan() {
		l := utf8.RuneCountInString(scanner.Text())
		lp := ((maxLen - l) / 2)
		rp := lp
		// if lp/rp are X.5 decimal values, add 1 right padd to account for the uneven split
		if float32((maxLen-l)/2) != float32(float32(maxLen-l)/2) {
			rp = rp + 1
		}
		fmt.Fprintf(&b, "│%s%s%s│\n", strings.Repeat(" ", lp), scanner.Text(), strings.Repeat(" ", rp))
	}
	fmt.Fprintln(&b, "└"+strings.Repeat("─", maxLen)+"┘")
	return &b
}

// Center align text to the width.
func Center(text string, width int) string {
	w := (width - len(text)) / 2
	if w > 0 {
		return strings.Repeat("\u0020", w) + text
	}
	return text
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
	if term == "none" {
		// user disabled color output, but it doesn't disable ANSI output
		fmt.Fprintln(w, source)
	} else if !TestMode && (fo.Mode()&os.ModeCharDevice) == 0 {
		// disable colour when piping, this will also trigger with go test
		fmt.Fprintln(w, source)
	} else if err := quick.Highlight(w, source, lexer, term, style); err != nil {
		return err
	}
	return nil
}

// NumberizeKeys uses ANSI to underline and prefix a sequential number in front of each key.
func NumberizeKeys(keys []string) string {
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
	return strings.Join(s, ", ")
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

// UnderlineChar uses ANSI to underline the first character of a string.
func UnderlineChar(c string) (s string, err error) {
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

// UnderlineKeys uses ANSI to underline the first letter of each key.
func UnderlineKeys(keys []string) string {
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
					log.Fatal(err)
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

// JSONExample is used for previewing color themes
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
	out, _ := json.MarshalIndent(s, "", "  ")
	if flag != "" {
		fmt.Println("\n" + color.Secondary.Sprintf("%s=%q", flag, s.Style.Name))
	}
	Highlight(string(out), "json", s.Style.Name)
}

// JSONStyles prints out a list of available YAML color styles.
func JSONStyles(cmd string) {
	for i, s := range styles.Names() {
		var styles JSONExample
		styles.Style.Name = s
		styles.Style.Count = i
		if s == "dracula" {
			styles.Style.Default = true
		}
		styles.String(cmd)
	}
}
