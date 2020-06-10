//Package transform is extends Go's x/text/encoding capability to convert legacy text
// to UTF-8.
package transform

import (
	"bytes"
	"log"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
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
	if s.Encoding, err = Encoding(name); err != nil {
		return runes, err
	}
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

	//fmt.Println("--> runes", string(s.R), string(s.B), fmt.Sprintf("%X", s.B))
	// for i, a := range s.R {
	// 	fmt.Printf("%2d. %v 0x%X %q\n", i, a, a, string(a))
	// 	fmt.Printf("    -> %v 0x%X\n", s.src[i], s.src[i])
	// }

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
	return ianaindex.IANA.Encoding(name)
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

// MakeMap generates an 8-bit unsigned int container ready to hold legacy code point values.
func MakeMap() (m []byte) {
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

// SwapAll transforms all common ...
func SwapAll(b []byte) []byte {
	var s Set
	s.B = b
	s.Newline = true
	s.Swap()
	return s.B
}

// Swap transforms common ...
func (s *Set) Swap() {
	// s.CutEOF()
	//s.SwapNuls()
	// s.SwapPipes()
	// s.SwapDels()
	//s.SwapNBSP()
	//s.SwapControlsDOS()
	//s.SwapControlPics()
	//s.SwapControlsIBM()
	//s.SwapANSI()
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
	case charmap.KOI8R, charmap.KOI8U:
		s.RunesKOI8()
	case charmap.Macintosh:
		// https://en.wikipedia.org/wiki/Mac_OS_Roman aka aliases cp10000 etc.
		s.RunesControls()
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

// Drop EUC-JP,ISO2022JP,

// SwapANSI replaces out all ←[ character combinations with the ANSI escape control.
func (s *Set) SwapANSI() {
	s.B = bytes.ReplaceAll(s.B, []byte("←["), []byte{27, 91})
}

// SwapNuls replaces the ASCII codepoint 0 NULL value with the Unicode 0020 SP space value.
func (s *Set) SwapNuls() {
	s.B = bytes.ReplaceAll(s.B, []byte{0}, []byte("\u0020"))
}

// SwapPipes replaces the ASCII codepoint 124 broken bar (or pipe) with the Unicode 00A6 ¦ broken pipe symbol.
func (s *Set) SwapPipes() {
	s.B = bytes.ReplaceAll(s.B, []byte{124}, []byte("\u00A6"))
}

// SwapDels replaces the ASCII codepoint 127 delete with the Unicode codepoint 2302 ⌂ house symbol.
func (s *Set) SwapDels() {
	s.B = bytes.ReplaceAll(s.B, []byte{127}, []byte("\u2302"))
}

// SwapNBSP replaces the ASCII codepoint 255 no-break-space with Unicode codepoint C2A0 no-break.
func (s *Set) SwapNBSP() {
	s.B = bytes.ReplaceAll(s.B, []byte{255}, []byte("\uC2A0"))
}

// RunesDOS switches out C0, C1 and other controls with PC/MS-DOS picture glyphs.
func (s *Set) RunesDOS() {
	if len(s.R) == 0 {
		return
	}
	var ctrls = append(cp437C0, cp437C1...)
	for i, r := range s.R {
		switch {
		case r == 0x00:
			s.R[i], _ = utf8.DecodeRune([]byte{0xE2, 0x90, 0x80}) // NUL
		case r > 0x00 && r <= 0x1f:
			s.R[i], _ = utf8.DecodeRuneInString(ctrls[r]) // c0, c1 controllers
		case r == 0x7c: // todo: flag option?
			s.R[i], _ = utf8.DecodeRune([]byte{0xc2, 0xa6}) // ¦
		case r == 0x7f:
			s.R[i], _ = utf8.DecodeRune([]byte{0xE2, 0x8c, 0x82}) // ⌂
		}
	}
}

// RunesControls switches out C0 and C1 ASCII controls with Unicode picture represenations.
func (s *Set) RunesControls() {
	if len(s.R) == 0 {
		return
	}
	const z = byte(0x80)
	for i, r := range s.R {
		switch {
		case r >= 0x00 && r <= 0x1f:
			b := byte(rune(z) + r)
			s.R[i], _ = utf8.DecodeRune([]byte{0xE2, 0x90, b})
		}
	}
}

// RunesKOI8 blanks out unused C0, C1 and other controls spaces for Russian sets.
func (s *Set) RunesKOI8() {
	if len(s.R) == 0 {
		return
	}
	for i, r := range s.R {
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
	for i, r := range s.R {
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
	for i, r := range s.R {
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
			s.R[i], _ = utf8.DecodeRune([]byte{0xe2, 0x90, 0xA1}) // DEL
		case 65533:
			s.R[i] = rune(32)
		}
	}
}

// RunesWindows tweaks some Unicode picture represenations for Windows-125x sets.
func (s *Set) RunesWindows() {
	for i, r := range s.R {
		switch r {
		case 0x7f:
			s.R[i], _ = utf8.DecodeRune([]byte{0xe2, 0x90, 0xA1}) // DEL
		case 65533:
			s.R[i] = rune(32)
		}
	}
}

func decode(b byte) (r rune) {
	r, _ = utf8.DecodeRune([]byte{0xe2, 0x90, b})
	return r
}

// RunesEBCDIC switches out EBCDIC IBM mainframe controls with Unicode picture represenations.
// Where no appropriate picture exists a space is used.
func (s *Set) RunesEBCDIC() {
	if len(s.R) == 0 {
		return
	}
	for i, r := range s.R {
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
		default:
			//fmt.Printf("[%d→%v/%X %q]\n", i, r, r, r)
		}
	}
}
