package nl

import "fmt"

type LineBreaks int

const (
	NL LineBreaks = iota
	Dos
	Win
	C64
	Darwin
	Mac
	Amiga
	Linux
	Unix
)

const (
	// Linefeed is a Linux/macOS line break.
	Linefeed rune = 10
	// CarriageReturn is a partial line break for Windows/DOS.
	CarriageReturn rune = 13
)

// LineBreak returns line break character for the system platform.
func LineBreak(platform LineBreaks) string {
	switch platform {
	case Dos, Win:
		return fmt.Sprintf("%x%x", CarriageReturn, Linefeed)
	case C64, Darwin, Mac:
		return fmt.Sprintf("%x", CarriageReturn)
	case Amiga, Linux, Unix:
		return fmt.Sprintf("%x", Linefeed)
	case NL: // use operating system default
		return "\n"
	}
	return "\n"
}
