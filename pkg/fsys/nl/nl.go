// Package nl provides line break characters for multiple system and microcomputer platforms.
package nl

import "fmt"

// LineBreaks is a system or microcomputer platform.
type LineBreaks int

const (
	NL     LineBreaks = iota // NL is the operating system default.
	Dos                      // Dos is an Microsoft DOS line break.
	Win                      // Win is a Windows line break.
	C64                      // C64 is a Commodore 64 line break.
	Darwin                   // Darwin is a macOS line break.
	Mac                      // Mac is an Apple Macintosh line break.
	Amiga                    // Amiga is a Commodore Amiga line break.
	Linux                    // Linux is a Linux line break.
	Unix                     // Unix is a Unix line break.
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
