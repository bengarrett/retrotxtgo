package convert

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
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
	LightVertical = 179
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
	// Open Box ␣.
	OpenBox = 9251
)

const (
	row8  = 128
	row8f = 143
	row9  = 144
	row9f = 159
	rowA  = 160
	rowE  = 224
	cp037 = "IBM037"
	cp858 = "IBM00858"
	ibm   = "ibm"
	iso11 = "iso-8859-11"
	msdos = "msdos"
	u32   = "UTF-32"
	u32be = "UTF-32BE"
	u32le = "UTF-32LE"
	win   = "windows"
)

// Encoder returns the named character set encoder.
func Encoder(name string) (encoding.Encoding, error) {
	// use charmap string
	for _, c := range charmap.All {
		if fmt.Sprint(c) == name {
			return c, nil
		}
	}
	// use iana names or alias
	if e, err := ianaindex.IANA.Encoding(name); err == nil && e != nil {
		return e, nil
	}
	// use html index name (only used with HTML compatible encodings)
	// https://encoding.spec.whatwg.org/#names-and-labels
	if e, err := htmlindex.Get(name); err == nil && e != nil {
		return e, nil
	}
	s := shorten(name)
	a := encodeAlias(s)
	if a == iso11 {
		// ISO-8859-11 uses the same characters as Windows 847
		// except for 9 characters in rows 8 and 9.
		// https://en.wikipedia.org/wiki/ISO/IEC_8859-11#Code_page_874_(IBM)_/_9066
		return charmap.XUserDefined, nil
	}
	if e := encodeUTF32(a); e != nil {
		return e, nil
	}
	e, err := ianaindex.IANA.Encoding(s)
	if err != nil {
		e, err = ianaindex.IANA.Encoding(a)
	}
	if err != nil {
		e, err = ianaindex.IANA.Encoding(encodeAlias(name))
	}
	if err != nil || e == nil {
		if a == "" {
			return e, fmt.Errorf("%q: %w", name, ErrName)
		}
		return e, fmt.Errorf("name %q or alias %q: %w", name, a, ErrName)
	}
	return e, nil
}

// encodeUTF32 initializes common UTF-32 encodings.
func encodeUTF32(name string) encoding.Encoding {
	// UTF-32... doesn't return a match in ianaindex.IANA
	switch strings.ToUpper(name) {
	case u32:
		return utf32.UTF32(utf32.LittleEndian, utf32.UseBOM)
	case u32be:
		return utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM)
	case u32le:
		return utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM)
	}
	return nil
}

// Humanize the encoding by using an shorter, less formal name.
func Humanize(name string) string {
	if _, err := Encoder(name); err != nil {
		return ""
	}
	return encodeAlias(shorten(name))
}

// shorten the name to a custom name, a common name or an alias.
func shorten(name string) string { // nolint:gocognit,cyclop
	n, l := strings.ToLower(name), len(name)
	switch {
	case l > 3 && n[:3] == "cp-":
		return n[3:]
	case l > 2 && n[:2] == "cp":
		return n[2:]
	case l > 14 && n[:14] == "ibm code page ":
		return ibm + n[14:]
	case l > 4 && n[:4] == "ibm-":
		return n[4:]
	case l > 3 && n[:3] == ibm:
		return n[3:]
	case l > 4 && n[:4] == "oem-":
		return n[4:]
	case n == "windows code page 858":
		return cp858
	case l > 8 && n[:8] == "windows-":
		return n[8:]
	case l > 7 && n[:7] == win:
		return n[7:]
	case l > 7 && n[:7] == "iso8859":
		return "iso-8859-" + n[7:]
	case l > 9 && n[:9] == "iso 8859-":
		return "iso-8859-" + n[9:]
	}
	return ""
}

// encodeAlias returns a valid IANA index encoding name from a shorten name or alias.
func encodeAlias(name string) string {
	// list of valid tables
	// https://github.com/golang/text/blob/v0.3.2/encoding/charmap/maketables.go
	// a switch is used instead of a map to avoid typos with duplicate values
	if s := encodingIBM(name); s != "" {
		return s
	}
	if s := encodingMisc(name); s != "" {
		return s
	}
	if s := encodingWin(name); s != "" {
		return s
	}
	if s := encodingISO(name); s != "" {
		return s
	}
	if s := encodingUnicode(name); s != "" {
		return s
	}
	return ""
}

// encodingIBM returns a valid IANA index encoding name for IBM codepages using a custom alias.
func encodingIBM(name string) string {
	switch name {
	case "37", "037":
		return cp037
	case "437", "dos", "ibmpc", msdos, "us", "pc-8", "latin-us":
		return "IBM437"
	case "850", "latini":
		return "IBM850"
	case "852", "latinii":
		return "IBM852"
	case "855":
		return "IBM855"
	case "858":
		return cp858
	case "860":
		return "IBM860"
	case "862":
		return "IBM862"
	case "863":
		return "IBM863"
	case "865":
		return "IBM865"
	case "866":
		return "IBM866"
	case "1047":
		return "IBM1047"
	case "1140", "ibm1140":
		return "IBM01140"
	}
	return ""
}

// encodingMisc returns a valid IANA index encoding name using a custom alias.
func encodingMisc(name string) string {
	switch name {
	case "878", "20866", "koi8r":
		return "KOI8-R"
	case "1168", "21866", "koi8u":
		return "KOI8-U"
	case "10000", "macroman", "mac-roman", "mac os roman":
		return "Macintosh"
	case "shift jis", "shiftjis":
		return "shift_jis"
	case "ebcdic", ibm:
		return cp037
	case "iso88598e", "iso88598i", "iso88596e", "iso88596i":
		l := len(name)
		return fmt.Sprintf("ISO-8859-%v-%v", name[l-2:l-1], name[l-1:])
	}
	return ""
}

// encodingIBM returns a valid IANA index encoding name for Microsoft codepages using a custom alias.
func encodingWin(name string) string {
	switch name {
	case "874":
		return "Windows-874"
	case "1250":
		return "Windows-1250"
	case "1251":
		return "Windows-1251"
	case "1252", "1004", "win", win:
		return "Windows-1252"
	case "1253":
		return "Windows-1253"
	case "1254":
		return "Windows-1254"
	case "1255":
		return "Windows-1255"
	case "1256":
		return "Windows-1256"
	case "1257":
		return "Windows-1257"
	case "1258":
		return "Windows-1258"
	}
	return ""
}

// encodingIBM returns a valid IANA index encoding name for ISO codepages using a custom alias.
func encodingISO(name string) string {
	switch name {
	case "5", "1124", "28595":
		return "ISO-8859-5"
	case "6", "1089", "28596":
		return "ISO-8859-6"
	case "7", "813", "28597":
		return "ISO-8859-7"
	case "8", "916", "1125", "28598":
		return "ISO-8859-8"
	case "9", "920", "28599":
		return "ISO-8859-9"
	case "10", "919", "28600":
		return "ISO-8859-10"
	case "11", iso11:
		return strings.ToUpper(iso11)
	case "13", "921", "28603":
		return "ISO-8859-13"
	case "14", "28604":
		return "ISO-8859-14"
	case "15", "923", "28605":
		return "ISO-8859-15"
	case "16", "28606", "iso885916":
		return "ISO-8859-16"
	default:
		if s := encodingEurope(name); s != "" {
			return s
		}
	}
	return ""
}

// encodingIBM returns a valid IANA index encoding name for Latin codepages using a custom alias.
func encodingEurope(name string) string {
	switch name {
	case "1", "819", "28591":
		return "ISO-8859-1"
	case "2", "1111", "28592":
		return "ISO-8859-2"
	case "3", "913", "28593":
		return "ISO-8859-3"
	case "4", "914", "28594":
		return "ISO-8859-4"
	}
	return ""
}

// encodingIBM returns a valid IANA index encoding name for Unicode using a custom alias.
func encodingUnicode(name string) string {
	switch strings.ToLower(name) {
	case "utf16":
		return "UTF-16" // Go will use the byte-order-mark
	case "16be", "utf16b", "utf16be", "utf-16-be":
		return "UTF-16BE"
	case "16le", "utf16l", "utf16le", "utf-16-le":
		return "UTF-16LE"
	case "utf32", "utf-32":
		return u32 // Go will use the byte-order-mark
	case "32be", "utf32b", "utf32be", "utf-32be", "utf-32-be":
		return u32be
	case "32le", "utf32l", "utf32le", "utf-32le", "utf-32-le":
		return u32le
	}
	return ""
}

// Swap transforms character map and control codes into UTF-8 unicode runes.
func (c *Convert) Swap() *Convert {
	if len(c.Output) == 0 {
		return nil
	}
	c = c.swaps()
	if c.lineBreaks {
		c.LineBreaks()
	}
	switch c.Input.Encoding {
	case charmap.CodePage037, charmap.CodePage1047, charmap.CodePage1140:
		if c.Input.table {
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
	return c
}

func (c *Convert) swaps() *Convert {
	if len(c.Flags.SwapChars) == 0 {
		return c
	}
	replacers, unknown := []rune{}, []string{}
	for _, val := range c.Flags.SwapChars {
		nr := newReplacer(val)
		if nr == nil {
			unknown = append(unknown, val)
		}
		replacers = append(replacers, nr...)
	}
	if len(unknown) > 0 {
		fmt.Fprintln(os.Stderr, str.Info(), "unsupported --swap-chars values:", strings.Join(unknown, ","))
	}
	replace := runes.Map(func(r rune) rune {
		if !replacer(r, replacers...) {
			return r
		}
		return swap(r)
	})
	s, _, err := transform.String(replace, string(c.Output))
	if err != nil {
		logs.Fatal(err)
	}
	c.Output = []rune(s)
	return c
}

func replacer(r rune, replacers ...rune) bool {
	for _, s := range replacers {
		if s == r {
			return true
		}
	}
	return false
}

func newReplacer(s string) []rune {
	switch strings.ToLower(s) {
	case "n", "nul", "null", "00", "0":
		return []rune{NUL, SymbolNUL}
	case "b", "bar", "7c", "124":
		return []rune{VerticalBar}
	case "h", "house", "home", "7f", "127":
		return []rune{DEL, House}
	case "p", "pipe", "b3", "179":
		return []rune{LightVertical, LightVerticalU}
	case "r", "root", "squareroot", "fb", "251":
		return []rune{SquareRoot, SquareRootU}
	case "s", "space", "20", "32":
		return []rune{SP}
	}
	return nil
}

func swap(code rune) rune {
	return map[rune]rune{
		NUL:       SP,
		SymbolNUL: SP,
		// ¦
		VerticalBar: BrokenBar,
		// Δ
		DEL:   Delta,
		House: Delta,
		// ⎮
		LightVertical:  IntegralExtension,
		LightVerticalU: IntegralExtension,
		// ✓
		SquareRoot:  CheckMark,
		SquareRootU: CheckMark,
		// ␣
		SP: OpenBox,
	}[code]
}

// ANSIControls replaces out all ←[ and ␛[ character matches with functional ANSI escape controls.
func (c *Convert) ANSIControls() *Convert {
	if c == nil {
		return nil
	}
	if len(c.Output) == 0 {
		log.Fatal(ErrChainANSI)
	}
	for i, r := range c.Output {
		if i+1 >= len(c.Output) {
			continue
		}
		if r == LeftwardsArrow && c.Output[i+1] == LeftSquareBracket {
			c.Output[i] = ESC // replace ←[
			continue
		}
		if r == SymbolESC && c.Output[i+1] == LeftSquareBracket {
			c.Output[i] = ESC // replace ␛[
		}
	}
	return c
}

// LineBreaks will try to guess the line break representation as a 2 byte value.
// A guess of Unix will return [10, 0], Windows [13, 10], otherwise a [0, 0] value is returned.
func (c *Convert) LineBreaks() {
	c.Input.lineBreak = filesystem.LineBreaks(true, c.Output...)
}

// RunesControls switches out C0 and C1 ASCII controls with Unicode Control Picture represenations.
func (c *Convert) RunesControls() {
	if len(c.Output) == 0 {
		return
	}
	const z = byte(row8)
	for i := 0; i < len(c.Output); i++ {
		r := c.Output[i]
		if c.skipIgnores(i) {
			continue
		}
		if c.skipLineBreaks(i) {
			if c.Input.lineBreak == [2]rune{CR, 0} {
				c.Output[i] = LF
			}
			i++
			continue
		}
		if r >= NUL && r <= US {
			c.Output[i] = picture(byte(rune(z) + r))
		}
	}
}

// RunesControlsEBCDIC switches out EBCDIC controls with Unicode Control Picture represenations.
func (c *Convert) RunesControlsEBCDIC() {
	if len(c.Output) == 0 {
		return
	}
	const z = byte(row8)
	for i := 0; i < len(c.Output); i++ {
		r := c.Output[i]
		if c.skipIgnores(i) {
			continue
		}
		if r >= NUL && r <= US {
			c.Output[i] = picture(byte(rune(z) + r))
		}
	}
}

// RunesDOS switches out C0, C1 and other controls with PC/MS-DOS picture glyphs.
func (c *Convert) RunesDOS() {
	if len(c.Output) == 0 {
		return
	}
	// ASCII C0 = row 1, C1 = row 2
	ctrls := [32]string{
		string(picture(row8 + byte(0))),
		"\u263A", "\u263B", "\u2665", "\u2666", "\u2663", "\u2660",
		"\u2022", "\u25D8", "\u25CB", "\u25D9", "\u2642", "\u2640",
		"\u266A", "\u266B", "\u263C", "\u25BA", "\u25C4", "\u2195",
		"\u203C", "\u00B6", "\u00A7", "\u25AC", "\u21A8", "\u2191",
		"\u2193", "\u2192", "\u2190", "\u221F", "\u2194", "\u25B2", "\u25BC",
	}
	for i := 0; i < len(c.Output); i++ {
		r := c.Output[i]
		if c.skipIgnores(i) {
			continue
		}
		if c.skipLineBreaks(i) {
			if c.Input.lineBreak == [2]rune{13, 0} {
				c.Output[i] = LF // swap CR with LF
			}
			i++
			continue
		}
		switch {
		case r >= NUL && r <= US:
			c.Output[i], _ = utf8.DecodeRuneInString(ctrls[r]) // c0, c1 controllers
		case r == DEL:
			house := [3]byte{0xE2, 0x8C, 0x82} // ⌂
			c.Output[i], _ = utf8.DecodeRune(house[:])
		case r == NBSP:
			noBreakSpace := [2]byte{0xC2, 0xA0} // NBSP
			c.Output[i], _ = utf8.DecodeRune(noBreakSpace[:])
		}
	}
}

// RunesEBCDIC switches out EBCDIC IBM mainframe controls with Unicode picture represenations.
// Where no appropriate picture exists a space placeholder is used.
func (c *Convert) RunesEBCDIC() {
	if len(c.Output) == 0 {
		return
	}
	for i := 0; i < len(c.Output); i++ {
		r := c.Output[i]
		if c.skipIgnores(i) {
			continue
		}
		if c.control(i, r) {
			continue
		}
	}
}

// control switches out an EBCDIC IBM mainframe control with Unicode picture representation.
func (c *Convert) control(i int, r rune) (skip bool) { //nolint:funlen,cyclop
	if i > len(c.Output) {
		return false
	}
	const (
		ht  = 0x89
		del = 0xA1
		nl  = 0xA4
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
	switch r {
	case HT:
		c.Output[i] = picture(ht)
		return false
	case DEL:
		c.Output[i] = picture(del)
		return false
	case nel:
		if c.lineBreaks {
			c.Output[i] = LF // Go will automatically convert this to CRLF on Windows
			return true
		}
		c.Output[i] = picture(nl)
		return false
	case BS:
		c.Output[i] = picture(bs)
		return false
	case LF:
		c.Output[i] = picture(lf)
		return false
	case ETB:
		c.Output[i] = picture(etb)
		return false
	case ESC:
		c.Output[i] = picture(esc)
		return false
	case ENQ:
		c.Output[i] = picture(enq)
		return false
	case ACK:
		c.Output[i] = picture(ack)
		return false
	case BEL:
		c.Output[i] = picture(bel)
		return false
	case SYN:
		c.Output[i] = picture(syn)
		return false
	case Dash:
		c.Output[i] = picture(eot)
		return false
	case DC4:
		c.Output[i] = picture(dc4)
		return false
	case NAK:
		c.Output[i] = picture(nak)
		return false
	case SUB:
		c.Output[i] = picture(sub)
		return false
	default:
		c.miscCtrls(i, r)
	}
	return false
}

// miscCtrls switches out an EBCDIC control with Unicode Control Picture representation.
// Controls included those shared with ASCII C0+C1, NBSP and unprintables.
func (c *Convert) miscCtrls(i int, r rune) {
	switch r {
	case Nbsp:
		noBreakSpace := [2]byte{0xC2, 0xA0}
		c.Output[i], _ = utf8.DecodeRune(noBreakSpace[:])
	case NUL, SOH, STX, ETX, VT, FF, CR, SO, SI, DLE, DC1, DC2, DC3, CAN, EM, FS, GS, RS, US:
		// shared controls with ASCII C0+C1
		c.Output[i] = picture(row8 + byte(r))
	case EOT, InvertedExclamation:
		// unprintable controls
		c.Output[i] = rune(SP)
	case Cent, Negation, PlusMinus:
		return
	default:
		c.outOfRange(i, r)
	}
}

// outOfRange replaces EBCDIC runes that are out of range of
// valid 8-bit ASCII tables with a space placeholder.
func (c *Convert) outOfRange(i int, r rune) {
	const skipA, skipB, skipC, skipD = 0xA0, 0xBF, 0xC0, 0xFF
	switch {
	case
		r >= skipA && r <= skipB,
		r >= skipC && r <= skipD,
		r >= row8 && r <= row8f,
		r >= row9 && r <= row9f:
		c.Output[i] = rune(SP)
	}
}

// RunesKOI8 blanks out unused C0, C1 and other controls spaces for Russian sets.
func (c *Convert) RunesKOI8() {
	if len(c.Output) == 0 {
		return
	}
	for i := 0; i < len(c.Output); i++ {
		r := c.Output[i]
		if c.skipLineBreaks(i) {
			i++
			continue
		}
		switch {
		case r >= NUL && r <= US:
			c.Output[i] = rune(SP)
		case r == DEL:
			c.Output[i] = rune(SP)
		case r == Replacement:
			c.Output[i] = rune(SP)
		}
	}
}

// RunesLatin blanks out unused C0, C1 and other controls spaces for ISO Latin sets.
func (c *Convert) RunesLatin() {
	if len(c.Output) == 0 {
		return
	}
	for i := 0; i < len(c.Output); i++ {
		r := c.Output[i]
		if c.skipLineBreaks(i) {
			i++
			continue
		}
		switch {
		case r >= NUL && r <= US:
			c.Output[i] = rune(SP)
		case r >= DEL && r <= row9f:
			c.Output[i] = rune(SP)
		case r == Replacement:
			c.Output[i] = rune(SP)
		}
	}
}

// RunesMacintosh replaces specific Mac OS Roman characters with Unicode picture represenations.
func (c *Convert) RunesMacintosh() {
	const (
		command = iota + 17 // ⌘
		shift               // ⇧
		option              // ⌥
		control             // ⌃

		z = byte(row8)
	)
	for i := 0; i < len(c.Output); i++ {
		r := c.Output[i]
		if c.skipLineBreaks(i) {
			i++
			continue
		}
		switch r {
		case command:
			commandR := [3]byte{0xE2, 0x8C, 0x98}
			c.Output[i], _ = utf8.DecodeRune(commandR[:])
		case shift:
			shiftR := [3]byte{0xE2, 0x87, 0xA7}
			c.Output[i], _ = utf8.DecodeRune(shiftR[:])
		case option:
			optionR := [3]byte{0xE2, 0x8C, 0xA5}
			c.Output[i], _ = utf8.DecodeRune(optionR[:])
		case control:
			controlR := [3]byte{0xE2, 0x8C, 0x83}
			c.Output[i], _ = utf8.DecodeRune(controlR[:])
		case DEL:
			c.Output[i] = SymbolDEL
		case Replacement:
			c.Output[i] = rune(SP)
		default:
			if r >= NUL && r <= US {
				c.Output[i] = picture(byte(rune(z) + r))
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
	for i, r := range c.Output {
		switch {
		case r == backslash:
			c.Output[i] = rune(yen)
		case r == tilde:
			c.Output[i] = rune(overline)
		case r == DEL:
			c.Output[i] = picture(del)
		case r > DEL && r <= rowA,
			r >= rowE && r <= NBSP:
			c.Output[i] = rune(SP)
		}
	}
}

// RunesWindows tweaks some Unicode picture represenations for Windows-125x sets.
func (c *Convert) RunesWindows() {
	for i, r := range c.Output {
		switch r {
		case DEL:
			c.Output[i] = SymbolDEL
		case Replacement:
			c.Output[i] = rune(SP)
		}
	}
}

// RunesUTF8 tweaks some Unicode picture represenations for UTF-8 Basic Latin.
func (c *Convert) RunesUTF8() {
	for i, r := range c.Output {
		switch {
		case r == DEL:
			c.Output[i] = SymbolDEL
		case r > DEL && r < rowA:
			c.Output[i] = rune(SP)
		}
	}
}

// picture converts a byte value to a Unicode Control Picture rune.
func picture(b byte) rune {
	// code points U+2400 to U+243F
	controlPicture := [3]byte{0xE2, 0x90, b}
	r, _ := utf8.DecodeRune(controlPicture[:])
	return r
}

// equalLB returns true when r matches the single or multi-byte, line break character runes.
func equalLB(r, nl [2]rune) bool {
	// single-byte line break
	if nl[1] == 0 {
		return nl[0] == r[0]
	}
	// mutli-byte
	return bytes.Equal([]byte{byte(r[0]), byte(r[1])},
		[]byte{byte(nl[0]), byte(nl[1])})
}

// skipLineBreaks returns true if rune is a linebreak.
func (c *Convert) skipLineBreaks(i int) bool {
	if !c.lineBreaks {
		return false
	}
	l, r0, r1 := len(c.Output)-1, c.Output[i], rune(0)
	if i < l {
		// check for multi-byte line breaks
		r1 = c.Output[i+1]
	}
	if equalLB([2]rune{r0, r1}, c.Input.lineBreak) {
		return true
	}
	return false
}

// skipIgnores returns true if the rune should be skipped.
func (c *Convert) skipIgnores(i int) bool {
	for _, ign := range c.ignores {
		if c.Output[i] == ign {
			return true
		}
	}
	return false
}
