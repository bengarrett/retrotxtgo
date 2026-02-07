// Package term provides colors and text formatting for the terminal.
package term

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/gookit/color"
)

var ErrRune = errors.New("invalid encoded rune")

const (
	// HBar is a the Unicode horizontal bar character.
	HBar = "\u2500"
	none = "none"
	term = "terminal"
)

// Terminal color support.
type Terminal int

const (
	TermMono Terminal = iota // monochrome with no color
	Term16                   // ANSI standard 16 color
	Term88                   // XTerm with 88 colors
	Term256                  // XTerm with 256 colors
	Term16M                  // ANSI high-color with 16 million colors
)

// String returns the terminal as a named color value.
func (t Terminal) String() string {
	return [...]string{none, term, term, "terminal256", "terminal16m"}[t]
}

// Border wraps the string around a single line border.
func Border(w io.Writer, s string) {
	if w == nil {
		w = io.Discard
	}
	const split = 2
	maxLen, scanner := 0, bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		l := utf8.RuneCountInString(scanner.Text())
		if l > maxLen {
			maxLen = l
		}
	}
	maxLen += split
	scanner = bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanLines)
	fmt.Fprintln(w, ("┌" + strings.Repeat("─", maxLen) + "┐"))
	for scanner.Scan() {
		l := utf8.RuneCountInString(scanner.Text())
		lp := ((maxLen - l) / split)
		rp := lp
		// if lp/rp are X.5 decimal values, add 1 right padd to account for the uneven split
		if float32((maxLen-l)/split) != float32(maxLen-l)/split {
			rp++
		}
		fmt.Fprintf(w, "│%s%s%s│\n", strings.Repeat(" ", lp), scanner.Text(), strings.Repeat(" ", rp))
	}
	fmt.Fprintln(w, "└"+strings.Repeat("─", maxLen)+"┘")
}

// Center align text to a the width of an area.
// If the width is less than the length of the string, the string is returned.
// There is no padding after the string.
func Center(width int, s string) string {
	const split, space = 2, "\u0020"
	if w := (width - len(s)) / split; w > 0 {
		return strings.Repeat(space, w) + s
	}
	return s
}

// GetEnv gets and formats the value of the environment variable named in the key.
func GetEnv(key string) string {
	return strings.TrimSpace(strings.ToLower(os.Getenv(key)))
}

// IsTerminal reports whether the output is a terminal.
func IsTerminal() bool {
	return color.Enable
}

// Head returns a colored and underlined string for use as a header.
// Provide a fixed width value for the underline border or set to zero.
// The header is colored with the fuzzy color.
func Head(w io.Writer, width int, s string) {
	if w == nil {
		w = io.Discard
	}
	r := color.OpFuzzy.Sprint(strings.Repeat(HBar, width))
	h := color.Primary.Sprint(Center(width, s))
	fmt.Fprintln(w, r)
	fmt.Fprintln(w, h)
}

// HR returns a horizontal ruler and a line break.
func HR(w io.Writer, width int) {
	if w == nil {
		w = io.Discard
	}
	fmt.Fprintf(w, " %s\n", Secondary(strings.Repeat(HBar, width)))
}

// Term determines the terminal type based on the COLORTERM and TERM environment variables.
//
// Possible reply values are: terminal, terminal16, terminal256, terminal16m.
//
// The value terminal is a monochrome terminal.
//
// The value terminal16 is a 4-bit, 16 color terminal.
//
// The value terminal256 is a 8-bit, 256 color terminal.
//
// The value terminal16m is a 24-bit, 16 million color terminal.
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
func UnderlineChar(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	if !utf8.ValidString(s) {
		return "", fmt.Errorf("underlinechar %q: %w", s, ErrRune)
	}
	if !color.Enable {
		return s, nil
	}
	b := &strings.Builder{}
	r, _ := utf8.DecodeRuneInString(s)
	t, err := template.New("underline").Parse("{{define \"TEXT\"}}\033[0m\033[4m{{.}}\033[0m{{end}}")
	if err != nil {
		return "", fmt.Errorf("underlinechar new template: %w", err)
	}
	if err := t.ExecuteTemplate(b, "TEXT", string(r)); err != nil {
		return "", fmt.Errorf("underlinechar execute template: %w", err)
	}
	return b.String(), nil
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
			if filepath.Ext(key) != ".min" {
				continue
			}
			s := strings.Split(keys[i], ".")
			if len(s) == 0 {
				continue
			}
			base := strings.Join(s[0:len(s)-1], ".")
			m, err := UnderlineChar("m")
			if err != nil {
				// must use standard log package
				log.Fatal("underline keys", keys, err)
			}
			keys[i] = fmt.Sprintf("%s.%sin", base, m)
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

// Alert returns the string "Problem:" using the error color.
func Alert() string {
	return color.Error.Sprint("Problem:") + "\n"
}

// Example returns the string using the debug color.
func Example(s string) string {
	return color.Debug.Sprint(s)
}

// Inform returns "Information:" using the info color.
func Inform() string {
	return color.Info.Sprint("Information:") + "\n"
}

// Bool returns a checkmark ✓ when true or a cross ✗ when false.
func Bool(b bool) string {
	const check, cross = "✓", "✗"
	if b {
		return color.Success.Sprint(check)
	}
	return color.Warn.Sprint(cross)
}

// Options writes the string and a sorted list of opts to the writer.
// If shorthand is true, the options are underlined.
// If flag is true, the string is prepended with "flag".
func Options(w io.Writer, s string, shorthand, flag bool, opts ...string) {
	if w == nil {
		w = io.Discard
	}
	if len(opts) == 0 {
		return
	}
	sort.Strings(opts)
	keys := strings.Join(opts, ", ")
	if shorthand {
		keys = UnderlineKeys(opts...)
	}
	if flag {
		fmt.Fprintln(w, s)
		fmt.Fprintf(w, "flag options: %s", color.Info.Sprint(keys))
		return
	}
	fmt.Fprintln(w, s+".")
	fmt.Fprintf(w, "  Options: %s", color.Info.Sprint(keys))
}

// Comment returns a string in the comment color.
func Comment(s string) string {
	return color.Comment.Sprint(s)
}

// Fuzzy returns a string in the fuzzy color.
func Fuzzy(s string) string {
	return color.OpFuzzy.Sprint(s)
}

// Info returns a string in the info color.
func Info(s string) string {
	return color.Info.Sprint(s)
}

// Secondary returns a string in the secondary color.
func Secondary(s string) string {
	return color.Secondary.Sprint(s)
}
