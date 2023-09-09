package str

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/gookit/color"
)

var ErrRune = errors.New("invalid encoded rune")

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
	// HBar is a the Unicode horizontal bar character.
	HBar = "\u2500"

	none = "none"
	term = "terminal"
)

func (t terminal) String() string {
	return [...]string{none, term, term, "terminal256", "terminal16m"}[t]
}

// Border wraps the string around a single line border.
func Border(s string) *bytes.Buffer {
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
	w := new(bytes.Buffer)
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
	return w
}

// Center align text to a the width of an area.
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

// Head returns a colored and underlined string for use as a header.
// Provide a fixed width value for the underline border or set to zero.
func Head(width int, s string) string {
	const div, padding = 2, 4
	var (
		h string
		p string
		r string
	)
	if width == 0 {
		r = strings.Repeat(HBar, len(s)+padding)
		p = strings.Repeat(" ", padding/div)
	} else {
		r = strings.Repeat(HBar, width)
		p = strings.Repeat(" ", (width-len(s))/div)
	}
	return fmt.Sprintf("\n%s%s%s\n%s", p, h, p, r)
}

func HeadDark(width int, s string) string {
	r := color.OpFuzzy.Sprint(strings.Repeat(HBar, width))
	h := color.Primary.Sprint(Center(width, s))
	return fmt.Sprintf("%s\n%s\n", r, h)
}

// HR returns a horizontal ruler or line break.
func HR(width int) string {
	return fmt.Sprintf(" %s", ColSec(strings.Repeat(HBar, width)))
}

func HRPad(width int) string {
	return fmt.Sprintf(" \n%s\n", ColSec(strings.Repeat(HBar, width)))
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
func UnderlineChar(c string) (string, error) {
	if c == "" {
		return "", nil
	}
	if !utf8.ValidString(c) {
		return "", fmt.Errorf("underlinechar %q: %w", c, ErrRune)
	}
	if !color.Enable {
		return c, nil
	}
	w := new(bytes.Buffer)
	r, _ := utf8.DecodeRuneInString(c)
	t, err := template.New("underline").Parse("{{define \"TEXT\"}}\033[0m\033[4m{{.}}\033[0m{{end}}")
	if err != nil {
		return "", fmt.Errorf("underlinechar new template: %w", err)
	}
	if err := t.ExecuteTemplate(w, "TEXT", string(r)); err != nil {
		return "", fmt.Errorf("underlinechar execute template: %w", err)
	}
	return w.String(), nil
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
