//Package transform is extends Go's x/text/encoding capability to convert legacy text
// to UTF-8.
package transform

import (
	"bytes"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/html/charset"
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

func MakeMap() [256]byte {
	var m [256]byte
	//encoding, _ := ianaindex.IANA.Encoding("cp437")
	for i := 0; i <= 255; i++ {
		// b := []byte{uint8(i)}
		// c := math.Mod(float64(i), 16)
		// t, _ := Transform(b, encoding)
		// if c == 0 {
		// 	fmt.Print("\n")
		// }
		// fmt.Printf("%x %s\t", b, t)
		m[i] = uint8(i)
	}
	return m
}

// ToBOM adds a UTF-8 byte order mark if it doesn't already exist.
func ToBOM(b []byte) []byte {
	if len(b) > 2 {
		if t := b[:3]; bytes.Equal(t, BOM()) {
			return b
		}
	}
	return append(BOM(), b...)
}

// UTF8 determines if a document is encoded as UTF-8.
func UTF8(b []byte) bool {
	_, name, _ := charset.DetermineEncoding(b, "text/plain")
	return name == "utf-8" // bool
}

func SwapAll(b []byte) []byte {
	var s Set
	s.Data = b
	s.SwapAll(true)
	return s.Data
}

func SwapRecommended(b []byte) []byte {
	var s Set
	s.Data = b
	s.SwapAll(false)
	return s.Data
}

// CutEOF cut text at the first DOS end-of-file marker.
func (s *Set) CutEOF() {
	if cut := bytes.IndexByte(s.Data, 26); cut > 0 {
		s.Data = s.Data[:cut]
	}
}

func (s *Set) SwapAll(nl bool) {
	s.CutEOF()
	s.SwapNuls()
	s.SwapPipes()
	s.SwapDels()
	s.SwapNBSP()
	s.SwapControls(nl)
	s.SwapANSI()
}

// SwapANSI switches out ←[ characters with ANSI escape codes.
func (s *Set) SwapANSI() {
	s.Data = bytes.ReplaceAll(s.Data, []byte("←["), []byte{27, 91})
}

// \u0000 should be swapped for SP \u0000 --nul-as-space (true)
func (s *Set) SwapNuls() {
	s.Data = bytes.ReplaceAll(s.Data, []byte{0}, []byte("\u0020"))
}

// \u007c (7C) [pipe] can be swapped for broken bar \u00A6 --pipe-as-broken-bar (false)
func (s *Set) SwapPipes() {
	s.Data = bytes.ReplaceAll(s.Data, []byte{124}, []byte("\u00A6"))
}

// // \u0127? (7F) [delete] can be swapped for a house \u2303 --del-as-house (false)
func (s *Set) SwapDels() {
	s.Data = bytes.ReplaceAll(s.Data, []byte{127}, []byte("\u2303"))
}

// // FF NBSP often displays a ?, it can be replaced with SP --nbsp-as-space (true)
func (s *Set) SwapNBSP() {
	s.Data = bytes.ReplaceAll(s.Data, []byte{255}, []byte("\u263B"))
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

// Transform ...
func (s *Set) Transform() (err error) {
	if s.Data, err = s.Encoding.NewDecoder().Bytes(s.Data); err != nil {
		return err
	}
	return nil
}
