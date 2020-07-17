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
	"golang.org/x/text/encoding/unicode"
	"retrotxt.com/retrotxt/lib/filesystem"
)

// Characters map code page 437 characters with alternative runes
var Characters = map[int]rune{
	0:   32,    // NUL
	124: 166,   // ¦
	127: 916,   // Δ
	178: 9134,  // ⎮
	251: 10003, // ✓
}

// Encoding returns the named character set encoding.
func Encoding(name string) (encoding.Encoding, error) {
	// use charmap string
	for _, c := range charmap.All {
		//fmt.Println("--->", fmt.Sprint(c), "--->", name)
		if fmt.Sprint(c) == name {
			return c, nil
		}
	}
	// use iana names or alias
	if enc, err := ianaindex.IANA.Encoding(name); err == nil && enc != nil {
		return enc, nil
	}
	// use html index name (only used with HTML compatible encodings)
	// https://encoding.spec.whatwg.org/#names-and-labels
	if enc, err := htmlindex.Get(name); err == nil && enc != nil {
		return enc, nil
	}
	s := shorten(name)
	a := encodingAlias(s)
	n := encodingAlias(name)
	enc, err := ianaindex.IANA.Encoding(s)
	if enc == nil {
		enc, err = ianaindex.IANA.Encoding(a)
	}
	if enc == nil {
		enc, err = ianaindex.IANA.Encoding(n)
	}
	if err != nil {
		return enc, fmt.Errorf("encoding could not match name %q or alias %q: %s",
			name, a, err)
	}
	return enc, nil
}

// Humanize the encoding by using an shorter, less formal name.
func Humanize(name string) string {
	if _, err := Encoding(name); err != nil {
		return ""
	}
	return encodingAlias(shorten(name))
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
	case "16", "28606", "iso885916":
		n = "ISO-8859-16"
	case "878", "20866", "koi8r":
		n = "KOI8-R"
	case "1168", "21866", "koi8u":
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
	case "shift jis", "shiftjis":
		n = "shift_jis"
	case "utf16":
		n = "UTF-16" // Go will use the byte-order-mark
	case "16be", "utf16b", "utf16be", "utf-16-be":
		n = "UTF-16BE"
	case "16le", "utf16l", "utf16le", "utf-16-le":
		n = "UTF-16LE"
	case "ebcdic", "ibm":
		n = "IBM037"
	case "iso88598e", "iso88598i", "iso88596e", "iso88596i":
		l := len(name)
		n = fmt.Sprintf("ISO-8859-%v-%v", name[l-2:l-1], name[l-1:])
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
	//println(fmt.Sprintf("newline detected: %+v", c.newlines))
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
		c.RunesControls()
		c.RunesShiftJIS()
	case unicode.UTF8, unicode.UTF8BOM:
		c.RunesControls()
		c.RunesUTF8()
	default:
	}
	if len(c.swapChars) > 0 {
		for i := 0; i < c.len; i++ {
			if s := c.runeSwap(c.Runes[i]); s >= 0 {
				c.Runes[i] = s
			}
		}
	}
	return c
}

// ANSI replaces out all ←[ and ␛[ character matches with functional ANSI escape controls.
func (c *Convert) ANSI() {
	if c == nil {
		return
	}
	if c.len == 0 {
		log.Fatal(errors.New("ansi() is a chain method that is to be used in conjuction with swap: c.swap().ansi()"))
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
	c.newlines = filesystem.Newlines(true, c.Runes...)
}

// RunesControls switches out C0 and C1 ASCII controls with Unicode picture represenations.
func (c *Convert) RunesControls() {
	if len(c.Runes) == 0 {
		return
	}
	const z = byte(0x80)
	for i := 0; i < c.len; i++ {
		r := c.Runes[i]
		if c.skipIgnores(i) {
			continue
		}
		if c.skipNewlines(i) {
			if c.newlines == [2]rune{13, 0} {
				c.Runes[i] = 10 // swap CR with LF
			}
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
	// ASCII C0 = row 1, C1 = row 2
	var ctrls = [32]string{"\u0000", "\u263A", "\u263B", "\u2665", "\u2666", "\u2663", "\u2660", "\u2022", "\u25D8", "\u25CB", "\u25D9", "\u2642", "\u2640", "\u266A", "\u266B", "\u263C",
		"\u25BA", "\u25C4", "\u2195", "\u203C", "\u00B6", "\u00A7", "\u25AC", "\u21A8", "\u2191", "\u2193", "\u2192", "\u2190", "\u221F", "\u2194", "\u25B2", "\u25BC"}
	for i := 0; i < c.len; i++ {
		r := c.Runes[i]
		if c.skipNewlines(i) {
			if c.newlines == [2]rune{13, 0} {
				c.Runes[i] = 10 // swap CR with LF
			}
			i++
			continue
		}
		switch {
		case r > 0x00 && r <= 0x1f:
			c.Runes[i], _ = utf8.DecodeRuneInString(ctrls[r]) // c0, c1 controllers
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
		if c.skipIgnores(i) {
			//i++
			continue
		}
		switch r {
		case 9:
			c.Runes[i] = decode(0x89) // HT
		case 127:
			c.Runes[i] = decode(0xa1) // DEL
		case 133:
			if c.newline {
				c.Runes[i] = 10 // Go will automatically convert this to CRLF on Windows
				continue
			}
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
		if c.skipNewlines(i) {
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
		if c.skipNewlines(i) {
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
		if c.skipNewlines(i) {
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

// RunesShiftJIS tweaks some Unicode picture represenations for Shift-JIS.
func (c *Convert) RunesShiftJIS() {
	for i, r := range c.Runes {
		switch {
		case r == 0x5c:
			c.Runes[i] = rune(0xa5) // ¥
		case r == 0x7e:
			c.Runes[i] = rune(0x203e) // ‾
		case r == 0x7f:
			c.Runes[i] = decode(0xa1) // DEL
		case r > 0x7f && r <= 0xa0,
			r >= 0xe0 && r <= 0xff:
			c.Runes[i] = rune(32)
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

// RunesUTF8 tweaks some Unicode picture represenations for UTF-8 Basic Latin.
func (c *Convert) RunesUTF8() {
	for i, r := range c.Runes {
		switch {
		case r == 0x7f:
			c.Runes[i] = decode(0xa1) // DEL
		case r > 0x7f && r < 0xa0:
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
		return nl[0] == r[0]
	}
	// mutli-byte
	return bytes.Equal([]byte{byte(r[0]), byte(r[1])},
		[]byte{byte(nl[0]), byte(nl[1])})
}

func (c Convert) skipNewlines(i int) bool {
	if !c.newline {
		return false
	}
	var l, r0, r1 = c.len - 1, c.Runes[i], rune(0)
	if i < l {
		// check for multi-byte newlines
		r1 = c.Runes[i+1]
	}
	if equalNL([2]rune{r0, r1}, c.newlines) {
		return true
	}
	return false
}

func (c Convert) runeSwap(r rune) rune {
	if !c.swap(r) {
		return -1
	}
	switch {
	case r == 0x00:
		return Characters[0] // NUL
	case r == 0x2400:
		return Characters[0]
	case r == 0x7c:
		return Characters[124] // ¦
	case r == 0x2302:
		return Characters[127] // ⌂
	case r == 0x2502:
		return Characters[178] // │
	case r == 0x221A:
		return Characters[251] // √
	}
	return -1
}

func (c Convert) skipIgnores(i int) bool {
	for _, ign := range c.ignores {
		if c.Runes[i] == ign {
			return true
		}
	}
	return false
}

func (c Convert) swap(r rune) bool {
	/*
		0:   32,    // NUL
		124: 166,   // ¦
		127: 916,   // Δ
		178: 9134,  // ⎮
		251: 10003, // ✓
	*/
	chk := 0
	switch r {
	case 0x7c:
		chk = 124
	case 0x2302:
		chk = 127
	case 0x2502:
		chk = 178
	case 0x221A:
		chk = 251
	}
	for _, c := range c.swapChars {
		if c == chk {
			return true
		}
	}
	return false
}
