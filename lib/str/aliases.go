// Package str manipulates strings and standard output text.
package str

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gookit/color"
)

// Alert prints "problem:" in Error color.
func Alert() string {
	s := "problem:"
	e := color.Error.Sprint(s)
	if e == "" {
		return s
	}
	return e
}

// Info prints "info:" in Info color.
func Info() string {
	s := "info:"
	e := color.Info.Sprint(s)
	if e == "" {
		return s
	}
	return e
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

// Example is intended for the cobra.Command Example fields.
func Example(s string) string {
	return color.Info.Sprint(s)
}

// Options appends options: ... to the usage string.
func Options(s string, shorthand bool, opts ...string) (usage string) {
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
	return fmt.Sprintf("%s\n options: %s", s, color.Info.Sprint(keys))
}

// Required appends (required) to the usage string.
func Required(s string) (usage string) {
	return fmt.Sprintf("%s (required)", color.Primary.Sprint(s))
}
