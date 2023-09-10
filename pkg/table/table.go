package table

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strings"
	"text/tabwriter"
	"unicode"

	"github.com/bengarrett/retrotxtgo/pkg/asa"
	"github.com/bengarrett/retrotxtgo/pkg/byter"
	"github.com/bengarrett/retrotxtgo/pkg/convert"
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/gookit/color"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	uni "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

var (
	ErrNoName = errors.New("there is no encoding name")
	ErrUTF16  = errors.New("utf-16 table encodings are not supported")
	ErrUTF32  = errors.New("utf-32 table encodings are not supported")
)

const (
	width = 67
)

// Table prints out all the characters in the named 8-bit character set.
func Table(name string) (*bytes.Buffer, error) { //nolint:funlen
	cp, err := CodePager(name)
	if err != nil {
		return nil, err
	}
	h := fmt.Sprintf("%s", cp)
	if ISO11(name) {
		h = "ISO 8859-11"
	}
	h += CharmapAlias(cp) + charmapStandard(cp)
	var buf bytes.Buffer
	const tabWidth = 8
	w := new(tabwriter.Writer).Init(&buf, 0, tabWidth, 0, '\t', 0)
	fmt.Fprint(w, " "+term.HeadDark(width, h))
	const start, end, max = 0, 15, 255
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
	b, conv, row := byter.MakeBytes(), converter(name, cp), 0
	runes, err := conv.Chars(b...)
	if err != nil {
		return nil, fmt.Errorf("table convert bytes error: %w", err)
	}
	cp = revert(name)
	const hex = 16
	for i, r := range runes {
		char := Character(cp, i, r)
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
	if err := w.Flush(); err != nil {
		return nil, fmt.Errorf("table tab writer failed to flush data: %w", err)
	}
	return &buf, nil
}

// ISO11 returns true if s matches an ISO-8859-11 name or alias.
func ISO11(s string) bool {
	switch strings.ToUpper(s) {
	case
		"ISO 8859-11",
		"ISO-8859-11",
		"ISO8859-11",
		"11",
		"ISO885911":
		return true
	}
	return false
}

// CodePager returns the encoding of the code page name or alias.
func CodePager(s string) (encoding.Encoding, error) { //nolint:ireturn
	if s == "" {
		return nil, ErrNoName
	}
	if ISO11(s) {
		return charmap.Windows874, nil
	}
	switch strings.ToLower(s) {
	case asa.Ascii63:
		return asa.ASAX34_1963, nil
	case asa.Ascii65:
		return asa.ASAX34_1965, nil
	case asa.Ascii67:
		return asa.ANSIX34_1967, nil
	default:
		return CodePage(s)
	}
}

// CodePage returns the encoding of the code page name or alias.
// But without any of the custom, ASA ASCII or ISO-8859-11 encodings.
func CodePage(s string) (encoding.Encoding, error) { //nolint:ireturn
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

func converter(name string, cp encoding.Encoding) convert.Convert {
	conv := convert.Convert{}
	switch strings.ToLower(name) {
	case asa.Ascii63:
		cp = asa.ASAX34_1963
	case asa.Ascii65:
		cp = asa.ASAX34_1965
	case asa.Ascii67:
		cp = asa.ANSIX34_1967
	}
	conv.Input.Encoding = cp
	return conv
}

func revert(name string) encoding.Encoding { //nolint:ireturn
	if ISO11(name) {
		return charmap.XUserDefined
	}
	switch strings.ToLower(name) {
	case asa.Ascii63:
		return asa.ASAX34_1963
	case asa.Ascii65:
		return asa.ASAX34_1965
	case asa.Ascii67:
		return asa.ANSIX34_1967
	}
	return nil
}

// CharISO11 returns a string for the ISO-8859-11 character codes.
func CharISO11(cp encoding.Encoding, code int) rune {
	// ISO-8859-11 is not included in Go so a user defined charmap is used.
	ISO8859_11 := charmap.XUserDefined
	if cp != ISO8859_11 {
		return -1
	}
	const pad, nbsp = 128, 160
	if code >= pad && code < nbsp {
		return ' '
	}
	return -1
}

// Character converts code or rune to an character mapped string.
func Character(cp encoding.Encoding, code int, r rune) string {
	if r := asa.Char(cp, code); r > -1 {
		return string(r)
	}
	if r := CharISO11(cp, code); r > -1 {
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
