package convert

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/japanese"
	"retrotxt.com/retrotxt/lib/filesystem"
)

// BOM is the UTF-8 byte order mark prefix.
func BOM() []byte {
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

// EndOfFile will cut text at the first DOS end-of-file marker.
func EndOfFile(b []byte) []byte {
	if cut := bytes.IndexByte(b, 26); cut > 0 {
		return b[:cut]
	}
	return b
}

// MakeBytes generates a 256 character or 8-bit container ready to hold legacy code point values.
func MakeBytes() (m []byte) {
	m = make([]byte, 256)
	for i := 0; i <= 255; i++ {
		m[i] = uint8(i)
	}
	return m
}

// Mark adds a UTF-8 byte order mark to the text if it doesn't already exist.
func Mark(b []byte) []byte {
	if len(b) > 2 {
		if t := b[:3]; bytes.Equal(t, BOM()) {
			return b
		}
	}
	return append(BOM(), b...)
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
	// a switch is used instead of a map to avoid typos with duplicate values
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
	case "utf16":
		n = "UTF-16" // Go will default to Big Endian
	case "16be", "utf16b", "utf16be", "utf-16-be":
		n = "UTF-16BE"
	case "16le", "utf16l", "utf16le", "utf-16-le":
		n = "UTF-16LE"
	}
	return n
}

// Swap transforms character map and control codes into UTF-8 unicode runes.
func (c *Convert) Swap() *Convert {
	if c.len == 0 {
		return nil
	}
	if c.newline {
		c.Newlines()
	}
	switch c.encode {
	case charmap.CodePage037, charmap.CodePage1047, charmap.CodePage1140:
		c.RunesEBCDIC()
	case charmap.CodePage437, charmap.CodePage850, charmap.CodePage852, charmap.CodePage855,
		charmap.CodePage858, charmap.CodePage860, charmap.CodePage862, charmap.CodePage863,
		charmap.CodePage865, charmap.CodePage866:
		c.RunesDOS()
	case charmap.ISO8859_1, charmap.ISO8859_2, charmap.ISO8859_3, charmap.ISO8859_4, charmap.ISO8859_5,
		charmap.ISO8859_6, charmap.ISO8859_7, charmap.ISO8859_8, charmap.ISO8859_9, charmap.ISO8859_10,
		charmap.ISO8859_13, charmap.ISO8859_14, charmap.ISO8859_15, charmap.ISO8859_16,
		charmap.Windows874:
		c.RunesLatin()
	case charmap.ISO8859_6E, charmap.ISO8859_6I, charmap.ISO8859_8E, charmap.ISO8859_8I:
		c.RunesControls()
		c.RunesLatin()
	case charmap.KOI8R, charmap.KOI8U:
		c.RunesKOI8()
	case charmap.Macintosh:
		c.RunesMacintosh()
	case charmap.Windows1250, charmap.Windows1251, charmap.Windows1252, charmap.Windows1253,
		charmap.Windows1254, charmap.Windows1255, charmap.Windows1256, charmap.Windows1257, charmap.Windows1258:
		c.RunesControls()
		c.RunesWindows()
	case japanese.ShiftJIS:
	default:
		c.RunesControls()
	}
	return c
}

// ANSI replaces out all ←[ and ␛[ character matches with functional ANSI escape controls.
func (c *Convert) ANSI() {
	if c == nil {
		return
	}
	if c.len == 0 {
		log.Fatal(errors.New("ANSI() is a chain method that is to be used in conjuction with Swap: c.Swap().ANSI()"))
	}
	for i, r := range c.Runes {
		if i+1 >= c.len {
			continue
		}
		if r == 8592 && c.Runes[i+1] == 91 {
			c.Runes[i] = 27 // replace ←[
		} else if r == 9243 && c.Runes[i+1] == 91 {
			c.Runes[i] = 27 // replace ␛[
		}
	}
}

// Newlines will try to guess the newline representation as a 2 byte value.
// A guess of Unix will return [10, 0], Windows [13, 10], otherwise a [0, 0] value is returned.
func (c *Convert) Newlines() {
	c.newlines = filesystem.Newlines(c.Runes)
}

// RunesControls switches out C0 and C1 ASCII controls with Unicode picture represenations.
func (c *Convert) RunesControls() {
	if len(c.Runes) == 0 {
		return
	}
	const z = byte(0x80)
	for i := 0; i < c.len; i++ {
		r := c.Runes[i]
		if c.skipRune(i) {
			i++
			continue
		}
		switch {
		case r >= 0x00 && r <= 0x1f:
			c.Runes[i] = decode(byte(rune(z) + r))
		}
	}
}

// RunesDOS switches out C0, C1 and other controls with PC/MS-DOS picture glyphs.
func (c *Convert) RunesDOS() {
	if len(c.Runes) == 0 {
		return
	}
	var (
		c0    = []string{"\u0000", "\u263A", "\u263B", "\u2665", "\u2666", "\u2663", "\u2660", "\u2022", "\u25D8", "\u25CB", "\u25D9", "\u2642", "\u2640", "\u266A", "\u266B", "\u263C"}
		c1    = []string{"\u25BA", "\u25C4", "\u2195", "\u203C", "\u00B6", "\u00A7", "\u25AC", "\u21A8", "\u2191", "\u2193", "\u2192", "\u2190", "\u221F", "\u2194", "\u25B2", "\u25BC"}
		ctrls = append(c0, c1...)
	)
	for i := 0; i < c.len; i++ {
		r := c.Runes[i]
		if c.skipRune(i) {
			i++
			continue
		}
		switch {
		case r == 0x00:
			c.Runes[i] = decode(0x80) // NUL
		case r > 0x00 && r <= 0x1f:
			c.Runes[i], _ = utf8.DecodeRuneInString(ctrls[r]) // c0, c1 controllers
		case r == 0x7c: // todo: add a user flag toggle
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xc2, 0xa6}) // ¦
		case r == 0x7f:
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xe2, 0x8c, 0x82}) // ⌂
		case r == 0xff:
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xc2, 0xa0}) // NBSP
		}
	}
}

// RunesEBCDIC switches out EBCDIC IBM mainframe controls with Unicode picture represenations.
// Where no appropriate picture exists a space is used.
func (c *Convert) RunesEBCDIC() {
	if len(c.Runes) == 0 {
		return
	}
	for i := 0; i < c.len; i++ {
		r := c.Runes[i]
		if c.skipRune(i) {
			i++
			continue
		}
		switch r {
		case 9:
			c.Runes[i] = decode(0x89) // HT
		case 127:
			c.Runes[i] = decode(0xa1) // DEL
		case 133:
			c.Runes[i] = decode(0xa4) // NL
		case 8:
			c.Runes[i] = decode(0x88) // BS
		case 10:
			c.Runes[i] = decode(0x8A) // LF
		case 23:
			c.Runes[i] = decode(0x97) // ETB
		case 27: // x27 = ESC = x1B unicode
			c.Runes[i] = decode(0x9B) // ESC
		case 5:
			c.Runes[i] = decode(0x85) // ENQ
		case 6:
			c.Runes[i] = decode(0x86) // ACK
		case 7:
			c.Runes[i] = decode(0x87) // BEL
		case 22:
			c.Runes[i] = decode(0x96) // SYN
		case 150:
			c.Runes[i] = decode(0x84) // EOT
		case 20:
			c.Runes[i] = decode(0x94) // DC4
		case 21:
			c.Runes[i] = decode(0x95) // NAK
		case 26:
			c.Runes[i] = decode(0x9A) // SUB
		case 160:
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xc2, 0xa0}) // NBSP
		case 0x00, 0x01, 0x02, 0x03, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
			0x10, 0x11, 0x12, 0x13, 0x18, 0x19, 0x1c, 0x1d, 0x1e, 0x1f:
			// shared controls with ASCII C0+C1
			c.Runes[i] = decode(0x80 + byte(r))
		case 0x9c, 0x86, 0x97, 0x8d, 0x8e,
			0x9D, 0x87, 0x92, 0x8f,
			0x80, 0x81, 0x82, 0x83, 0xA1, 0x84, 0x88, 0x89, 0x8a, 0x8b, 0x8c,
			0x90, 0x91, 0x93, 0x94, 0x95, 0x04, 0x98, 0x99, 0x9a, 0x9b, 0x9e,
			0x9f:
			// unprintable controls
			c.Runes[i] = rune(32)
			// default:
			// 	fmt.Printf("[%d→%v/%X %q]\n", i, r, r, r)
		}
	}
}

// RunesKOI8 blanks out unused C0, C1 and other controls spaces for Russian sets.
func (c *Convert) RunesKOI8() {
	if len(c.Runes) == 0 {
		return
	}
	for i := 0; i < c.len; i++ {
		r := c.Runes[i]
		if c.skipRune(i) {
			i++
			continue
		}
		switch {
		case r >= 0x00 && r <= 0x1f:
			c.Runes[i] = rune(32)
		case r == 0x7f:
			c.Runes[i] = rune(32)
		case r == 65533:
			c.Runes[i] = rune(32)
		}
	}
}

// RunesLatin blanks out unused C0, C1 and other controls spaces for ISO Latin sets.
func (c *Convert) RunesLatin() {
	if len(c.Runes) == 0 {
		return
	}
	for i := 0; i < c.len; i++ {
		r := c.Runes[i]
		if c.skipRune(i) {
			i++
			continue
		}
		switch {
		case r >= 0x00 && r <= 0x1f:
			c.Runes[i] = rune(32)
		case r >= 0x7f && r <= 0x9f:
			c.Runes[i] = rune(32)
		case r == 65533:
			c.Runes[i] = rune(32)
		}
	}
}

// RunesMacintosh replaces specific Mac OS Roman characters with Unicode picture represenations.
func (c *Convert) RunesMacintosh() {
	const z = byte(0x80)
	for i := 0; i < c.len; i++ {
		r := c.Runes[i]
		if c.skipRune(i) {
			i++
			continue
		}
		switch r {
		case 0x11:
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xe2, 0x8c, 0x98}) // ⌘
		case 0x12:
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xe2, 0x87, 0xa7}) // ⇧
		case 0x13:
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xe2, 0x8c, 0xa5}) // ⌥
		case 0x14:
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xe2, 0x8c, 0x83}) // ⌃
		case 0x7f:
			c.Runes[i] = decode(0xa1) // DEL
		case 65533:
			c.Runes[i] = rune(32)
		default:
			switch {
			case r >= 0x00 && r <= 0x1f:
				c.Runes[i] = decode(byte(rune(z) + r))
			}
		}
	}
}

// RunesWindows tweaks some Unicode picture represenations for Windows-125x sets.
func (c *Convert) RunesWindows() {
	for i, r := range c.Runes {
		switch r {
		case 0x7f:
			c.Runes[i] = decode(0xa1) // DEL
		case 65533:
			c.Runes[i] = rune(32)
		}
	}
}

func (c *Convert) ignore(r rune) {
	c.ignores = append(c.ignores, r)
}

// utf8.DecodeRune([]byte{0xe2, 0x90, b})
func decode(b byte) (r rune) {
	r, _ = utf8.DecodeRune([]byte{0xe2, 0x90, b})
	return r
}

// equalNL reports whether r matches the single or multi-byte, newline character runes.
func equalNL(r [2]rune, nl [2]rune) bool {
	// single-byte newline
	if nl[1] == 0 {
		if nl[0] == r[0] {
			return true
		}
		return false
	}
	// mutli-byte
	return bytes.Equal([]byte{byte(r[0]), byte(r[1])},
		[]byte{byte(nl[0]), byte(nl[1])})
}

func (c *Convert) skipRune(i int) bool {
	var l, r0, r1 = c.len - 1, c.Runes[i], rune(0)
	if i < l {
		// check for multi-byte newlines
		r1 = c.Runes[i+1]
	}
	if c.newline && equalNL([2]rune{r0, r1}, c.newlines) {
		return true
	}
	for _, ign := range c.ignores {
		if r0 == ign {
			return true
		}
	}
	return false
}
