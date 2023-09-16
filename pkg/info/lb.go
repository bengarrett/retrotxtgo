package info

import (
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/pkg/fsys"
)

// LineBreak contains details on the line break sequence used to create a new line in a text file.
type LineBreak struct {
	Abbr    string  `json:"string"   xml:"string,attr"` // Abbr is the abbreviation for the line break.
	Escape  string  `json:"escape"   xml:"-"`           // Escape is the escape sequence for the line break.
	Decimal [2]rune `json:"decimal" xml:"decimal"`      // Decimal is the numeric character code for the line break.
}

// Find determines the new lines characters found in the rune pair.
func (lb *LineBreak) Find(r [2]rune) {
	a, e := "", ""
	switch r {
	case [2]rune{lf}:
		a = "lf"
		e = "\n"
	case [2]rune{cr}:
		a = "cr"
		e = "\r"
	case [2]rune{cr, lf}:
		a = "crlf"
		e = "\r\n"
	case [2]rune{lf, cr}:
		a = "lfcr"
		e = "\n\r"
	case [2]rune{nl}, [2]rune{nel}:
		a = "nl"
		e = "\025"
	}
	lb.Decimal = r
	lb.Abbr = strings.ToUpper(a)
	lb.Escape = e
}

// Total counts the number of lines in the named file
// based on the line break sequence.
func (lb *LineBreak) Total(name string) (int, error) {
	f, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	l, err := fsys.Lines(f, lb.Decimal)
	if err != nil {
		return 0, err
	}
	return l, f.Close()
}
