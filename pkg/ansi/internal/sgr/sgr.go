package sgr

import (
	"math"
	"strconv"
)

// Ps is the value of control sequence parameter. That is
// often referred to as Ps in VT100 and ANSI documentation.
type Ps int

const (
	Invalid             Ps = iota - 1  // A placeholder for unknown Ps values.
	Normal                             // Reset or default attributes.
	Bold                               // Increased font intensity.
	Faint                              // Decreased font intensity.
	Italic                             // Use an italic type font.
	Underline                          // Apply an underline text decoration.
	Blink                              // Animate the text with a standard speed blink.
	BlinkFast                          // Animate the text with a rapid blink.
	Inverse                            // Invert the color of the text.
	Conceal                            // Sets the text color to match the background color.
	StrikeThrough                      // Apply a horizontal line decoration through the middle of the text.
	Font0                              // Alternative font 1 (Eagle Sprint CGA). NOTE: invalid position
	Font1                              // Alternative font 2 (IBM BIOS).
	Font2                              // Alternative font 3 (IBM CGA).
	Font3                              // Alternative font 4 (IBM CGA Thin).
	Font4                              // Alternative font 5 (Amiga Topaz).
	Font5                              // Alternative font 6 (IBM EGA 8px).
	Font6                              // Alternative font 7 (IBM EGA 9px).
	Font7                              // Alternative font 8 (IBM VGA 8px).
	Font8                              // Alternative font 9 (IBM VGA 9px).
	Font9                              // Alternative font 10 (IBM MDA).
	Fraktur                            // Fraktur or Gothic font.
	Underline2x                        // Apply a doubled underline text decoration.
	NotBoldFaint                       // Unset Bold and Faint.
	NotItalicFraktur                   // Unset Italic and Fraktur.
	NotUnderline                       // Unset Underline and Underline2x.
	Steady                             // Unset Blink and BlinkFast.
	PositiveImg         Ps = iota      // Unset Inverse.
	Revealed                           // Unset Conceal.
	NotStrikeThrough                   // Unset StrikeThrough.
	Black                              // Text color black.
	Red                                // Text color red.
	Green                              // Text color green.
	Yellow                             // Text color yellow.
	Blue                               // Text color blue.
	Magenta                            // Text color magenta.
	Cyan                               // Text color cyan.
	White                              // Text color white.
	Extension                          //
	BlackB              Ps = iota + 1  // Background color black.
	RedB                               // Background color red.
	GreenB                             // Background color green.
	YellowB                            // Background color yellow.
	BlueB                              // Background color blue.
	MagentaB                           // Background color magenta.
	CyanB                              // Background color cyan.
	WhiteB                             // Background color white.
	ExtensionB                         //
	RevertB                            // Revert background color to default.
	Framed              Ps = iota + 2  // Framed text decoration.
	Encircled                          // Encircled text decoration.
	Overlined                          // Overlined text decoration.
	NotFramedEnbcircled                // Unset Framed and Encircled.
	NotOverlined                       // Unset overlined.
	BoldBlack           Ps = iota + 36 // Bold black text.
	BoldRed                            // Bold red text.
	BoldGreen                          // Bold green text.
	BoldYellow                         // Bold yellow text.
	BoldBlue                           // Bold blue text.
	BoldMagenta                        // Bold mageta text.
	BoldCyan                           // Bold cyan text.
	BoldWhite                          // Bold white text.
	BrightBlack         Ps = iota + 38 // Bright background black.
	BrightRed                          // Bright background red.
	BrightGreen                        // Bright background green.
	BrightYellow                       // Bright background yellow.
	BrightBlue                         // Bright background blue.
	BrightMagenta                      // Bright background magenta.
	BrightCyan                         // Bright background cyan.
	BrightWhite                        // Bright background white.
)

// String returns the CSS class name of p.
func (p Ps) String() string {
	if !p.Valid() {
		return ""
	}
	return "SGR" + strconv.Itoa(int(p))
}

// Valid returns true if p is a known control sequence parameter.
func (p Ps) Valid() bool {
	if p == Invalid {
		return false
	}
	switch {
	case
		p >= Normal && p <= Steady,
		p >= PositiveImg && p <= RevertB,
		p >= Framed && p <= NotOverlined,
		p >= BoldBlack && p <= BoldWhite,
		p >= BrightBlack && p <= BrightWhite:
		return true
	}
	return false
}

// Uint8 returns true if i is a valid 8 bit integer.
// That is a number between 0 and 255.
func Uint8(i int) bool {
	switch {
	case math.MaxUint8 < i, i < 0:
		return false
	}
	return true
}
