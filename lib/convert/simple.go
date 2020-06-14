package convert

import (
	"errors"
	"io/ioutil"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// DString decodes simple character encoding text.
func DString(s string, c charmap.Charmap) (result []byte, err error) {
	decoder := c.NewDecoder()
	reader := transform.NewReader(strings.NewReader(s), decoder)
	return ioutil.ReadAll(reader)
}

// EString encodes text into a simple character encoding.
func EString(s string, c charmap.Charmap) (result []byte, err error) {
	if !utf8.Valid([]byte(s)) {
		return result, errors.New("convert.estring: string is not encoded as utf8")
	}
	encoder := c.NewEncoder()
	reader := transform.NewReader(strings.NewReader(s), encoder)
	return ioutil.ReadAll(reader)
}

// D437 decodes IBM Code Page 437 encoded text.
func D437(s string) (result []byte, err error) {
	return DString(s, *charmap.CodePage437)
}

// E437 encodes text into IBM Code Page 437.
func E437(s string) (result []byte, err error) {
	return EString(s, *charmap.CodePage437)
}
