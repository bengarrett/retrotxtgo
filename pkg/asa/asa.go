// Package asa provides the 1960s, American Standards Association (ASA) ASCII character encodings.
package asa

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

const (
	Ascii63 = "ascii-63"
	Ascii65 = "ascii-65"
	Ascii67 = "ascii-67"
)

// Encoding is an implementation of the Encoding interface that adds the String
// and ID methods to an existing encoding.
type Encoding struct {
	encoding.Encoding
	Name string
}

var (
	// ASAX34_1963 ASA X3.4 1963.
	ASAX34_1963 encoding.Encoding = &x34_1963 //nolint: gochecknoglobals

	// AsaX34_1965 ASA X3.4 1965.
	ASAX34_1965 encoding.Encoding = &x34_1965 //nolint: gochecknoglobals

	// AnsiX34_1967 ANSI X3.4 1967/77/86.
	ANSIX34_1967 encoding.Encoding = &x34_1967 //nolint: gochecknoglobals

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

// Name returns a named value for the legacy ASA ASCII character encodings.
func Name(e encoding.Encoding) string {
	switch e {
	case ASAX34_1963:
		return Ascii63
	case ASAX34_1965:
		return Ascii65
	case ANSIX34_1967:
		return Ascii67
	}
	return ""
}

// CharX3463 returns a string for the legacy ASA X3.4 character codes.
// If the code is not defined in the encoding, then a space is returned.
// If the code matches an existing Windows-1252 character, then -1 is returned.
func Char(e encoding.Encoding, code int) rune {
	switch e {
	case ASAX34_1963:
		return CharX3463(code)
	case ASAX34_1965:
		return CharX3465(code)
	case ANSIX34_1967:
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
