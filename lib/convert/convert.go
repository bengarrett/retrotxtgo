// Package convert is extends Go's x/text/encoding capability
// to convert legacy encoded text to a modern UTF-8 encoding.
package convert

import (
	"bytes"
	"strings"

	"golang.org/x/text/encoding"
)

// Convert 8-bit legacy or other Unicode text to UTF-8.
type Convert struct {
	Flags      Flag   // Commandline supplied flag values.
	Input      In     // Input text for transformation.
	Output     []rune // Transformed UTF-8 runes.
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
	SwapChars []string // Swap out these characters with UTF-8 alternatives.
	MaxWidth  int      // Maximum text width per-line.
}

const (
	DosSUB    = 8594
	SymbolSUB = 9242
)

// BOM is the UTF-8 byte order mark prefix.
func BOM() []byte {
	const ef, bb, bf = 239, 187, 191
	return []byte{ef, bb, bf}
}

// TrimEOF will cut text at the first occurrence of the SUB character.
// The SUB is used by DOS and CP/M as an end-of-file marker.
func TrimEOF(b []byte) []byte {
	// ASCII control code
	if cut := bytes.IndexByte(b, SUB); cut > 0 {
		return b[:cut]
	}
	// UTF-8 symbol for substitute character
	if cut := strings.IndexRune(string(b), SymbolSUB); cut > 0 {
		return []byte(string(b)[:cut])
	}
	// UTF-8 right-arrow which is displayed for the CP-437 substitute character code point 26
	if cut := strings.IndexRune(string(b), DosSUB); cut > 0 {
		return []byte(string(b)[:cut])
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
