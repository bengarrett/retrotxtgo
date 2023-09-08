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
func DString(s string, c *charmap.Charmap) ([]byte, error) {
	decoder := c.NewDecoder()
	reader := transform.NewReader(strings.NewReader(s), decoder)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("dstring ioutil readall error: %w", err)
	}
	return b, nil
}

// EString encodes text into a simple character encoding.
func EString(s string, c *charmap.Charmap) ([]byte, error) {
	if !utf8.Valid([]byte(s)) {
		return nil, fmt.Errorf("estring: %w", ErrUTF8)
	}
	encoder := c.NewEncoder()
	reader := transform.NewReader(strings.NewReader(s), encoder)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("estring ioutil readall error: %w", err)
	}
	return b, nil
}

// D437 decodes IBM Code Page 437 encoded text.
func D437(s string) ([]byte, error) {
	b, err := DString(s, charmap.CodePage437)
	if err != nil {
		return nil, fmt.Errorf("decode code page 437: %w", err)
	}
	return b, nil
}

// E437 encodes text into IBM Code Page 437.
func E437(s string) ([]byte, error) {
	b, err := EString(s, charmap.CodePage437)
	if err != nil {
		return nil, fmt.Errorf("encode code page 437: %w", err)
	}
	return b, nil
}
