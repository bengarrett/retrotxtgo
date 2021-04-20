// Package convert is extends Go's x/text/encoding capability
// to convert legacy encoded text to a modern UTF-8 encoding.
package convert

import (
	"bytes"

	"golang.org/x/text/encoding"
)

// Convert 8-bit legacy or other Unicode text to UTF-8.
type Convert struct {
	Source Source // Source text for conversion.
	Output Output // Output UTF-8 text.
	Flags  Flags  // User supplied flag values.
}

// Source text for conversion.
type Source struct {
	B         []byte            // Source text as bytes.
	E         encoding.Encoding // Text encoding.
	table     bool              // flag Source.B as text for display as a codepage table.
	lineBreak [2]rune           // line break controls used by the text.
}

// Output UTF-8 text.
type Output struct {
	R          []rune // Output text as runes.
	ignores    []rune // runes to be ignored.
	len        int    // R (runes) count.
	lineBreaks bool   // use line break controls.
}

// Flags are the user supplied values.
type Flags struct {
	Controls  []string // Always use these control codes.
	SwapChars []int    // Swap out these characters with UTF-8 alternatives.
	Width     int      // Maximum text width per-line.
}

const (
	// EOF end of file character.
	EOF = 26
)

// BOM is the UTF-8 byte order mark prefix.
func BOM() []byte {
	return []byte{239, 187, 191} // 0xEF,0xBB,0xBF
}

// EndOfFile will cut text at the first DOS end-of-file marker.
func EndOfFile(b ...byte) []byte {
	if cut := bytes.IndexByte(b, EOF); cut > 0 {
		return b[:cut]
	}
	return b
}

// MakeBytes generates a 256 character or 8-bit container ready to hold legacy code point values.
func MakeBytes() (m []byte) {
	const max = 256
	m = make([]byte, max)
	for i := 0; i < max; i++ {
		m[i] = byte(i)
	}
	return m
}

// Mark adds a UTF-8 byte order mark to the text if it doesn't already exist.
func Mark(b ...byte) []byte {
	const min = 3
	if len(b) >= min {
		if t := b[:3]; bytes.Equal(t, BOM()) {
			return b
		}
	}
	return append(BOM(), b...)
}
