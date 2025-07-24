// Package xud provides X User Defined, character encodings.
// It does not encode or decode text, only provides information about the encodings.
//
// This includes the early American Standards Association (ASA) ASCII character encodings.
// There are three ASA encodings, X3.4-1963, X3.4-1965, X3.4-1967 and one missing ISO 8859-11 encoding.
// These encodings are not compatible with each other.
//
// But the X3.4-1967 character codes are compatible with the ANSI X3.4-1977 and ANSI X3.4-1986 encodings.
// Which are also compatible with many of the IBM Code Page and ISO 8859-X encodings, as-well as Unicode.
package xud

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

// Named, numeric and alias values for the legacy ASA ASCII character encodings.
const (
	Name11  = "iso-8859-11" // name of ISO 8859-11
	Name63  = "ascii-63"    // name of ASA X3.4 1963
	Name65  = "ascii-65"    // name of ASA X3.4 1965
	Name67  = "ascii-67"    // name of ANSI X3.4 1967/77/86
	Numr11  = "11"          // numeric value for ISO 8859-11
	Numr63  = "1963"        // numeric value for ASA X3.4 1963
	Numr65  = "1965"        // numeric value for ASA X3.4 1965
	Numr67  = "1967"        // numeric value for ANSI X3.4 1967/77/86
	Alias11 = "iso885911"   // alias for ISO 8859-11
	Alias67 = "ansi"        // alias for ANSI X3.4 1967/77/86
)

var ErrName = errors.New("there is no encoding name")

// Encoding is an implementation of the Encoding interface that adds a formal name
// to a custom encoding.
type Encoding struct {
	encoding.Encoding // Encoding is the underlying encoding.

	Name string // Name is the formal name of the character encoding.
}

var (
	// XUserDefinedISO11 ISO-8859-11.
	XUserDefinedISO11 encoding.Encoding = &xThaiISO11
	// XUserDefined1963 ASA X3.4 1963.
	XUserDefined1963 encoding.Encoding = &x34_1963
	// XUserDefined1965 ASA X3.4 1965.
	XUserDefined1965 encoding.Encoding = &x34_1965
	// XUserDefined1967 ANSI X3.4 1967/77/86.
	XUserDefined1967 encoding.Encoding = &x34_1967

	xThaiISO11 = Encoding{
		Encoding: charmap.Windows874,
		Name:     "ISO-8859-11",
	}

	x34_1963 = Encoding{
		Encoding: charmap.Windows1252,
		Name:     "ASA X3.4 1963",
	}
	x34_1965 = Encoding{
		Encoding: charmap.Windows1252,
		Name:     "ASA X3.4 1965",
	}
	x34_1967 = Encoding{
		Encoding: charmap.Windows1252,
		Name:     "ANSI X3.4 1967/77/86",
	}
)

// String returns the formal name of the ASA encoding.
func (e Encoding) String() string {
	return e.Name
}

// CodePage returns the encoding of the code page name or alias.
func CodePage(name string) encoding.Encoding {
	switch strings.ToLower(name) {
	case Name11, Numr11, Alias11:
		return XUserDefinedISO11
	case Name63, Numr63:
		return XUserDefined1963
	case Name65, Numr65:
		return XUserDefined1965
	case Name67, Numr67, Alias67:
		return XUserDefined1967
	default:
		return nil
	}
}

// Code7bit reports whether the encoding is a 7-bit ASCII encoding.
// The 7-bit encodings are limited to 127 characters.
// The more common 8-bit encodings are limited to 256 characters.
func Code7bit(e encoding.Encoding) bool {
	switch e {
	case XUserDefined1963, XUserDefined1965, XUserDefined1967:
		return true
	}
	return false
}

// Name returns a named value for the legacy ASA ASCII character encodings.
func Name(e encoding.Encoding) string {
	switch e {
	case XUserDefinedISO11:
		return Name11
	case XUserDefined1963:
		return Name63
	case XUserDefined1965:
		return Name65
	case XUserDefined1967:
		return Name67
	}
	return ""
}

// Numeric returns a numeric value for the legacy ASA ASCII character encodings.
func Numeric(e encoding.Encoding) string {
	switch e {
	case XUserDefinedISO11:
		return Numr11
	case XUserDefined1963:
		return Numr63
	case XUserDefined1965:
		return Numr65
	case XUserDefined1967:
		return Numr67
	}
	return ""
}

// Alias returns an alias value for the legacy ASA ASCII character encodings.
func Alias(e encoding.Encoding) string {
	switch e {
	case XUserDefinedISO11:
		return Alias11
	case XUserDefined1967:
		return Alias67
	}
	return ""
}

// Footnote returns a footnote value for the legacy ASA ASCII character encodings.
func Footnote(w io.Writer, e encoding.Encoding) {
	if w == nil {
		w = io.Discard
	}
	switch e {
	case XUserDefined1963:
		fmt.Fprintln(w)
		fmt.Fprintln(w, "* ASA X3.4 1963 has a number of historic control codes in"+
			"\n  rows 0 and 1 that are not printable in Unicode.")
	case XUserDefined1965:
		fmt.Fprintln(w)
		fmt.Fprintln(w, "* ASA X3.4 1965 cell 1-A is SUB, but it is not printable in Unicode.")
	}
}

// Char returns a string for the 8-bit, character encoding decimal code.
// If the code is not defined in the encoding, then a space is returned.
// If the code matches an existing Windows-1252 character, then -1 is returned.
func Char(e encoding.Encoding, code int) rune {
	switch e {
	case XUserDefinedISO11:
		return CharISO885911(code)
	case XUserDefined1963:
		return CharX3463(code)
	case XUserDefined1965:
		return CharX3465(code)
	case XUserDefined1967:
		return CharX3467(code)
	}
	return -1
}

// CharISO885911 returns a rune for the ISO-8859-11 character code.
// If the code is not defined in the encoding, then a space is returned.
// If the code matches an existing Windows-1252 character, then -1 is returned.
func CharISO885911(code int) rune {
	const pad, nbsp = 128, 160
	if code >= pad && code < nbsp {
		return ' '
	}
	return -1
}

// CharX3463 returns a rune for the legacy ASA X3.4 1963 character code.
// If the code is not defined in the encoding, then a space is returned.
// If the code matches an existing Windows-1252 character, then -1 is returned.
func CharX3463(code int) rune {
	const blank = ' '
	const us, end = 31, 128
	if code >= end || code == 125 {
		return blank
	}
	if x := mapX3493(code); x > 0 {
		return x
	}
	if code <= us {
		return blank
	}
	if code >= 96 && code <= 123 {
		return blank
	}
	return rune(code)
}

func mapX3493(i int) rune {
	m := map[int]rune{
		0:   '␀',
		4:   '␄',
		7:   '␇',
		9:   '␉',
		10:  '␊',
		11:  '␋',
		12:  '␌',
		13:  '␍',
		14:  '␎',
		15:  '␏',
		17:  '␑',
		18:  '␒',
		19:  '␓',
		20:  '␔',
		94:  '↑',
		95:  '←',
		124: '␆',
		126: '␛',
		127: '␡',
	}
	return m[i]
}

// CharX3465 returns a string for the legacy ASA X3.4 1965 character code.
// If the code is not defined in the encoding, then a space is returned.
// If the code matches an existing Windows-1252 character, then -1 is returned.
func CharX3465(code int) rune {
	const sub, grave, tilde, at, not, bar, end = 26, 64, 92, 96, 124, 126, 128
	if code >= end {
		return ' '
	}
	switch code {
	case sub:
		return ' '
	case grave:
		return '`'
	case tilde:
		return '~'
	case at:
		return '@'
	case not:
		return '¬'
	case bar:
		return '|'
	}
	return -1
}

// CharX3467 returns a string for the legacy ASA X3.4 1967 character code.
// If the code is not defined in the encoding, then a space is returned.
// If the code matches an existing Windows-1252 character, then -1 is returned.
func CharX3467(code int) rune {
	const end = 128
	if code >= end {
		return ' '
	}
	return -1
}
