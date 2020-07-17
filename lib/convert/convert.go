//Package convert is extends Go's x/text/encoding capability to convert legacy text
// to UTF-8.
package convert

import (
	"bytes"
)

// BOM is the UTF-8 byte order mark prefix.
func BOM() []byte {
	return []byte{239, 187, 191} // 0xEF,0xBB,0xBF
}

// EndOfFile will cut text at the first DOS end-of-file marker.
func EndOfFile(b ...byte) []byte {
	if cut := bytes.IndexByte(b, 26); cut > 0 {
		return b[:cut]
	}
	return b
}

// MakeBytes generates a 256 character or 8-bit container ready to hold legacy code point values.
func MakeBytes() (m []byte) {
	m = make([]byte, 256)
	for i := 0; i <= 255; i++ {
		m[i] = uint8(byte(i))
	}
	return m
}

// Mark adds a UTF-8 byte order mark to the text if it doesn't already exist.
func Mark(b ...byte) []byte {
	if len(b) > 2 {
		if t := b[:3]; bytes.Equal(t, BOM()) {
			return b
		}
	}
	return append(BOM(), b...)
}
