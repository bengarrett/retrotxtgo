package layout

import "github.com/bengarrett/retrotxtgo/lib/logs"

// Layout are HTML template variations.
type Layout int

const (
	// use 0 as an error placeholder.
	_ Layout = iota
	// Standard template with external CSS, JS, fonts.
	Standard
	// Inline template with CSS and JS embedded.
	Inline
	// Compact template with external CSS, JS, fonts and no meta-tags.
	Compact
	// None template, just print the generated HTML.
	None
)

const (
	none     = "none"
	compact  = "compact"
	inline   = "inline"
	standard = "standard"
	unknown  = "unknown"

	ZipName = "retrotxt.zip"
)

// ParseLayout parses possible --layout argument values.
func ParseLayout(name string) (Layout, error) {
	switch name {
	case standard, "s":
		return Standard, nil
	case inline, "i":
		return Inline, nil
	case compact, "c":
		return Compact, nil
	case none, "n":
		return None, nil
	}
	return 0, logs.ErrTmplName
}

// Pack is the packed name of the HTML template.
func (l Layout) Pack() string {
	return [...]string{unknown, standard, standard, standard, none}[l]
}

func (l Layout) String() string {
	return [...]string{unknown, standard, inline, compact, none}[l]
}

func UseCSS(l Layout) bool {
	switch l {
	case Standard, Compact:
		return true
	case Inline, None:
		return false
	}
	return false
}

func UseFontCSS(l Layout) bool {
	switch l {
	case Standard, Compact:
		return true
	case Inline, None:
		return false
	}
	return false
}

func UseIcon(l Layout) bool {
	switch l {
	case Standard, Compact:
		return true
	case Inline, None:
		return false
	}
	return false
}

func UseJS(l Layout) bool {
	return false
}
