//Package convert is extends Go's x/text/encoding capability to convert legacy text
// to UTF-8.
package convert

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/japanese"
)

// Set data for transformation into UTF-8.
type Set struct {
	B        []byte
	R        []rune
	Encoding encoding.Encoding
	Newline  bool
	src      []byte
}

// Transform byte data from named character map text encoding into UTF-8.
func (s *Set) Transform(name string) (runes int, err error) {
	if name == "" {
		name = "UTF-8"
	}
	if s.Encoding, err = Encoding(name); err != nil {
		return runes, err
	}
	//fmt.Println("Encoding", s.Encoding, "is UTF-8", utf8.Valid(s.B), "length:", len(s.B))
	//fmt.Println(convert.HexEncode(string(s.B)))
	if len(s.B) == 0 {
		return runes, nil
	}
	s.src = s.B
	// only convert if data is not UTF-8
	if utf8.Valid(s.B) {
		s.R = bytes.Runes(s.B)
		return utf8.RuneCount(s.B), nil
	}
	if s.B, err = s.Encoding.NewDecoder().Bytes(s.B); err != nil {
		return runes, err
	}
	s.R = bytes.Runes(s.B)
	return utf8.RuneCount(s.B), nil
}

var (
	cp437C0 = []string{"\u0000", "\u263A", "\u263B", "\u2665", "\u2666", "\u2663", "\u2660", "\u2022", "\u25D8", "\u25CB", "\u25D9", "\u2642", "\u2640", "\u266A", "\u266B", "\u263C"}
	cp437C1 = []string{"\u25BA", "\u25C4", "\u2195", "\u203C", "\u00B6", "\u00A7", "\u25AC", "\u21A8", "\u2191", "\u2193", "\u2192", "\u2190", "\u221F", "\u2194", "\u25B2", "\u25BC"}
)

// BOM is the UTF-8 byte order mark prefix.
var BOM = func() []byte {
	return []byte{239, 187, 191} // 0xEF,0xBB,0xBF
}

// Encoding returns the named character set encoding.
func Encoding(name string) (encoding.Encoding, error) {
	// use iana names or alias
	if enc, err := ianaindex.IANA.Encoding(name); err == nil && enc != nil {
		return enc, err
	}
	// use html index name (only used with HTML compatible encodings)
	// https://encoding.spec.whatwg.org/#names-and-labels
	if enc, err := htmlindex.Get(name); err == nil && enc != nil {
		return enc, err
	}
	n := encodingAlias(shorten(name))
	enc, err := ianaindex.IANA.Encoding(n)
	if err != nil {
		err = fmt.Errorf("%s %q → %q", err, name, n)
		return enc, err
	}
	return enc, err
}

// shorten name to a custom/common names or aliases
func shorten(name string) (n string) {
	n = strings.ToLower(name)
	switch {
	case len(n) > 3 && n[:3] == "cp-":
		n = n[3:]
	case len(n) > 2 && n[:2] == "cp":
		n = n[2:]
	case len(n) > 14 && n[:14] == "ibm code page ":
		n = "ibm" + n[14:]
	case len(n) > 4 && n[:4] == "ibm-":
		n = n[4:]
	case len(n) > 3 && n[:3] == "ibm":
		n = n[3:]
	case len(n) > 4 && n[:4] == "oem-":
		n = n[4:]
	case n == "windows code page 858":
		n = "IBM00858"
	case len(n) > 8 && n[:8] == "windows-":
		n = n[8:]
	case len(n) > 7 && n[:7] == "windows":
		n = n[7:]
	case len(n) > 7 && n[:7] == "iso8859":
		n = "iso-8859-" + n[7:]
	case len(n) > 9 && n[:9] == "iso 8859-":
		n = "iso-8859-" + n[9:]
	}
	return n
}

// encodingAlias returns a valid IANA index encoding name from a shorten name or alias.
func encodingAlias(name string) (n string) {
	// list of valid tables
	// https://github.com/golang/text/blob/v0.3.2/encoding/charmap/maketables.go
	switch name {
	case "37", "037":
		n = "IBM037"
	case "437", "dos", "ibmpc", "msdos", "us", "pc-8", "latin-us":
		n = "IBM437"
	case "850", "latini":
		n = "IBM850"
	case "852", "latinii":
		n = "IBM852"
	case "855":
		n = "IBM855"
	case "858":
		n = "IBM00858"
	case "860":
		n = "IBM860"
	case "862":
		n = "IBM862"
	case "863":
		n = "IBM863"
	case "865":
		n = "IBM865"
	case "866":
		n = "IBM866"
	case "1047":
		n = "IBM1047"
	case "1140", "ibm1140":
		n = "IBM01140"
	case "1", "819", "28591":
		n = "ISO-8859-1"
	case "2", "1111", "28592":
		n = "ISO-8859-2"
	case "3", "913", "28593":
		n = "ISO-8859-3"
	case "4", "914", "28594":
		n = "ISO-8859-4"
	case "5", "1124", "28595":
		n = "ISO-8859-5"
	case "6", "1089", "28596":
		n = "ISO-8859-6"
	case "7", "813", "28597":
		n = "ISO-8859-7"
	case "8", "916", "1125", "28598":
		n = "ISO-8859-8"
	case "9", "920", "28599":
		n = "ISO-8859-9"
	case "10", "919", "28600":
		n = "ISO-8859-10"
	case "11", "874", "iso-8859-11":
		n = "Windows-874" // "iso-8859-11" causes a panic
	case "13", "921", "28603":
		n = "ISO-8859-13"
	case "14", "28604":
		n = "ISO-8859-14"
	case "15", "923", "28605":
		n = "ISO-8859-15"
	case "16", "28606":
		n = "ISO-8859-16"
	case "878", "20866":
		n = "KOI8-R"
	case "1168", "21866":
		n = "KOI8-U"
	case "10000", "macroman", "mac-roman", "mac os roman":
		n = "Macintosh"
	case "1250":
		n = "Windows-1250"
	case "1251":
		n = "Windows-1251"
	case "1252", "1004", "win", "windows":
		n = "Windows-1252"
	case "1253":
		n = "Windows-1253"
	case "1254":
		n = "Windows-1254"
	case "1255":
		n = "Windows-1255"
	case "1256":
		n = "Windows-1256"
	case "1257":
		n = "Windows-1257"
	case "1258":
		n = "Windows-1258"
	case "koi8r":
		n = "KOI8R"
	case "koi8u":
		n = "KOI8U"
	case "shift jis":
		n = "ShiftJIS"
	}
	return n
}

// MakeBytes generates an 8-bit unsigned int container ready to hold legacy code point values.
func MakeBytes() (m []byte) {
	m = make([]byte, 256)
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
	if cut := bytes.IndexByte(s.B, 26); cut > 0 {
		s.B = s.B[:cut]
	}
}

// Newlines will try to guess the newline representation as a 2 byte value.
// A guess of Unix will return [10, 0], Windows [13, 10], otherwise a [0, 0] value is returned.
func (s Set) Newlines() [2]rune {
	// scan data for possible newlines
	c := []struct {
		abbr  string
		count int
	}{
		{"lf", 0},   // linux, unix, amiga...
		{"cr", 0},   // 8-bit micros
		{"crlf", 0}, // windows, dos, cp/m...
		{"lfcr", 0}, // acorn bbc micro
		{"nl", 0},   // ibm ebcdic encodings
	}
	l := len(s.R) - 1 // range limit
	for i, r := range s.R {
		switch r {
		case 10:
			if i < l && s.R[i+1] == 13 {
				c[3].count++ // lfcr
				continue
			}
			if i != 0 && s.R[i-1] == 13 {
				// crlf (already counted)
				continue
			}
			c[0].count++
		case 13:
			if i < l && s.R[i+1] == 10 {
				c[2].count++ // crlf
				continue
			}
			if i != 0 && s.R[i-1] == 10 {
				// lfcr (already counted)
				continue
			}
			c[1].count++
		case 21:
			c[4].count++
		case 155:
			// atascii (not currently used)
		}
	}
	// sort results
	sort.SliceStable(c, func(i, j int) bool {
		return c[i].count > c[j].count
	})
	switch c[0].abbr {
	case "lf":
		return [2]rune{10}
	case "cr":
		return [2]rune{13}
	case "crlf":
		return [2]rune{13, 10}
	case "lfcr":
		return [2]rune{10, 13}
	case "nl":
		return [2]rune{21}
	}
	return [2]rune{}
}

// Skip the rune if it matches the newline characters.
func skip(r rune, nl [2]rune) bool {
	switch r {
	case 0: // avoid false positives with 1 byte newlines
		return false
	case nl[0], nl[1]:
		return true
	}
	return false
}

// Swap transforms common ...
func (s *Set) Swap() {
	switch s.Encoding {
	case charmap.CodePage037, charmap.CodePage1047, charmap.CodePage1140:
		s.RunesEBCDIC()
	case charmap.CodePage437, charmap.CodePage850, charmap.CodePage852, charmap.CodePage855,
		charmap.CodePage858, charmap.CodePage860, charmap.CodePage862, charmap.CodePage863,
		charmap.CodePage865, charmap.CodePage866:
		s.RunesDOS()
	case charmap.ISO8859_1, charmap.ISO8859_2, charmap.ISO8859_3, charmap.ISO8859_4, charmap.ISO8859_5,
		charmap.ISO8859_6, charmap.ISO8859_7, charmap.ISO8859_8, charmap.ISO8859_9, charmap.ISO8859_10,
		charmap.ISO8859_13, charmap.ISO8859_14, charmap.ISO8859_15, charmap.ISO8859_16,
		charmap.Windows874:
		s.RunesLatin()
	case charmap.ISO8859_6E, charmap.ISO8859_6I, charmap.ISO8859_8E, charmap.ISO8859_8I:
		s.RunesControls()
		s.RunesLatin()
	case charmap.KOI8R, charmap.KOI8U:
		s.RunesKOI8()
	case charmap.Macintosh:
		s.RunesMacintosh()
	case charmap.Windows1250, charmap.Windows1251, charmap.Windows1252, charmap.Windows1253,
		charmap.Windows1254, charmap.Windows1255, charmap.Windows1256, charmap.Windows1257, charmap.Windows1258:
		s.RunesControls()
		s.RunesWindows()
	case japanese.ShiftJIS:
		log.Fatal("japanese character sets are not working")
		// s.RunesControls()
		// s.RunesShiftJIS()
	default:
		s.RunesControls()
	}
}

// SwapANSI replaces out all ←[ and ␛[ character matches with functional ANSI escape controls.
func (s *Set) SwapANSI() {
	for i, r := range s.R {
		if i+1 >= len(s.R) {
			continue
		}
		if r == 8592 && s.R[i+1] == 91 {
			s.R[i] = 27 // replace ←[
		} else if r == 9243 && s.R[i+1] == 91 {
			s.R[i] = 27 // replace ␛[
		}
	}
}

// RunesControls switches out C0 and C1 ASCII controls with Unicode picture represenations.
func (s *Set) RunesControls() {
	if len(s.R) == 0 {
		return
	}
	const z = byte(0x80)
	var nl [2]rune
	if s.Newline {
		nl = s.Newlines()
	}
	for i, r := range s.R {
		if s.Newline && skip(r, nl) {
			continue
		}
		switch {
		case r >= 0x00 && r <= 0x1f:
			s.R[i] = decode(byte(rune(z) + r))
		}
	}
}

// RunesDOS switches out C0, C1 and other controls with PC/MS-DOS picture glyphs.
func (s *Set) RunesDOS() {
	if len(s.R) == 0 {
		return
	}
	var (
		ctrls = append(cp437C0, cp437C1...)
		nl    [2]rune
	)
	if s.Newline {
		nl = s.Newlines()
	}
	for i, r := range s.R {
		if s.Newline && skip(r, nl) {
			continue
		}
		switch {
		case r == 0x00:
			s.R[i] = decode(0x80) // NUL
		case r > 0x00 && r <= 0x1f:
			s.R[i], _ = utf8.DecodeRuneInString(ctrls[r]) // c0, c1 controllers
		case r == 0x7c: // todo: flag option?
			s.R[i], _ = utf8.DecodeRune([]byte{0xc2, 0xa6}) // ¦
		case r == 0x7f:
			s.R[i], _ = utf8.DecodeRune([]byte{0xe2, 0x8c, 0x82}) // ⌂
		case r == 0xff:
			s.R[i], _ = utf8.DecodeRune([]byte{0xc2, 0xa0}) // NBSP
		case r == 0x0d:
			s.R[i], _ = utf8.DecodeRune([]byte("\n"))
		}
	}
}

// RunesEBCDIC switches out EBCDIC IBM mainframe controls with Unicode picture represenations.
// Where no appropriate picture exists a space is used.
func (s *Set) RunesEBCDIC() {
	if len(s.R) == 0 {
		return
	}
	var nl [2]rune
	if s.Newline {
		nl = s.Newlines()
	}
	for i, r := range s.R {
		if s.Newline && skip(r, nl) {
			continue
		}
		switch r {
		case 9:
			s.R[i] = decode(0x89) // HT
		case 127:
			s.R[i] = decode(0xa1) // DEL
		case 133:
			s.R[i] = decode(0xa4) // NL
		case 8:
			s.R[i] = decode(0x88) // BS
		case 10:
			s.R[i] = decode(0x8A) // LF
		case 23:
			s.R[i] = decode(0x97) // ETB
		case 27: // x27 = ESC = x1B unicode
			s.R[i] = decode(0x9B) // ESC
		case 5:
			s.R[i] = decode(0x85) // ENQ
		case 6:
			s.R[i] = decode(0x86) // ACK
		case 7:
			s.R[i] = decode(0x87) // BEL
		case 22:
			s.R[i] = decode(0x96) // SYN
		case 150:
			s.R[i] = decode(0x84) // EOT
		case 20:
			s.R[i] = decode(0x94) // DC4
		case 21:
			s.R[i] = decode(0x95) // NAK
		case 26:
			s.R[i] = decode(0x9A) // SUB
		case 160:
			s.R[i], _ = utf8.DecodeRune([]byte{0xc2, 0xa0}) // NBSP
		case 0x00, 0x01, 0x02, 0x03, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
			0x10, 0x11, 0x12, 0x13, 0x18, 0x19, 0x1c, 0x1d, 0x1e, 0x1f:
			// shared controls with ASCII C0+C1
			s.R[i] = decode(0x80 + byte(r))
		case 0x9c, 0x86, 0x97, 0x8d, 0x8e,
			0x9D, 0x87, 0x92, 0x8f,
			0x80, 0x81, 0x82, 0x83, 0xA1, 0x84, 0x88, 0x89, 0x8a, 0x8b, 0x8c,
			0x90, 0x91, 0x93, 0x94, 0x95, 0x04, 0x98, 0x99, 0x9a, 0x9b, 0x9e,
			0x9f:
			// unprintable controls
			s.R[i] = rune(32)
			// default:
			// 	fmt.Printf("[%d→%v/%X %q]\n", i, r, r, r)
		}
	}
}

// RunesKOI8 blanks out unused C0, C1 and other controls spaces for Russian sets.
func (s *Set) RunesKOI8() {
	if len(s.R) == 0 {
		return
	}
	var nl [2]rune
	if s.Newline {
		nl = s.Newlines()
	}
	for i, r := range s.R {
		if s.Newline && skip(r, nl) {
			continue
		}
		switch {
		case r >= 0x00 && r <= 0x1f:
			s.R[i] = rune(32)
		case r == 0x7f:
			s.R[i] = rune(32)
		case r == 65533:
			s.R[i] = rune(32)
		}
	}
}

// RunesLatin blanks out unused C0, C1 and other controls spaces for ISO Latin sets.
func (s *Set) RunesLatin() {
	if len(s.R) == 0 {
		return
	}
	var nl [2]rune
	if s.Newline {
		nl = s.Newlines()
	}
	for i, r := range s.R {
		if s.Newline && skip(r, nl) {
			continue
		}
		switch {
		case r >= 0x00 && r <= 0x1f:
			s.R[i] = rune(32)
		case r >= 0x7f && r <= 0x9f:
			s.R[i] = rune(32)
		case r == 65533:
			s.R[i] = rune(32)
		}
	}
}

// RunesMacintosh replaces specific Mac OS Roman characters with Unicode picture represenations.
func (s *Set) RunesMacintosh() {
	const z = byte(0x80)
	var nl [2]rune
	if s.Newline {
		nl = s.Newlines()
	}
	for i, r := range s.R {
		if s.Newline && skip(r, nl) {
			continue
		}
		switch r {
		case 0x11:
			s.R[i], _ = utf8.DecodeRune([]byte{0xe2, 0x8c, 0x98}) // ⌘
		case 0x12:
			s.R[i], _ = utf8.DecodeRune([]byte{0xe2, 0x87, 0xa7}) // ⇧
		case 0x13:
			s.R[i], _ = utf8.DecodeRune([]byte{0xe2, 0x8c, 0xa5}) // ⌥
		case 0x14:
			s.R[i], _ = utf8.DecodeRune([]byte{0xe2, 0x8c, 0x83}) // ⌃
		case 0x7f:
			s.R[i] = decode(0xa1) // DEL
		case 65533:
			s.R[i] = rune(32)
		default:
			switch {
			case r >= 0x00 && r <= 0x1f:
				s.R[i] = decode(byte(rune(z) + r))
			}
		}
	}
}

// RunesWindows tweaks some Unicode picture represenations for Windows-125x sets.
func (s *Set) RunesWindows() {
	for i, r := range s.R {
		switch r {
		case 0x7f:
			s.R[i] = decode(0xa1) // DEL
		case 65533:
			s.R[i] = rune(32)
		}
	}
}

// utf8.DecodeRune([]byte{0xe2, 0x90, b})
func decode(b byte) (r rune) {
	r, _ = utf8.DecodeRune([]byte{0xe2, 0x90, b})
	return r
}
