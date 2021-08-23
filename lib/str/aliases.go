// Package str manipulates strings and standard output text.
package str

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gookit/color"
)

type verbs uint

const (
	confirm verbs = iota
	info
	problem
)

func sprint(v verbs) string {
	s, c := "", ""
	switch v {
	case confirm:
		s = "Confirm:"
		c = color.Question.Sprint(s)
	case info:
		s = "Information:"
		c = color.Info.Sprint(s)
	case problem:
		s = "Problem:"
		c = color.Error.Sprint(s)
	}
	if c == "" {
		return fmt.Sprintf("%s\n", s)
	}
	return fmt.Sprintf("%s\n", c)
}

// Alert prints "Problem" using the Error color.
func Alert() string {
	return sprint(problem)
}

// Confirm prints the "Confirm" string using the Question color.
func Confirm() string {
	return sprint(confirm)
}

// Example prints the string using the Debug color.
func Example(s string) string {
	return color.Debug.Sprint(s)
}

// Info prints "Information" using the Info color.
func Info() string {
	return sprint(info)
}

func Path(name string) string {
	return color.Secondary.Sprint(name)
}

// Bool returns a checkmark ✓ when true or a cross ✗ when false.
func Bool(b bool) string {
	const check, cross = "✓", "✗"
	if b {
		return color.Success.Sprint(check)
	}
	return color.Warn.Sprint(cross)
}

// Default appends (default ...) to the usage string.
func Default(s, def string) string {
	return fmt.Sprintf("%s (default \"%s\")", s, def)
}

// Options appends options: ... to the usage string.
func Options(s string, shorthand, flagHelp bool, opts ...string) (usage string) {
	var keys string
	if len(opts) == 0 {
		return s
	}
	sort.Strings(opts)
	if shorthand {
		keys = UnderlineKeys(opts...)
	} else {
		keys = strings.Join(opts, ", ")
	}
	if flagHelp {
		return fmt.Sprintf("%s\nflag options: %s", s, color.Info.Sprint(keys))
	}
	return fmt.Sprintf("%s.\n  Options: %s", s, color.Info.Sprint(keys))
}

// Required appends (required) to the string.
func Required(s string) string {
	return fmt.Sprintf("%s (required)", color.Primary.Sprint(s))
}

// color aliases.

// Cb returns a string in the color named secondary.
func Cb(t string) string {
	return color.Secondary.Sprint(t)
}

// Cc returns a string in the color named comment.
func Cc(t string) string {
	return color.Comment.Sprint(t)
}

// Ce returns a string in the color named warn.
func Ce(t string) string {
	return color.Warn.Sprint(t)
}

// Cf returns a string in the style named fuzzy.
func Cf(t string) string {
	return color.OpFuzzy.Sprint(t)
}

// Ci returns a string in the style named italic.
func Ci(t string) string {
	return color.OpItalic.Sprint(t)
}

// Cinf returns a string in the color named info.
func Cinf(t string) string {
	return color.Info.Sprint(t)
}

// Cp returns a string in the color named primary.
func Cp(t string) string {
	return color.Primary.Sprint(t)
}

// Cs returns a string in the color named success.
func Cs(t string) string {
	return color.Success.Sprint(t)
}
