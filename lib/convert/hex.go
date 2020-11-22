package convert

import (
	"encoding/hex"
	"fmt"
)

// HexDecode decodes hexadecimal into string bytes.
func HexDecode(hexadecimal string) (result []byte, err error) {
	src := []byte(hexadecimal)
	result = make([]byte, hex.DecodedLen(len(src)))
	if _, err = hex.Decode(result, src); err != nil {
		return nil, fmt.Errorf("could not decode hexadecimal string: %q: %w", hexadecimal, err)
	}
	return result, err
}

// HexEncode encodes string into hexadecimal encoded bytes.
func HexEncode(s string) (result []byte) {
	src := []byte(s)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}
