//Package codepage is extends Go's x/text/encoding capability to convert legacy text
// to UTF-8.
package codepage

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"text/tabwriter"

	"github.com/gookit/color"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
)

// Set blah
type Set struct {
	Data []byte
}

var (
	asciiC0 = []string{"\u0000", "\u263A", "\u263B", "\u2665", "\u2666", "\u2663", "\u2660", "\u2022", "\u25D8", "\u25CB", "\u25D9", "\u2642", "\u2640", "\u266A", "\u266B", "\u263C"}
	asciiC1 = []string{"\u25BA", "\u25C4", "\u2195", "\u203C", "\u00B6", "\u00A7", "\u25AC", "\u21A8", "\u2191", "\u2193", "\u2192", "\u2190", "\u221F", "\u2194", "\u25B2", "\u25BC"}
)

// BOM is the UTF-8 byte order mark prefix.
var BOM = func() []byte {
	return []byte{239, 187, 191} // 0xEF,0xBB,0xBF
}

func MakeMap() [256]byte {
	var m [256]byte
	//encoding, _ := ianaindex.IANA.Encoding("cp437")
	for i := 0; i <= 255; i++ {
		// b := []byte{uint8(i)}
		// c := math.Mod(float64(i), 16)
		// t, _ := Transform(b, encoding)
		// if c == 0 {
		// 	fmt.Print("\n")
		// }
		// fmt.Printf("%x %s\t", b, t)
		m[i] = uint8(i)
	}
	return m
}

func Center(text string, width int) string {
	l := len(text)
	w := (width - l) / 2
	if w > 0 {
		return strings.Repeat("\u0020", w) + text
	}
	return text
}

func Table(codepage string) string {
	e, _ := ianaindex.IANA.Encoding("cp437")

	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 0, 8, 0, '\t', 0)

	units := MakeMap()
	fmt.Fprintln(w, " "+color.OpFuzzy.Sprint(strings.Repeat("\u2015", 67)))
	fmt.Fprintln(w, color.Primary.Sprint(Center(fmt.Sprintf("%s", e), 67)))
	//fmt.Fprint(w, "\n")

	for i := 0; i < 16; i++ {
		switch {
		case i == 0:
			fmt.Fprint(w, color.OpFuzzy.Sprintf("     %X  ", i))
		case i == 15:
			fmt.Fprintf(w, color.OpFuzzy.Sprintf(" %X  \n", i))
		default:
			fmt.Fprintf(w, color.OpFuzzy.Sprintf(" %X  ", i))
		}
	}
	row := 0

	for i, m := range units {
		t, _ := Transform([]byte{m}, e)
		t = SwapRecommended(t)
		switch {
		case i == 0:
			fmt.Fprintf(w, " %s %s %s %s", color.OpFuzzy.Sprint("0"), color.OpFuzzy.Sprint("|"), t, color.OpFuzzy.Sprint("|"))
		case i == 255:
			fmt.Fprintf(w, " %s %s\n", t, color.OpFuzzy.Sprint("|"))
		case math.Mod(float64(i+1), 16) == 0: // on every 16th loop
			row++
			fmt.Fprintf(w, " %s %s\n %s %s", t, color.OpFuzzy.Sprint("|"), color.OpFuzzy.Sprintf("%X", row), color.OpFuzzy.Sprint("|"))
		default:
			fmt.Fprintf(w, " %s %s", t, color.OpFuzzy.Sprint("|"))
		}
	}
	fmt.Fprint(w, "\n")
	w.Flush()
	return buf.String()
}

// ToBOM adds a UTF-8 byte order mark if it doesn't already exist.
func ToBOM(b []byte) []byte {
	if len(b) > 2 {
		if t := b[:3]; bytes.Equal(t, BOM()) == true {
			return b
		}
	}
	return append(BOM(), b...)
}

// UTF8 determines if a document is encoded as UTF-8.
func UTF8(b []byte) bool {
	_, name, _ := charset.DetermineEncoding(b, "text/plain")
	if name == "utf-8" {
		return true
	}
	return false
}

func SwapAll(b []byte) []byte {
	var s Set
	s.Data = b
	s.SwapAll(true)
	return s.Data
}

func SwapRecommended(b []byte) []byte {
	var s Set
	s.Data = b
	s.SwapAll(false)
	return s.Data
}

func (s *Set) SwapAll(nl bool) {
	s.SwapNuls()
	s.SwapPipes()
	s.SwapDels()
	s.SwapNBSP()
	s.SwapControls(nl)
}

// \u0000 should be swapped for SP \u0000 --nul-as-space (true)
func (s *Set) SwapNuls() {
	s.Data = bytes.ReplaceAll(s.Data, []byte{0}, []byte("\u0020"))
}

// \u007c (7C) [pipe] can be swapped for broken bar \u00A6 --pipe-as-broken-bar (false)
func (s *Set) SwapPipes() {
	s.Data = bytes.ReplaceAll(s.Data, []byte{124}, []byte("\u00A6"))
}

// // \u0127? (7F) [delete] can be swapped for a house \u2303 --del-as-house (false)
func (s *Set) SwapDels() {
	s.Data = bytes.ReplaceAll(s.Data, []byte{127}, []byte("\u2303"))
}

// // FF NBSP often displays a ?, it can be replaced with SP --nbsp-as-space (true)
func (s *Set) SwapNBSP() {
	s.Data = bytes.ReplaceAll(s.Data, []byte{255}, []byte("\u263B"))
}

func (s *Set) SwapControls(nl bool) {
	// notes
	for i, u := range append(asciiC0, asciiC1...) {
		if nl == true {
			switch i {
			case 10, 13:
				continue
			}
		}
		s.Data = bytes.ReplaceAll(s.Data, []byte{uint8(i)}, []byte(u))
	}
}

func Transform(text []byte, e encoding.Encoding) ([]byte, error) {
	b, err := e.NewDecoder().Bytes(text)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Set) Transform(text []byte, e encoding.Encoding) error {
	var err error
	s.Data, err = Transform(text, e)
	if err != nil {
		return err
	}
	return nil
}
