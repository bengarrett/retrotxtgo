package logs

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gookit/color"
)

// Alert prints "problem:" in Error color.
var Alert = func() string {
	return color.Error.Sprint("problem:")
}

// Info prints "info:" in Info color.
var Info = func() string {
	return color.Info.Sprint("info:")
}

// Example is intended for the cobra.Command Example fields.
var Example = func(t string) string {
	return color.Info.Sprint(t)
}

// color aliases
var (
	Cb = func(t string) string {
		return color.Secondary.Sprint(t)
	}
	Cc = func(t string) string {
		return color.Comment.Sprint(t)
	}
	Ce = func(t string) string {
		return color.Warn.Sprint(t)
	}
	Cf = func(t string) string {
		return color.OpFuzzy.Sprint(t)
	}
	Ci = func(t string) string {
		return color.OpItalic.Sprint(t)
	}
	Cinf = func(t string) string {
		return color.Info.Sprint(t)
	}
	Cp = func(t string) string {
		return color.Primary.Sprint(t)
	}
	Cs = func(t string) string {
		return color.Success.Sprint(t)
	}
)

// Options appends options: ... to the usage string.
func Options(s string, opts []string, shorthand bool) (usage string) {
	var keys string
	if len(opts) == 0 {
		return s
	}
	sort.Strings(opts)
	if shorthand {
		keys = UnderlineKeys(opts)
	} else {
		keys = strings.Join(opts, ", ")
	}
	return fmt.Sprintf("%s\noptions: %s", s, color.Info.Sprint(keys))
}

// Required appends (required) to the usage string.
func Required(s string) (usage string) {
	return fmt.Sprintf("%s (required)", color.Primary.Sprint(s))
}
