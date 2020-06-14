package convert

import (
	"encoding/base64"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

const (
	bom  = unicode.UseBOM
	_bom = unicode.IgnoreBOM
	be   = unicode.BigEndian
	le   = unicode.LittleEndian
)

// Base64Decode decodes a base64 string.
func Base64Decode(s string) (result []byte, err error) {
	return base64.StdEncoding.DecodeString(s)
}

// Base64Encode encodes a string to base64.
func Base64Encode(s string) (result string) {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// UTF8 transforms string to UTF8 encoding.
func UTF8(s string) (result []byte, length int, err error) {
	return transform.Bytes(unicode.UTF8.NewEncoder(), []byte(s))
}

// UTF16BE transforms string to UTF16 big-endian encoding.
func UTF16BE(s string) (result []byte, length int, err error) {
	return transform.Bytes(unicode.UTF16(be, bom).NewEncoder(), []byte(s))
}

// UTF16LE transforms string to UTF16 little-endian encoding.
func UTF16LE(s string) (result []byte, length int, err error) {
	return transform.Bytes(unicode.UTF16(le, bom).NewEncoder(), []byte(s))
}
