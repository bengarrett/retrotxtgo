package convert

import (
	"encoding/hex"
	"fmt"
)

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
