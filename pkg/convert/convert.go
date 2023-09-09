// Package convert is extends Go's x/text/encoding capability
// to convert legacy encoded text to a modern UTF-8 encoding.
package convert

import (
	"golang.org/x/text/encoding"
)

// Convert legacy 8-bit codepage encode or Unicode byte array text to UTF-8 runes.
type Convert struct {
	Flags Flag // Flags are the cmd supplied flag values.
	Input struct {
		Encoding  encoding.Encoding // Encoding are the encoding of the input text.
		Bytes     []byte            // Bytes are the input text as bytes.
		LineBreak [2]rune           // Line break controls used by the text.
		Table     bool              // Table flags this text as a codepage table.
	}
	Output     []rune // Output are the transformed UTF-8 runes.
	Ignores    []rune // Ignores these runes.
	LineBreaks bool   // LineBreaks uses line break controls.
}

// Flag are the user supplied values.
type Flag struct {
	Controls  []string // Always use these control codes.
	SwapChars []string // Swap out these characters with UTF-8 alternatives.
	MaxWidth  int      // Maximum text width per-line.
}
