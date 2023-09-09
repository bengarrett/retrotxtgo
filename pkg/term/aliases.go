// Package term manipulates strings and standard output text.
package term

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gookit/color"
)

type verbs uint

const (
	info verbs = iota
	problem
)

func (v verbs) String() string {
	s, c := "", ""
	switch v {
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
	return problem.String()
}

// Example prints the string using the Debug color.
func Example(s string) string {
	return color.Debug.Sprint(s)
}

// Info prints "Information" using the Info color.
func Info() string {
	return info.String()
}

// Bool returns a checkmark ✓ when true or a cross ✗ when false.
func Bool(b bool) string {
	const check, cross = "✓", "✗"
	if b {
		return color.Success.Sprint(check)
	}
	return color.Warn.Sprint(cross)
}

// Options appends options: ... to the usage string.
func Options(s string, shorthand, flagHelp bool, opts ...string) string {
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

// color aliases.

// ColCmt returns a string in the comment color.
func ColCmt(s string) string {
	return color.Comment.Sprint(s)
}

// ColFuz returns a string in the fuzzy color.
func ColFuz(s string) string {
	return color.OpFuzzy.Sprint(s)
}

// ColInf returns a string in the info color.
func ColInf(s string) string {
	return color.Info.Sprint(s)
}

// ColSec returns a string in the secondary color.
func ColSec(s string) string {
	return color.Secondary.Sprint(s)
}
