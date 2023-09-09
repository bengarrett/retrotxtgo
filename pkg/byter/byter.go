package byter

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

var (
	ErrCharmap = fmt.Errorf("the c charmap cannot be nil")
	ErrUTF8    = errors.New("string cannot encode to utf-8")
)

const (
	SUB       = 26   // SUB is the ASCII control code for substitute.
	DosSUB    = 8594 // DosSub is the Unicode for the right-arrow.
	SymbolSUB = 9242 // SymbolSUB is the Unicode for the substitute character.
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
	s := string(b)
	if cut := strings.IndexRune(s, SymbolSUB); cut > 0 {
		return []byte(s[:cut])
	}
	// UTF-8 right-arrow which is displayed for the CP-437 substitute character code point 26
	if cut := strings.IndexRune(s, DosSUB); cut > 0 {
		return []byte(s[:cut])
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

// HexDecode decodes a hexadecimal string into bytes.
func HexDecode(s string) ([]byte, error) {
	src := []byte(s)
	dst := make([]byte, hex.DecodedLen(len(src)))
	if _, err := hex.Decode(dst, src); err != nil {
		return nil, fmt.Errorf("could not decode hexadecimal string: %q: %w", s, err)
	}
	return dst, nil
}

// HexEncode encodes a string into hexadecimal bytes.
func HexEncode(s string) []byte {
	src := []byte(s)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

// Decode a string using the character map.
func Decode(c *charmap.Charmap, s string) ([]byte, error) {
	if c == nil {
		return nil, ErrCharmap
	}
	decoder := c.NewDecoder()
	reader := transform.NewReader(strings.NewReader(s), decoder)
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("dstring io.readall error: %w", err)
	}
	return b, nil
}

// Encode the string using the character map.
func Encode(c *charmap.Charmap, s string) ([]byte, error) {
	if c == nil {
		return nil, ErrCharmap
	}
	b := []byte(s)
	if !utf8.Valid(b) {
		return nil, fmt.Errorf("estring: %w", ErrUTF8)
	}
	encoder := c.NewEncoder()
	reader := transform.NewReader(strings.NewReader(s), encoder)
	p, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("estring io.readall error: %w", err)
	}
	return p, nil
}
