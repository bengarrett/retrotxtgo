package datastream

import (
	"bytes"
	"strconv"
	"unicode"

	"github.com/bengarrett/retrotxtgo/pkg/ansi/internal/sgr"
)

// Erase parameter number; either 0, 1, 2 or -1 for not found.
type Erase int

const (
	NoErase                Erase = iota - 1 // Erase sequnce not found.
	EraseFromCursorToEnd                    // Erase from the cursor to end inclusive (default)
	EraseFromStartToCursor                  // Erase from start to cursor inclusive.
	EraseAll                                // Erase all.
)

// Number determines if b is comprised of decimal digits.
func Number(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	rs := bytes.Runes(b)
	for _, r := range rs {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// Cursor looks in b for the c character with a preceding uint8 number.
// If no match is found then 0 is returned.
// If the preceding number is greater than 255 then 0 is returned.
func cursor(b []byte, c byte) uint8 {
	i := bytes.IndexByte(b, c)
	if i == -1 {
		return 0
	}
	v := b[0:i]
	if !Number(v) {
		return 0
	}
	s, err := strconv.Atoi(string(v))
	if err != nil {
		return 0
	}
	if s > 255 {
		return 0
	}
	return uint8(s)
}

// CursorPos looks in b for a line and/or column values in b.
// The two values must be separated with a semicolumn.
// If b is invalid then both line and col return 0.
// If b equals c, then Home is requested and line and col will return 1.
func cursorPos(b []byte, c byte) (line uint8, col uint8) {
	if len(b) == 0 {
		return 0, 0
	}
	if len(b) == 1 && b[0] == c {
		return 1, 1
	}
	if l := cursor(b, c); l > 0 {
		return l, 1
	}
	if b[0] == byte(';') {
		if col = cursor(b[1:], c); col > 0 {
			return 1, col
		}
	}
	s := bytes.Split(b, []byte(";"))
	if len(s) != 2 {
		return 0, 0
	}
	if !Number(s[0]) {
		return 0, 0
	}
	r, err := strconv.Atoi(string(s[0]))
	if err != nil {
		return 0, 0
	}
	if r > 255 {
		return 0, 0
	}
	col = cursor(s[1], c)
	if col == 0 {
		return 0, 0
	}
	return uint8(r), col
}

// Erase looks in b for the c character with a preceding number of 0, 1 or 2.
// If no match is found then -1 is returned.
// If the preceding number is greater than 2 then -1 is returned.
func erase(b []byte, c byte) Erase {
	i := bytes.IndexByte(b, c)
	if i == -1 {
		return NoErase
	}
	v := b[0:i]
	if !Number(v) {
		return NoErase
	}
	s, err := strconv.Atoi(string(v))
	if err != nil {
		return NoErase
	}
	if s > 2 {
		return NoErase
	}
	return Erase(s)
}

// CUU Cursor Up. Move the cursor up by a number of lines.
// The cursor can never be outside of the page boundary.
func CUU(b []byte) (up uint8) {
	return cursor(b, byte('A'))
}

// CUD Cursor Down. Move the cursor down by a number of lines.
func CUD(b []byte) (down uint8) {
	return cursor(b, byte('B'))
}

// CUF Cursor Forward. Move the cursor right by a number of columns.
func CUF(b []byte) (right uint8) {
	return cursor(b, byte('C'))
}

// CUB Cursor Backward. Move the cursor left by a number of columns.
func CUB(b []byte) (left uint8) {
	return cursor(b, byte('D'))
}

// HPR Horizontal Position Relative.
// Same as Cursor Forward, except limited to the active line.
func HPR(b []byte) (right uint8) {
	return cursor(b, byte('a'))
}

// VPR Vertical Position Relative.
// Same as Cursor Down. Move the cursor down by a number of lines.
func VPR(b []byte) (down uint8) {
	return cursor(b, byte('e'))
}

// CNL Cursor Next Line. Moves the cursor down by a number of lines
// and to the beginning of the line.
func CNL(b []byte) (down uint8) {
	return cursor(b, byte('E'))
}

// CNL Cursor Preceding Line. Moves the cursor up by a number of lines
// and to the beginning of the line.
func CPL(b []byte) (up uint8) {
	return cursor(b, byte('F'))
}

// HPA Horizontal Position Absolute.
// Positions the cursor to the column (in the same line).
func HPA(b []byte) (col uint8) {
	return cursor(b, byte('`'))
}

// VPA Vertical Position Absolute.
// Positions the cursor to the line (in the same column).
func VPA(b []byte) (line uint8) {
	return cursor(b, byte('d'))
}

// CHA Cursor Horizontal Absolute.
// Positions the cursor to the column (in the same line).
func CHA(b []byte) (col uint8) {
	return cursor(b, byte('G'))
}

// CUP Cursor Position. (ESC [ Pn 1; Pn 2 H).
// Positions the cursor to line Pn1 and column Pn2.
func CUP(b []byte) (line uint8, col uint8) {
	return cursorPos(b, byte('H'))
}

// HVP, Horizontal & Vertical Position. (ESC [ Pn 1; Pn 2 f).
// Positions the cursor to line Pn1 and column Pn2.
func HVP(b []byte) (line uint8, col uint8) {
	return cursorPos(b, byte('f'))
}

// DCH Delete Character.
// Returns the number of characters to delete starting with the
// character at the cursor.
func DCH(b []byte) (chars uint8) {
	return cursor(b, byte('P'))
}

// DL Delete Line.
// Returns the number of lines to delete starting with the
// active line at the cursor.
func DL(b []byte) (lines uint8) {
	return cursor(b, byte('M'))
}

// ED Erase in Display.
// Erases some or all of the characters in the page as selected by Erase.
// An invalid value will return -1.
func ED(b []byte) Erase {
	return erase(b, byte('J'))
}

// EL Erase in Line.
// Erases some or all of the characters in the active line as selected by Ps.
// An invalid value will return -1.
func EL(b []byte) Erase {
	return erase(b, byte('K'))
}

// ICH Insert Character.
// Returns the number of space characters to insert at the cursor.
func ICH(b []byte) (spaces uint8) {
	return cursor(b, byte('@'))
}

// IL Insert Line.
// Returns the number of lines to insert at the cursor.
// The cursor position does not change.
func IL(b []byte) (lines uint8) {
	return cursor(b, byte('@'))
}

// REP Repeat.
// Causes the character immediately preceding the control to be repeated.
func REP(b []byte) (repeat uint8) {
	return cursor(b, byte('b'))
}

// SGR Select Graphic Rendition.
// Set a register that determines the graphic rendition with which characters
// subsequently entered are displayed as selected by each Ps.
func SGR(b []byte) sgr.Attributes {
	return sgr.DataStream(b)
}
