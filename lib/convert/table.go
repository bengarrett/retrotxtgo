package convert

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"text/tabwriter"
	"unicode"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/gookit/color"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	uni "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

const width = 67

// Table prints out all the characters in the named 8-bit character set.
func Table(name string) (*bytes.Buffer, error) {
	cp, err := tblEncode(name)
	if err != nil {
		return nil, err
	}
	h := fmt.Sprintf("%s", cp)
	if a := encodingAlias(shorten(name)); a == "iso-8859-11" {
		h = "ISO 8859-11"
		cp = charmap.XUserDefined
	}
	h += charmapAlias(cp)
	h += charmapStandard(cp)
	var buf bytes.Buffer
	const tabWidth, horizontalBar = 8, "\u2015"
	w := new(tabwriter.Writer).Init(&buf, 0, tabWidth, 0, '\t', 0)
	fmt.Fprintln(w, " "+color.OpFuzzy.Sprint(strings.Repeat(horizontalBar, width)))
	fmt.Fprintln(w, color.Primary.Sprint(str.Center(width, h)))
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
	var conv = Convert{}
	conv.Source.E = cp
	var b, row = MakeBytes(), 0
	runes, err := conv.Chars(&b)
	if err != nil {
		return nil, fmt.Errorf("table convert bytes error: %w", err)
	}
	const hex = 16
	for i, r := range runes {
		char := character(i, r, cp)
		switch {
		case i == 0:
			fmt.Fprintf(w, " %s %s %s %s",
				color.OpFuzzy.Sprint("0"),
				color.OpFuzzy.Sprint("|"),
				char, color.OpFuzzy.Sprint("|"))
		case i == max:
			fmt.Fprintf(w, " %s %s\n", char,
				color.OpFuzzy.Sprint("|"))
		case math.Mod(float64(i+1), hex) == 0:
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

func tblEncode(name string) (encoding.Encoding, error) {
	cp, err := Encoding(name)
	if err != nil {
		return nil, fmt.Errorf("table encoding error: %w", err)
	}
	switch cp {
	case uni.UTF16(uni.BigEndian, uni.UseBOM),
		uni.UTF16(uni.BigEndian, uni.IgnoreBOM),
		uni.UTF16(uni.LittleEndian, uni.IgnoreBOM):
		return nil, ErrUTF16
	case utf32.UTF32(utf32.BigEndian, utf32.UseBOM),
		utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM),
		utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM):
		return nil, ErrUTF32
	}
	return cp, nil
}

// Character converts rune to an encoded string.
func character(i int, r rune, cp encoding.Encoding) string {
	// ISO-8859-11 is not included in Go so a user defined charmap is used.
	var iso8859_11 = charmap.XUserDefined
	if cp == iso8859_11 {
		const PAD, NBSP = 128, 160
		if i >= PAD && i < NBSP {
			return " "
		}
	}
	// non-spacing mark characters
	if unicode.In(r, unicode.Mn) {
		// these require an additional monospace
		return fmt.Sprintf(" %s", string(r))
	}
	// format, other
	if unicode.In(r, unicode.Cf) {
		const ZWNJ, ZWJ, LRM, RLM = 8204, 8205, 8206, 8207
		switch r {
		case ZWNJ, ZWJ, LRM, RLM:
			// no suitable control character symbols exist
			return " "
		}
	}
	// unicode latin-1 supplement
	if cp == uni.UTF8 || cp == uni.UTF8BOM {
		const PAD, NBSP = 128, 160
		switch {
		case i >= PAD && i < NBSP:
			return " "
		case i >= NBSP:
			return string(rune(i))
		}
	}
	// rune to string
	return string(r)
}

// CharmapAlias humanizes ISO encodings.
func charmapAlias(cp encoding.Encoding) string { // nolint:gocyclo
	if c := charmapDOS(cp); c != "" {
		return c
	}
	switch cp {
	case charmap.CodePage1047:
		return " (C programming language)"
	case charmap.CodePage1140:
		return " (US/Canada Latin 1 plus €)"
	case charmap.ISO8859_1, charmap.Windows1252:
		return " (Western European)"
	case charmap.ISO8859_2, charmap.Windows1250:
		return " (Central European)"
	case charmap.ISO8859_3:
		return " (South European)"
	case charmap.ISO8859_4:
		return " (North European)"
	case charmap.ISO8859_5, charmap.Windows1251:
		return " (Cyrillic)"
	case charmap.ISO8859_6, charmap.Windows1256:
		return " (Arabic)"
	case charmap.ISO8859_7, charmap.Windows1253:
		return " (Greek)"
	case charmap.ISO8859_8, charmap.Windows1255:
		return " (Hebrew)"
	case charmap.ISO8859_9, charmap.Windows1254:
		return " (Turkish)"
	case charmap.ISO8859_10:
		return " (Nordic)"
	case charmap.Windows874, charmap.XUserDefined:
		return " (Thai)"
	case charmap.ISO8859_13, charmap.Windows1257:
		return " (Baltic Rim)"
	case charmap.ISO8859_14:
		return " (Celtic)"
	case charmap.ISO8859_15:
		return " (Western European, 1999)"
	case charmap.ISO8859_16:
		return " (South-Eastern European)"
	}
	if c := charmapMisc(cp); c != "" {
		return c
	}
	return ""
}

// CharmapDOS humanizes DOS encodings.
func charmapDOS(cp encoding.Encoding) string {
	switch cp {
	case charmap.CodePage037:
		return " (US/Canada Latin 1)"
	case charmap.CodePage437:
		return " (DOS, OEM-US)"
	case charmap.CodePage850:
		return " (DOS, Latin 1)"
	case charmap.CodePage852:
		return " (DOS, Latin 2)"
	case charmap.CodePage855:
		return " (DOS, Cyrillic)"
	case charmap.CodePage858:
		return " (DOS, Western Europe)"
	case charmap.CodePage860:
		return " (DOS, Portuguese)"
	case charmap.CodePage862:
		return " (DOS, Hebrew)"
	case charmap.CodePage863:
		return " (DOS, French Canada)"
	case charmap.CodePage865:
		return " (DOS, Nordic)"
	case charmap.CodePage866:
		return " (DOS, Cyrillic Russian)"
	}
	return ""
}

// CharmapDOS humanizes encodings.
func charmapMisc(cp encoding.Encoding) string {
	switch cp {
	case charmap.KOI8R:
		return " (Russian)"
	case charmap.KOI8U:
		return " (Ukrainian)"
	case charmap.Macintosh:
		return " (Mac OS Roman)"
	case charmap.Windows1258:
		return " (Vietnamese)"
	case japanese.ShiftJIS:
		return " (Japanese)"
	}
	return ""
}

// CharmapDOS humanizes common encodings.
func charmapStandard(cp encoding.Encoding) string {
	switch cp {
	case charmap.CodePage037, charmap.CodePage1047, charmap.CodePage1140:
		return " - EBCDIC"
	case uni.UTF8, uni.UTF8BOM:
		return " - Unicode"
	default:
		return " - Extended ASCII"
	}
}
