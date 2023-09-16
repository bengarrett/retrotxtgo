// Package nl provides line break characters for multiple system and microcomputer platforms.
package nl

// System is a micro or computer platform.
type System int

const (
	NL        System = iota // NL is the host operating system default.
	Amiga                   // Amiga is the Commodore Amiga.
	Commodore               // Commodore is the Commodore 64 and compatibles.
	Darwin                  // Darwin is Apple macOS.
	Linux                   // Linux is the Linux operating system.
	Macintosh               // Macintosh is the Apple Macintosh.
	PCDos                   // PCDos is an IBM PC and Microsoft DOS and compatibles.
	Unix                    // Unix are both Unix and BSD.
	Windows                 // Windows is Microsoft Windows.
)

const (
	LF rune = 10 // LF is the control code for a line feed.
	CR rune = 13 // CR is the control code for a carriage return.
)

// NewLine returns a new line or line break characters for the system platform.
func NewLine(s System) string {
	switch s {
	case PCDos, Windows:
		return string(CR) + string(LF)
	case Commodore, Darwin, Macintosh:
		return string(CR)
	case Amiga, Linux, Unix:
		return string(LF)
	case NL:
		return "\n"
	}
	return ""
}
