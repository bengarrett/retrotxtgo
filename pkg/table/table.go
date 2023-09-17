// Package table creates a table of all the characters in the named 8-bit character set.
package table

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"text/tabwriter"
	"unicode"

	"github.com/bengarrett/retrotxtgo/pkg/byter"
	"github.com/bengarrett/retrotxtgo/pkg/convert"
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/bengarrett/retrotxtgo/pkg/xud"
	"github.com/gookit/color"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	uni "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

var (
	ErrUTF16 = errors.New("utf-16 table encodings are not supported")
	ErrUTF32 = errors.New("utf-32 table encodings are not supported")
)

const width = 68 // width of the table in characters.

// Table prints, aligns and formats to the writer all characters in the named 8-bit character set.
func Table(wr io.Writer, name string) error { //nolint:funlen
	if wr == nil {
		wr = io.Discard
	}
	cp, err := xud.CodePage(name)
	if err != nil {
		return err
	}
	if cp == nil {
		cp, err = CodePage(name)
		if err != nil {
			return err
		}
	}
	h := fmt.Sprintf("%s", cp)
	h += CharmapAlias(cp) + charmapStandard(cp)
	const tabWidth = 8
	w := tabwriter.NewWriter(wr, 0, tabWidth, 0, '\t', 0)
	term.Head(w, width, " "+h)
	columns(w)
	if x := swapper(name); x != nil {
		cp = x
	}
	c := convert.Convert{}
	c.Input.Encoding = cp
	p := byter.MakeBytes()
	runes, err := c.Chars(p...)
	if err != nil {
		return fmt.Errorf("table convert bytes error: %w", err)
	}
	enc := reverter(name)
	const hex, max = 16, 255
	row := 0
out:
	for i, r := range runes {
		char := Character(enc, i, r)
		if x := Replacement(name, i); x != "" {
			char = x
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
		case math.Mod(float64(i+1), hex) == 0:
			// every 16th loop
			row++
			if xud.Code7bit(enc) && row >= 8 {
				// exit after the 7th row
				fmt.Fprintf(w, " %s %s\n", char,
					color.OpFuzzy.Sprint("|"))
				break out
			}
			fmt.Fprintf(w, " %s %s\n %s %s", char,
				color.OpFuzzy.Sprint("|"),
				color.OpFuzzy.Sprintf("%X", row),
				color.OpFuzzy.Sprint("|"))
		default:
			fmt.Fprintf(w, " %s %s", char,
				color.OpFuzzy.Sprint("|"))
		}
	}
	xud.Footnote(w, cp)
	Footnote(w, name)
	fmt.Fprint(w, "\n")
	return w.Flush()
}

func columns(w io.Writer) {
	if w == nil {
		w = io.Discard
	}
	const start, end = 0, 15
	for i := 0; i < 16; i++ {
		switch i {
		case start:
			fmt.Fprintf(w, "%s", color.OpFuzzy.Sprintf("     %X  ", i))
		case end:
			fmt.Fprintf(w, "%s", color.OpFuzzy.Sprintf(" %X  \n", i))
		default:
			fmt.Fprintf(w, "%s", color.OpFuzzy.Sprintf(" %X  ", i))
		}
	}
}

// Footnote writes a footnote with extra details or corrections for the named
// code page or encoding. It is called by Table.
func Footnote(w io.Writer, name string) {
	if w == nil {
		w = io.Discard
	}
	if SHY173(name) {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "* Cell A-D is SHY (soft hyphen), but it is not printable in Unicode.")
		return
	}
	x, _ := CodePage(name)
	if SHY240(x) {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "* Cell F-0 is SHY (soft hyphen), but it is not printable in Unicode.")
	}
}

// CodePage returns the encoding of the code page name or alias.
// But without any of the custom, ASA ASCII or ISO-8859-11 encodings.
func CodePage(s string) (encoding.Encoding, error) {
	cp, err := convert.Encoder(s)
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

// Swapper returns the Windows1252 charmap for use as the base template
// for the ASA ASCII encodings.
func swapper(name string) encoding.Encoding {
	switch strings.ToLower(name) {
	case xud.Name11, xud.Numr11, xud.Alias11:
		return charmap.Windows874
	case xud.Name63, xud.Numr63:
		return charmap.Windows1252
	case xud.Name65, xud.Numr65:
		return charmap.Windows1252
	case xud.Name67, xud.Numr67, xud.Alias67:
		return charmap.Windows1252
	}
	return nil
}

func reverter(name string) encoding.Encoding {
	switch strings.ToLower(name) {
	case xud.Name11, xud.Numr11, xud.Alias11:
		return nil // we don't want to revert to the custom charmap
	case xud.Name63, xud.Numr63:
		return xud.XUserDefined1963
	case xud.Name65, xud.Numr65:
		return xud.XUserDefined1965
	case xud.Name67, xud.Numr67, xud.Alias67:
		return xud.XUserDefined1967
	}
	return nil
}

// SHY240 returns true if the code page has a SHY (soft hyphen) at code 240.
func SHY240(x encoding.Encoding) bool {
	switch x {
	case charmap.CodePage850,
		charmap.CodePage852,
		charmap.CodePage855,
		charmap.CodePage858:
		return true
	}
	return false
}

// SHY173 returns true if the code page has a SHY (soft hyphen) at code 173.
func SHY173(name string) bool {
	s := strings.ToLower(name)
	if strings.Contains(s, "windows 125") ||
		strings.Contains(s, "iso 8859-") ||
		strings.Contains(s, "iso-8859-") {
		return true
	}
	return false
}

// Replacement returns a replacement character for the code page.
func Replacement(name string, code int) string {
	x, _ := CodePage(name)
	switch x {
	case charmap.CodePage850,
		charmap.CodePage858,
		charmap.CodePage865,
		charmap.CodePage437:
		const shy = 152
		if code == shy {
			return "\u00FF"
		}
	}
	if SHY240(x) {
		const shy = 240
		if code == shy {
			return "-"
		}
		return ""
	}
	if SHY173(name) {
		const shy = 173
		if code == shy {
			return "-"
		}
	}
	return ""
}

// Character converts code or rune to an character mapped string.
func Character(cp encoding.Encoding, code int, r rune) string {
	if xud.Name(cp) != "" {
		if x := xud.Char(cp, code); x > -1 {
			return string(x)
		}
		return string(r)
	}
	// non-spacing mark characters
	if unicode.In(r, unicode.Mn) {
		// these require an additional monospace
		return fmt.Sprintf(" %s", string(r))
	}
	// format, other
	if unicode.In(r, unicode.Cf) {
		const zwnj, zwj, lrm, rlm = 8204, 8205, 8206, 8207
		switch r {
		case zwnj, zwj, lrm, rlm:
			return " "
		}
	}
	// unicode latin-1 supplement
	if cp == uni.UTF8 || cp == uni.UTF8BOM {
		const pad, nbsp = 128, 160
		switch {
		case code >= pad && code < nbsp:
			return " "
		case code >= nbsp:
			return string(rune(code))
		}
	}
	return string(r)
}

// CharmapAlias humanizes encodings.
func CharmapAlias(cp encoding.Encoding) string { //nolint:cyclop
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
	case charmap.Windows874, xud.XUserDefinedISO11:
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

// charmapDOS humanizes DOS encodings.
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

// charmapMisc humanizes miscellaneous encodings.
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

// charmapStandard humanizes common encodings.
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
