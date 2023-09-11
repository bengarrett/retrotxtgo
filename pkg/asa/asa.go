// Package asa provides the 1960s, American Standards Association (ASA) ASCII character encodings.
package asa

import (
	"fmt"
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

const (
	Text63  = "ascii-63"
	Numr63  = "1963"
	Text65  = "ascii-65"
	Numr65  = "1965"
	Text67  = "ascii-67"
	Numr67  = "1967"
	Alias67 = "ansi"
)

// Encoding is an implementation of the Encoding interface that adds the String
// and ID methods to an existing encoding.
type Encoding struct {
	encoding.Encoding
	Name string
}

var (
	// XUserDefined1963 ASA X3.4 1963.
	XUserDefined1963 encoding.Encoding = &x34_1963 //nolint: gochecknoglobals

	// XUserDefined1965 ASA X3.4 1965.
	XUserDefined1965 encoding.Encoding = &x34_1965 //nolint: gochecknoglobals

	// XUserDefined1967 ANSI X3.4 1967/77/86.
	XUserDefined1967 encoding.Encoding = &x34_1967 //nolint: gochecknoglobals

	x34_1963 = Encoding{ //nolint: gochecknoglobals
		Encoding: charmap.Windows1252,
		Name:     "ASA X3.4 1963",
	}
	x34_1965 = Encoding{ //nolint: gochecknoglobals
		Encoding: charmap.Windows1252,
		Name:     "ASA X3.4 1965",
	}
	x34_1967 = Encoding{ //nolint: gochecknoglobals
		Encoding: charmap.Windows1252,
		Name:     "ANSI X3.4 1967/77/86",
	}
)

func (e Encoding) String() string {
	return e.Name
}

// Code7bit returns true if the encoding is a 7-bit ASCII encoding.
// The 7-bit encodings are limited to 127 characters.
// The more common 8-bit encodings are limited to 255 characters.
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
	case XUserDefined1963:
		return Text63
	case XUserDefined1965:
		return Text65
	case XUserDefined1967:
		return Text67
	}
	return ""
}

// Numeric returns a numeric value for the legacy ASA ASCII character encodings.
func Numeric(e encoding.Encoding) string {
	switch e {
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
	if e == XUserDefined1967 {
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

// CharX3463 returns a string for the legacy ASA X3.4 character codes.
// If the code is not defined in the encoding, then a space is returned.
// If the code matches an existing Windows-1252 character, then -1 is returned.
func Char(e encoding.Encoding, code int) rune {
	switch e {
	case XUserDefined1963:
		return CharX3463(code)
	case XUserDefined1965:
		return CharX3465(code)
	case XUserDefined1967:
		return CharX3467(code)
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
