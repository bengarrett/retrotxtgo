package convert

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// DString decodes simple character encoding text.
func DString(s string, c *charmap.Charmap) (result []byte, err error) {
	decoder := c.NewDecoder()
	reader := transform.NewReader(strings.NewReader(s), decoder)
	result, err = ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("dstring ioutil readall error: %w", err)
	}
	return result, nil
}

// EString encodes text into a simple character encoding.
func EString(s string, c *charmap.Charmap) (result []byte, err error) {
	if !utf8.Valid([]byte(s)) {
		return result, fmt.Errorf("estring: %w", ErrUTF8)
	}
	encoder := c.NewEncoder()
	reader := transform.NewReader(strings.NewReader(s), encoder)
	result, err = ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("estring ioutil readall error: %w", err)
	}
	return result, nil
}

// D437 decodes IBM Code Page 437 encoded text.
func D437(s string) (result []byte, err error) {
	result, err = DString(s, charmap.CodePage437)
	if err != nil {
		return nil, fmt.Errorf("decode code page 437: %w", err)
	}
	return result, nil
}

// E437 encodes text into IBM Code Page 437.
func E437(s string) (result []byte, err error) {
	result, err = EString(s, charmap.CodePage437)
	if err != nil {
		return nil, fmt.Errorf("encode code page 437: %w", err)
	}
	return result, nil
}
