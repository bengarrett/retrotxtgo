// Package nl provides line break characters for multiple system and microcomputer platforms.
package nl

// System is a system or microcomputer platform.
type System int

const (
	NL     System = iota // NL is the operating system default.
	Dos                  // Dos is an Microsoft DOS line break.
	Win                  // Win is a Windows line break.
	C64                  // C64 is a Commodore 64 line break.
	Darwin               // Darwin is a macOS line break.
	Mac                  // Mac is an Apple Macintosh line break.
	Amiga                // Amiga is a Commodore Amiga line break.
	Linux                // Linux is a Linux line break.
	Unix                 // Unix is a Unix line break.
)

const (
	Linefeed       rune = 10 // Linefeed is a Linux/macOS line break.
	CarriageReturn rune = 13 // CarriageReturn is a partial line break for Windows/DOS.
)

// NewLine returns a new line or line break characters for the system platform.
func NewLine(s System) string {
	switch s {
	case Dos, Win:
		return string(CarriageReturn) + string(Linefeed)
	case C64, Darwin, Mac:
		return string(CarriageReturn)
	case Amiga, Linux, Unix:
		return string(Linefeed)
	case NL: // use operating system default
		return "\n"
	}
	return "\n"
}
