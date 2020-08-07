package str

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gookit/color"
)

// Alert prints "problem:" in Error color.
func Alert() string {
	return color.Error.Sprint("problem:")
}

// Info prints "info:" in Info color.
func Info() string {
	return color.Info.Sprint("info:")
}

// color aliases.

// Cb secondary.
func Cb(t string) string {
	return color.Secondary.Sprint(t)
}

// Cc comment.
func Cc(t string) string {
	return color.Comment.Sprint(t)
}

// Ce warn.
func Ce(t string) string {
	return color.Warn.Sprint(t)
}

// Cf fuzzy.
func Cf(t string) string {
	return color.OpFuzzy.Sprint(t)
}

// Ci italic.
func Ci(t string) string {
	return color.OpItalic.Sprint(t)
}

// Cinf info.
func Cinf(t string) string {
	return color.Info.Sprint(t)
}

// Cp primary.
func Cp(t string) string {
	return color.Primary.Sprint(t)
}

// Cs success.
func Cs(t string) string {
	return color.Success.Sprint(t)
}

// Bool returns a ✓ or ✗.
func Bool(b bool) string {
	switch b {
	case true:
		return color.Success.Sprint("✓")
	default:
		return color.Warn.Sprint("✗")
	}
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
	return fmt.Sprintf("%s\noptions: %s", s, color.Info.Sprint(keys))
}

// Required appends (required) to the usage string.
func Required(s string) (usage string) {
	return fmt.Sprintf("%s (required)", color.Primary.Sprint(s))
}
