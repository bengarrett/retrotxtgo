package convert

import "encoding/hex"

// HexDecode decodes hexadecimal into string bytes.
func HexDecode(hexadecimal string) (result []byte, err error) {
	src := []byte(hexadecimal)
	result = make([]byte, hex.DecodedLen(len(src)))
	_, err = hex.Decode(result, src)
	return result, err
}

// HexEncode encodes string into hexadecimal encoded bytes.
func HexEncode(s string) (result []byte) {
	src := []byte(s)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}
