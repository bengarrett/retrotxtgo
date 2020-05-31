//Package transform is extends Go's x/text/encoding capability to convert legacy text
// to UTF-8.
package transform

import (
	"bytes"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
)

// EOF
// is an ASCII Control-Z, code 26

// Set blah
type Set struct {
	Data     []byte
	Encoding encoding.Encoding
}

var (
	asciiC0 = []string{"\u0000", "\u263A", "\u263B", "\u2665", "\u2666", "\u2663", "\u2660", "\u2022", "\u25D8", "\u25CB", "\u25D9", "\u2642", "\u2640", "\u266A", "\u266B", "\u263C"}
	asciiC1 = []string{"\u25BA", "\u25C4", "\u2195", "\u203C", "\u00B6", "\u00A7", "\u25AC", "\u21A8", "\u2191", "\u2193", "\u2192", "\u2190", "\u221F", "\u2194", "\u25B2", "\u25BC"}
)

// BOM is the UTF-8 byte order mark prefix.
var BOM = func() []byte {
	return []byte{239, 187, 191} // 0xEF,0xBB,0xBF
}

// Encoding returns the named character set encoding.
func Encoding(name string) (encoding.Encoding, error) {
	return ianaindex.IANA.Encoding(Replace(name))
}

// Replace normalizes common unofficial character set aliases.
func Replace(name string) string {
	name = strings.ToLower(name)
	reg, _ := regexp.Compile("[^a-z0-9]+")
	name = reg.ReplaceAllString(name, "")
	name = strings.ReplaceAll(name, "iso8859", "iso-8859-")
	name = strings.ReplaceAll(name, "isolatin", "latin")
	name = strings.ReplaceAll(name, "windows125", "windows-125")
	name = strings.ReplaceAll(name, "win125", "windows-125")
	name = strings.ReplaceAll(name, "cp125", "windows-125")
	return name
}

// Valid determines if the named character set or alias is known.
func Valid(name string) bool {
	if name == "" {
		return false
	}
	name = Replace(name)
	if _, err := ianaindex.IANA.Encoding(name); err != nil {
		return false
	}
	return true
}

// Transform byte data from charmap text encoding to UTF-8.
func Transform(m *charmap.Charmap, px *[]byte) (runes int, encoded []byte, err error) {
	p := *px
	if len(p) == 0 {
		return 0, encoded, nil
	}
	// confirm encoding is not utf8
	if utf8.Valid(p) {
		return utf8.RuneCount(p), p, nil
	}
	// use cp437 by default if text is not utf8
	// TODO: add default-unknown.encoding setting
	if m == nil {
		m = charmap.CodePage437
	}
	// convert to utf8
	if encoded, err = m.NewDecoder().Bytes(p); err != nil {
		return 0, encoded, err
	}
	return utf8.RuneCount(encoded), encoded, nil
}

// MakeMap generates an 8-bit unsigned int container ready to hold legacy code point values.
func MakeMap() (m [256]byte) {
	for i := 0; i <= 255; i++ {
		m[i] = uint8(i)
	}
	return m
}

// AddBOM adds a UTF-8 byte order mark if it doesn't already exist.
func AddBOM(b []byte) []byte {
	if len(b) > 2 {
		if t := b[:3]; bytes.Equal(t, BOM()) {
			return b
		}
	}
	return append(BOM(), b...)
}

// CutEOF cut text at the first DOS end-of-file marker.
func (s *Set) CutEOF() {
	if cut := bytes.IndexByte(s.Data, 26); cut > 0 {
		s.Data = s.Data[:cut]
	}
}

// SwapAll transforms all common ...
func SwapAll(b []byte) []byte {
	var s Set
	s.Data = b
	s.Swap(true)
	return s.Data
}

// Swap transforms common ...
func (s *Set) Swap(nl bool) {
	s.CutEOF()
	s.SwapNuls()
	s.SwapPipes()
	s.SwapDels()
	s.SwapNBSP()
	s.SwapControls(nl)
	s.SwapANSI()
}

// SwapANSI replaces out all ←[ character combinations with the ANSI escape control.
func (s *Set) SwapANSI() {
	s.Data = bytes.ReplaceAll(s.Data, []byte("←["), []byte{27, 91})
}

// SwapNuls replaces the ASCII codepoint 0 NULL value with the Unicode 0020 SP space value.
func (s *Set) SwapNuls() {
	s.Data = bytes.ReplaceAll(s.Data, []byte{0}, []byte("\u0020"))
}

// SwapPipes replaces the ASCII codepoint 124 broken bar (or pipe) with the Unicode 00A6 ¦ broken pipe symbol.
func (s *Set) SwapPipes() {
	s.Data = bytes.ReplaceAll(s.Data, []byte{124}, []byte("\u00A6"))
}

// SwapDels replaces the ASCII codepoint 127 delete with the Unicode codepoint 2302 ⌂ house symbol.
func (s *Set) SwapDels() {
	s.Data = bytes.ReplaceAll(s.Data, []byte{127}, []byte("\u2302"))
}

// SwapNBSP replaces the ASCII codepoint 255 no-break-space with Unicode codepoint C2A0 no-break.
func (s *Set) SwapNBSP() {
	s.Data = bytes.ReplaceAll(s.Data, []byte{255}, []byte("\uC2A0"))
}

// SwapControls switches out C0 and C1 ASCII controls except for newlines.
func (s *Set) SwapControls(nl bool) {
	for i, u := range append(asciiC0, asciiC1...) {
		if nl {
			switch i {
			case 10, 13: // newlines
				continue
			}
		}
		s.Data = bytes.ReplaceAll(s.Data, []byte{uint8(i)}, []byte(u))
	}
}

// Transform byte data from charmap text encoding to UTF-8.
func (s *Set) Transform() (err error) {
	if s.Data, err = s.Encoding.NewDecoder().Bytes(s.Data); err != nil {
		return err
	}
	return nil
}
