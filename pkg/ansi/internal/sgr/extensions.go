package sgr

import "github.com/bengarrett/retrotxtgo/lib/ansi/internal/xterm"

// Extensions are custom color parameters that were not defined
// in the original VT100 and ANSI documentations. These extend
// the color support beyond the 16 named colors.
type Extensions struct {
	BG    bool // Background toggle.
	FG    bool // Foreground text toggle.
	Xterm bool // 256 color toggle.
	Color int8 // 256 color value.
	RGB   bool // RGB color toggle.
	R     int8 // Red color value.
	G     int8 // Green color value.
	B     int8 // Blue color value.
}

func (e *Extensions) Reset() {
	e.BG = false
	e.Xterm = false
	e.Color = -1
	e.RGB = false
	e.R = -1
	e.G = -1
	e.B = -1
}

func (e *Extensions) Scan(s int) (cont bool) {
	switch Ps(s) {
	case Extension:
		e.SetFG()
		return true
	case ExtensionB:
		e.SetBG()
		return true
	}
	switch {
	// case e.FG && e.RGB:
	// 	//
	// 	return true
	case (e.FG || e.BG) && e.Xterm:
		if !Uint8(s) {
			// bad xterm value
			e.Reset()
		}
		e.Color = int8(s)
		return true
	// case e.BG && Ps(s) == 2:
	// 	e.SetRGB()
	// 	return true
	case e.BG && Ps(s) == xterm.Control:
		e.SetXterm()
		return true
	case e.BG:
		e.Reset()
		return true
	}
	return false
}

func (e *Extensions) SetBG() {
	e.BG = true
}

func (e *Extensions) SetFG() {
	e.FG = true
}

func (e *Extensions) SetXterm() {
	e.Xterm = true
}

func (e *Extensions) SetRGB() {
	e.RGB = true
}

func (e *Extensions) SetRed(i uint8) {
	e.R = int8(i)
}

func (e *Extensions) SetGreen(i uint8) {
	e.G = int8(i)
}

func (e *Extensions) SetBlue(i uint8) {
	e.B = int8(i)
}
