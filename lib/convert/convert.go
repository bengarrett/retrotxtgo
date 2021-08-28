// Package convert is extends Go's x/text/encoding capability
// to convert legacy encoded text to a modern UTF-8 encoding.
package convert

import (
	"bytes"

	"golang.org/x/text/encoding"
)

// Convert 8-bit legacy or other Unicode text to UTF-8.
type Convert struct {
	Flags      Flag   // Commandline supplied flag values.
	Input      In     // Input text for conversion.
	Output     []rune // Output text as UTF-8 runes.
	ignores    []rune // runes to ignore.
	lineBreaks bool   // use line break controls?
}

// In is the text input for conversion.
type In struct {
	Encoding  encoding.Encoding // Bytes text encoding.
	Bytes     []byte            // Input text as bytes.
	lineBreak [2]rune           // line break controls used by the text.
	table     bool              // flag this text as a codepage table.
}

// Flag are the user supplied values.
type Flag struct {
	Controls  []string // Always use these control codes.
	SwapChars []int    // Swap out these characters with UTF-8 alternatives.
	MaxWidth  int      // Maximum text width per-line.
}

const (
	// EOF end of file character.
	EOF = 26
)

// BOM is the UTF-8 byte order mark prefix.
func BOM() []byte {
	const ef, bb, bf = 239, 187, 191
	return []byte{ef, bb, bf}
}

// EndOfFile will cut text at the first DOS end-of-file marker.
func EndOfFile(b []byte) []byte {
	if cut := bytes.IndexByte(b, EOF); cut > 0 {
		return b[:cut]
	}
	return b
}

// MakeBytes generates a 256 character or 8-bit container ready to hold legacy code point values.
func MakeBytes() []byte {
	const max = 256
	m := make([]byte, max)
	for i := 0; i < max; i++ {
		m[i] = byte(i)
	}
	return m
}

// Mark adds a UTF-8 byte order mark to the text if it doesn't already exist.
func Mark(b []byte) []byte {
	const min = 3
	if len(b) >= min {
		if t := b[:3]; bytes.Equal(t, BOM()) {
			return b
		}
	}
	return append(BOM(), b...)
}
