//Package encoding to handle the opening and reading of text files
package encoding

//https://stackoverflow.com/questions/36111777/golang-how-to-read-a-text-file

import (
	"bytes"
	"golang.org/x/net/html/charset"
)

// BOM is the UTF-8 byte order mark prefix.
var BOM = func() []byte {
	return []byte{239, 187, 191} // 0xEF,0xBB,0xBF
}

// ToBOM adds a UTF-8 byte order mark if it doesn't already exist.
func ToBOM(b []byte) []byte {
	if len(b) > 2 {
		if t := b[:3]; bytes.Equal(t, BOM()) == true {
			return b
		}
	}
	return append(BOM(), b...)
}

// UTF8 determines if a document is encoded as UTF-8.
func UTF8(b []byte) bool {
	_, name, _ := charset.DetermineEncoding(b, "text/plain")
	if name == "utf-8" {
		return true
	}
	return false
}
