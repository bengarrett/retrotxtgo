package convert

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"text/tabwriter"
	"unicode"

	"github.com/gookit/color"
	"golang.org/x/text/encoding/charmap"
	uni "golang.org/x/text/encoding/unicode"
	"retrotxt.com/retrotxt/lib/str"
)

// Table prints out all the characters in the named 8-bit character set.
func Table(name string) (*bytes.Buffer, error) {
	cp, err := Encoding(name)
	if err != nil {
		return nil, fmt.Errorf("table encoding error: %w", err)
	}
	h := fmt.Sprintf("%s", cp)
	switch cp {
	case uni.UTF16(uni.BigEndian, uni.UseBOM),
		uni.UTF16(uni.BigEndian, uni.IgnoreBOM),
		uni.UTF16(uni.LittleEndian, uni.IgnoreBOM):
		return nil, fmt.Errorf("utf-16 table encodings are not supported")
	case charmap.CodePage037, charmap.CodePage1047, charmap.CodePage1140:
		h += " - EBCDIC"
	case uni.UTF8, uni.UTF8BOM:
		h += " - Unicode"
	default:
		h += " - Extended ASCII"
	}
	var buf bytes.Buffer
	w := new(tabwriter.Writer).Init(&buf, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, " "+color.OpFuzzy.Sprint(strings.Repeat("\u2015", 67)))
	fmt.Fprintln(w, color.Primary.Sprint(str.Center(h, 67)))
	const start, end, max = 0, 15, 255
	for i := 0; i < 16; i++ {
		switch {
		case i == start:
			fmt.Fprintf(w, "%s", color.OpFuzzy.Sprintf("     %X  ", i))
		case i == end:
			fmt.Fprintf(w, "%s", color.OpFuzzy.Sprintf(" %X  \n", i))
		default:
			fmt.Fprintf(w, "%s", color.OpFuzzy.Sprintf(" %X  ", i))
		}
	}
	var conv = Args{Encoding: name}
	var b, row = MakeBytes(), 0
	runes, err := conv.Chars(&b)
	if err != nil {
		return nil, fmt.Errorf("table convert bytes error: %w", err)
	}
	for i, r := range runes {
		char := string(r)
		// non-spacing mark characters require an additional space
		if unicode.In(r, unicode.Mn) {
			char = fmt.Sprintf(" %s", string(r))
		}
		// format, other
		if unicode.In(r, unicode.Cf) {
			const ZWNJ, ZWJ, LRM, RLM = 8204, 8205, 8206, 8207
			switch r {
			case ZWNJ, ZWJ, LRM, RLM:
				// no suitable control character symbols exist
				char = " "
			}
		}
		// Latin-1 Supplement
		if cp == uni.UTF8 || cp == uni.UTF8BOM {
			const PAD, NBSP = 128, 160
			switch {
			case i >= PAD && i < NBSP:
				char = " "
			case i >= NBSP:
				char = string(i)
			}
		}
		switch {
		case i == 0:
			fmt.Fprintf(w, " %s %s %s %s",
				color.OpFuzzy.Sprint("0"),
				color.OpFuzzy.Sprint("|"),
				char, color.OpFuzzy.Sprint("|"))
		case i == max:
			fmt.Fprintf(w, " %s %s\n", char,
				color.OpFuzzy.Sprint("|"))
		case math.Mod(float64(i+1), 16) == 0:
			// every 16th loop
			row++
			fmt.Fprintf(w, " %s %s\n %s %s", char,
				color.OpFuzzy.Sprint("|"),
				color.OpFuzzy.Sprintf("%X", row),
				color.OpFuzzy.Sprint("|"))
		default:
			fmt.Fprintf(w, " %s %s", char,
				color.OpFuzzy.Sprint("|"))
		}
	}
	fmt.Fprint(w, "\n")
	if err = w.Flush(); err != nil {
		return nil, fmt.Errorf("table tab writer failed to flush data: %w", err)
	}
	return &buf, nil
}
