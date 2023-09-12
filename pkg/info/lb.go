package info

import (
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/pkg/fsys"
)

// LineBreaks for new line toggles.
type LineBreaks struct {
	Abbr     string  `json:"string"   xml:"string,attr"`
	Escape   string  `json:"escape"   xml:"-"`
	Decimals [2]rune `json:"decimals" xml:"decimal"`
}

// Find determines the new lines characters found in the rune pair.
func (lb *LineBreaks) Find(r [2]rune) {
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
	lb.Decimals = r
	lb.Abbr = strings.ToUpper(a)
	lb.Escape = e
}

// Total counts the totals lines in the named file.
func (lb *LineBreaks) Total(name string) (int, error) {
	f, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	l, err := fsys.Lines(f, lb.Decimals)
	if err != nil {
		return 0, err
	}
	return l, f.Close()
}
