package create

import (
	"fmt"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

// Font enum.
type Font uint

const (
	// Automatic uses AutoFont to suggest a font.
	Automatic Font = iota
	// Mona is a Japanese language font for ShiftJIS encoding.
	Mona
	// VGA is an all-purpose 8 pixel IBM/MS-DOS era VGA font.
	VGA
)

func (f Font) String() string {
	return [...]string{"automatic", "mona", "vga"}[f]
}

// File is the packed filename of the font.
func (f Font) File() string {
	files := [...]string{"ibm-vga8", "mona", "ibm-vga8"}
	return fmt.Sprintf("%s.woff2", files[f])
}

// AutoFont applies the automatic font-family setting to suggest a font based on the given encoding.
func AutoFont(e encoding.Encoding) Font {
	if e == japanese.ShiftJIS {
		return Mona
	}
	return VGA
}

// Family returns the named font.
func Family(name string) Font {
	const a, m, v = "a", "m", "v"
	switch strings.ToLower(name) {
	case Automatic.String(), a:
		return Automatic
	case Mona.String(), m:
		return Mona
	case VGA.String(), v:
		return VGA
	default:
		return Automatic
	}
}

// Fonts are values for the CSS font-family attribute.
func Fonts() []string {
	return []string{Automatic.String(), Mona.String(), VGA.String()}
}
