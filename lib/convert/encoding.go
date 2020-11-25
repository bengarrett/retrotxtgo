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

const (
	// NUL Null control code.
	NUL = iota
	// SOH Start of heading.
	SOH
	// STX Start of text.
	STX
	// ETX End of text.
	ETX
	// EOT End of transmission.
	EOT
	// ENQ Enquiry.
	ENQ
	// ACK Acknowledge.
	ACK
	// BEL Bell or alert.
	BEL
	// BS Backspace.
	BS
	// HT Horizontal tabulation.
	HT
	// LF Line feed.
	LF
	// VT Vertical tabulation.
	VT
	// FF Form feed.
	FF
	// CR Carriage return.
	CR
	// SO Shift out.
	SO
	// SI Shift in.
	SI
	// DLE Data Link Escape.
	DLE
	// DC1 Device control one.
	DC1
	// DC2 Device control two.
	DC2
	// DC3 Device control three.
	DC3
	// DC4 Device control four.
	DC4
	// NAK Negative acknowledge.
	NAK
	// SYN Synchronous idle.
	SYN
	// ETB End of transmission block.
	ETB
	// CAN Cancel.
	CAN
	// EM End of medium.
	EM
	// SUB Substitute.
	SUB
	// ESC Escape.
	ESC
	// FS File separator.
	FS
	// GS Group separator.
	GS
	// RS Record separator.
	RS
	// US Unit separator.
	US
	// SP Space.
	SP
)

const (
	// LeftSquareBracket [.
	LeftSquareBracket = 91
	// VerticalBar |.
	VerticalBar = 124
	// DEL Delete.
	DEL = 127
	// Dash Hyphen -.
	Dash = 150
	// Nbsp Non-breaking space.
	Nbsp = 160
	// InvertedExclamation ¡.
	InvertedExclamation = 161
	// Cent ¢.
	Cent = 162
	// BrokenBar ¦.
	BrokenBar = 166
	// Negation ¬.
	Negation = 172
	// PlusMinus ±.
	PlusMinus = 177
	// LightVertical light vertical │.
	LightVertical = 179 // TODO: test 178 vs 179
	// SquareRoot Square root √.
	SquareRoot = 251
	// NBSP Non-breaking space.
	NBSP = 255
	// Delta Δ.
	Delta = 916
	// LeftwardsArrow ←.
	LeftwardsArrow = 8592
	// SquareRootU Unicode square root √.
	SquareRootU = 8730
	// House ⌂.
	House = 8962
	// IntegralExtension ⎮.
	IntegralExtension = 9134
	// SymbolNUL ␀.
	SymbolNUL = 9216
	// SymbolESC ␛.
	SymbolESC = 9243
	// SymbolDEL ␡.
	SymbolDEL = 9249
	// LightVerticalU Box drawing light vertical │.
	LightVerticalU = 9474
	// CheckMark ✓.
	CheckMark = 10003
	// Replacement character �.
	Replacement = 65533
)

const (
	row8  = 128
	row8f = 143
	row9  = 144
	row9f = 159
	rowA  = 160
	rowE  = 224
	iso11 = "iso-8859-11"
)

// Chars are characters with alternative runes.
type Chars map[int]rune

var (
	// ErrChainANSI ansi() is a chain method.
	ErrChainANSI = errors.New("ansi() is a chain method that is to be used in conjunction with swap: c.swap().ansi()")
	// ErrName unknown encoding.
	ErrName = errors.New("encoding could not match name or alias")
)

// Characters map code page 437 characters with alternative runes.
func Characters() Chars {
	return Chars{
		NUL:           SP,
		VerticalBar:   BrokenBar,
		DEL:           Delta,             // Δ
		LightVertical: IntegralExtension, // ⎮
		SquareRoot:    CheckMark,         // ✓
	}
}

// Encoding returns the named character set encoding.
func Encoding(name string) (encoding.Encoding, error) {
	// use charmap string
	for _, c := range charmap.All {
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
	if a == iso11 {
		// ISO-8859-11 uses the same characters as Windows 847
		// except for 9 characters in rows 8 and 9.
		// https://en.wikipedia.org/wiki/ISO/IEC_8859-11#Code_page_874_(IBM)_/_9066
		return charmap.Windows874, nil
	}
	enc, err := ianaindex.IANA.Encoding(s)
	if err != nil {
		enc, err = ianaindex.IANA.Encoding(a)
	}
	if err != nil {
		enc, err = ianaindex.IANA.Encoding(encodingAlias(name))
	}
	if err != nil || enc == nil {
		if a == "" {
			return enc, fmt.Errorf("%q: %w", name, ErrName)
		}
		return enc, fmt.Errorf("name %q or alias %q: %w", name, a, ErrName)
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

// shorten name to a custom/common names or aliases.
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
	case "11", iso11:
		n = strings.ToUpper(iso11)
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
	case "874":
		n = "Windows-874"
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
	const debug = false
	if c.useBreaks {
		c.LineBreaks()
	}
	if debug {
		println(fmt.Sprintf("line break detected: %+v", c.useBreaks))
	}
	switch c.Source.E {
	case charmap.CodePage037, charmap.CodePage1047, charmap.CodePage1140:
		if c.table {
			c.RunesEBCDIC()
		}
		c.RunesControlsEBCDIC()
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
	if len(c.Flags.SwapChars) > 0 {
		for i := 0; i < c.len; i++ {
			if s := c.runeSwap(c.Runes[i]); s >= 0 {
				c.Runes[i] = s
			}
		}
	}
	return c
}

// ANSIControls replaces out all ←[ and ␛[ character matches with functional ANSI escape controls.
func (c *Convert) ANSIControls() {
	if c == nil {
		return
	}
	if c.len == 0 {
		log.Fatal(ErrChainANSI)
	}
	for i, r := range c.Runes {
		if i+1 >= c.len {
			continue
		}
		if r == LeftwardsArrow && c.Runes[i+1] == LeftSquareBracket {
			c.Runes[i] = ESC // replace ←[
		} else if r == SymbolESC && c.Runes[i+1] == LeftSquareBracket {
			c.Runes[i] = ESC // replace ␛[
		}
	}
}

// LineBreaks will try to guess the line break representation as a 2 byte value.
// A guess of Unix will return [10, 0], Windows [13, 10], otherwise a [0, 0] value is returned.
func (c *Convert) LineBreaks() {
	c.lineBreak = filesystem.LineBreaks(true, c.Runes...)
}

// RunesControls switches out C0 and C1 ASCII controls with Unicode picture represenations.
func (c *Convert) RunesControls() {
	if len(c.Runes) == 0 {
		return
	}
	const z = byte(row8)
	for i := 0; i < c.len; i++ {
		r := c.Runes[i]
		if c.skipIgnores(i) {
			continue
		}
		if c.skipLineBreaks(i) {
			if c.lineBreak == [2]rune{CR, 0} {
				c.Runes[i] = LF
			}
			i++
			continue
		}
		if r >= NUL && r <= US {
			c.Runes[i] = decode(byte(rune(z) + r))
		}
	}
}

// RunesControlsEBCDIC switches out EBCDIC controls with Unicode picture represenations.
func (c *Convert) RunesControlsEBCDIC() {
	if len(c.Runes) == 0 {
		return
	}
	//
	const z = byte(row8)
	for i := 0; i < c.len; i++ {
		r := c.Runes[i]
		if c.skipIgnores(i) {
			continue
		}
		if r >= NUL && r <= US {
			c.Runes[i] = decode(byte(rune(z) + r))
		}
	}
	fmt.Println(string(c.Runes))
}

// RunesDOS switches out C0, C1 and other controls with PC/MS-DOS picture glyphs.
func (c *Convert) RunesDOS() {
	if len(c.Runes) == 0 {
		return
	}
	// ASCII C0 = row 1, C1 = row 2
	var ctrls = [32]string{string(decode(row8 + byte(0))),
		"\u263A", "\u263B", "\u2665", "\u2666", "\u2663", "\u2660",
		"\u2022", "\u25D8", "\u25CB", "\u25D9", "\u2642", "\u2640",
		"\u266A", "\u266B", "\u263C", "\u25BA", "\u25C4", "\u2195",
		"\u203C", "\u00B6", "\u00A7", "\u25AC", "\u21A8", "\u2191",
		"\u2193", "\u2192", "\u2190", "\u221F", "\u2194", "\u25B2", "\u25BC"}
	for i := 0; i < c.len; i++ {
		r := c.Runes[i]
		if c.skipIgnores(i) {
			continue
		}
		if c.skipLineBreaks(i) {
			if c.lineBreak == [2]rune{13, 0} {
				c.Runes[i] = LF // swap CR with LF
			}
			i++
			continue
		}
		switch {
		case r >= NUL && r <= US:
			c.Runes[i], _ = utf8.DecodeRuneInString(ctrls[r]) // c0, c1 controllers
		case r == DEL:
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xe2, 0x8c, 0x82}) // ⌂
		case r == NBSP:
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
	const (
		ht  = 0x89
		del = 0xa1
		nl  = 0xa4
		bs  = 0x88
		lf  = 0x8A
		etb = 0x97
		esc = 0x9B
		enq = 0x85
		ack = 0x86
		bel = 0x87
		syn = 0x96
		eot = 0x84
		dc4 = 0x94
		nak = 0x95
		sub = 0x9A
		nel = 133
	)
	const skipA, skipB, skipC, skipD = 0xA0, 0xBF, 0xC0, 0xFF
	for i := 0; i < c.len; i++ {
		r := c.Runes[i]
		if c.skipIgnores(i) {
			continue
		}
		switch r {
		case HT:
			c.Runes[i] = decode(ht)
		case DEL:
			c.Runes[i] = decode(del)
		case nel:
			if c.useBreaks {
				// Go will automatically convert this to CRLF on Windows
				c.Runes[i] = LF
				continue
			}
			c.Runes[i] = decode(nl)
		case BS:
			c.Runes[i] = decode(bs)
		case LF:
			c.Runes[i] = decode(lf)
		case ETB:
			c.Runes[i] = decode(etb)
		case ESC:
			c.Runes[i] = decode(esc)
		case ENQ:
			c.Runes[i] = decode(enq)
		case ACK:
			c.Runes[i] = decode(ack)
		case BEL:
			c.Runes[i] = decode(bel)
		case SYN:
			c.Runes[i] = decode(syn)
		case Dash:
			c.Runes[i] = decode(eot)
		case DC4:
			c.Runes[i] = decode(dc4)
		case NAK:
			c.Runes[i] = decode(nak)
		case SUB:
			c.Runes[i] = decode(sub)
		case Nbsp:
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xc2, 0xa0})
		case NUL, SOH, STX, ETX, VT, FF, CR, SO, SI, DLE, DC1, DC2, DC3, CAN, EM, FS, GS, RS, US:
			// shared controls with ASCII C0+C1
			c.Runes[i] = decode(row8 + byte(r))
		case EOT, InvertedExclamation:
			// unprintable controls
			c.Runes[i] = rune(SP)
		case Cent, Negation, PlusMinus:
			// keep these symbols
		default:
			switch {
			case
				r >= skipA && r <= skipB,
				r >= skipC && r <= skipD,
				r >= row8 && r <= row8f,
				r >= row9 && r <= row9f:
				c.Runes[i] = rune(SP)
			}
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
		if c.skipLineBreaks(i) {
			i++
			continue
		}
		switch {
		case r >= NUL && r <= US:
			c.Runes[i] = rune(SP)
		case r == DEL:
			c.Runes[i] = rune(SP)
		case r == Replacement:
			c.Runes[i] = rune(SP)
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
		if c.skipLineBreaks(i) {
			i++
			continue
		}
		switch {
		case r >= NUL && r <= US:
			c.Runes[i] = rune(SP)
		case r >= DEL && r <= row9f:
			c.Runes[i] = rune(SP)
		case r == Replacement:
			c.Runes[i] = rune(SP)
		}
	}
}

// RunesMacintosh replaces specific Mac OS Roman characters with Unicode picture represenations.
func (c *Convert) RunesMacintosh() {
	const z = byte(row8)
	const (
		command = iota + 17 // ⌘
		shift               // ⇧
		option              // ⌥
		control             // ⌃
	)
	for i := 0; i < c.len; i++ {
		r := c.Runes[i]
		if c.skipLineBreaks(i) {
			i++
			continue
		}
		switch r {
		case command:
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xe2, 0x8c, 0x98})
		case shift:
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xe2, 0x87, 0xa7})
		case option:
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xe2, 0x8c, 0xa5})
		case control:
			c.Runes[i], _ = utf8.DecodeRune([]byte{0xe2, 0x8c, 0x83})
		case DEL:
			c.Runes[i] = SymbolDEL
		case Replacement:
			c.Runes[i] = rune(SP)
		default:
			if r >= NUL && r <= US {
				c.Runes[i] = decode(byte(rune(z) + r))
			}
		}
	}
}

// RunesShiftJIS tweaks some Unicode picture represenations for Shift-JIS.
func (c *Convert) RunesShiftJIS() {
	const (
		backslash = 92     // /
		tilde     = 126    // ~
		yen       = 0xa5   // ¥
		overline  = 0x203e // ‾
		del       = 0xa1   // del
	)
	for i, r := range c.Runes {
		switch {
		case r == backslash:
			c.Runes[i] = rune(yen)
		case r == tilde:
			c.Runes[i] = rune(overline)
		case r == DEL:
			c.Runes[i] = decode(del)
		case r > DEL && r <= rowA,
			r >= rowE && r <= NBSP:
			c.Runes[i] = rune(SP)
		}
	}
}

// RunesWindows tweaks some Unicode picture represenations for Windows-125x sets.
func (c *Convert) RunesWindows() {
	for i, r := range c.Runes {
		switch r {
		case DEL:
			c.Runes[i] = SymbolDEL
		case Replacement:
			c.Runes[i] = rune(SP)
		}
	}
}

// RunesUTF8 tweaks some Unicode picture represenations for UTF-8 Basic Latin.
func (c *Convert) RunesUTF8() {
	for i, r := range c.Runes {
		switch {
		case r == DEL:
			c.Runes[i] = SymbolDEL
		case r > DEL && r < rowA:
			c.Runes[i] = rune(SP)
		}
	}
}

func decode(b byte) (r rune) {
	r, _ = utf8.DecodeRune([]byte{0xe2, 0x90, b})
	return r
}

// equalLB reports whether r matches the single or multi-byte, line break character runes.
func equalLB(r, nl [2]rune) bool {
	// single-byte line break
	if nl[1] == 0 {
		return nl[0] == r[0]
	}
	// mutli-byte
	return bytes.Equal([]byte{byte(r[0]), byte(r[1])},
		[]byte{byte(nl[0]), byte(nl[1])})
}

func (c *Convert) skipLineBreaks(i int) bool {
	if !c.useBreaks {
		return false
	}
	var l, r0, r1 = c.len - 1, c.Runes[i], rune(0)
	if i < l {
		// check for multi-byte line breaks
		r1 = c.Runes[i+1]
	}
	if equalLB([2]rune{r0, r1}, c.lineBreak) {
		return true
	}
	return false
}

func (c *Convert) runeSwap(r rune) rune {
	// todo: use transform
	// https://play.golang.org/p/unix7YjB8Dw
	// https://pkg.go.dev/golang.org/x/text/runes
	if !c.swap(r) {
		return -1
	}
	switch {
	case r == NUL:
		return Characters()[0]
	case r == SymbolNUL:
		return Characters()[0]
	case r == VerticalBar:
		return Characters()[VerticalBar]
	case r == House:
		return Characters()[DEL]
	case r == LightVerticalU:
		return Characters()[LightVertical]
	case r == SquareRootU:
		return Characters()[SquareRoot]
	}
	return -1
}

func (c *Convert) skipIgnores(i int) bool {
	for _, ign := range c.Output.ignores {
		if c.Runes[i] == ign {
			return true
		}
	}
	return false
}

func (c *Convert) swap(r rune) bool {
	chk := NUL
	switch r {
	case VerticalBar:
		chk = VerticalBar
	case House:
		chk = DEL
	case LightVerticalU:
		chk = LightVertical
	case SquareRootU:
		chk = SquareRoot
	}
	for _, c := range c.Flags.SwapChars {
		if c == chk {
			return true
		}
	}
	return false
}
