package convert

import (
	"bytes"
	"errors"
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

const (
	width   = 67
	ascii63 = "ascii-63"
	ascii65 = "ascii-65"
	ascii67 = "ascii-67"
)

var (
	ErrNoName = errors.New("there is no encoding name")
)

func ISO11Name(name string) bool {
	switch strings.ToUpper(name) {
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

// Table prints out all the characters in the named 8-bit character set.
func Table(name string) (*bytes.Buffer, error) { //nolint:funlen
	cp, err := codepager(name)
	if err != nil {
		return nil, err
	}
	h := fmt.Sprintf("%s", cp)
	if ISO11Name(name) {
		h = "ISO 8859-11"
	}
	h += charmapAlias(cp) + charmapStandard(cp)
	var buf bytes.Buffer
	const tabWidth = 8
	w := new(tabwriter.Writer).Init(&buf, 0, tabWidth, 0, '\t', 0)
	fmt.Fprint(w, " "+str.HeadDark(width, h))
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
	b, conv, row := MakeBytes(), encoder(name, cp), 0
	runes, err := conv.Chars(b...)
	if err != nil {
		return nil, fmt.Errorf("table convert bytes error: %w", err)
	}
	cp = revert(name)
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

func codepager(name string) (encoding.Encoding, error) {
	if name == "" {
		return nil, ErrNoName
	}
	if ISO11Name(name) {
		return charmap.Windows874, nil
	}
	switch strings.ToLower(name) {
	case ascii63:
		return AsaX34_1963, nil
	case ascii65:
		return AsaX34_1965, nil
	case ascii67:
		return x34_1967, nil
	default:
		cp, err := defaultCP(name)
		if err != nil {
			return nil, err
		}
		return cp, nil
	}
}

func defaultCP(name string) (encoding.Encoding, error) {
	cp, err := Encoder(name)
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

func encoder(name string, cp encoding.Encoding) Convert {
	conv := Convert{}
	switch strings.ToLower(name) {
	case ascii63, ascii65, ascii67:
		cp = charmap.Windows1252
	}
	conv.Input.Encoding = cp
	return conv
}

func revert(name string) encoding.Encoding {
	if ISO11Name(name) {
		return charmap.XUserDefined
	}
	switch strings.ToLower(name) {
	case ascii63:
		return AsaX34_1963
	case ascii65:
		return AsaX34_1965
	case ascii67:
		return AnsiX34_1967
	}
	return nil
}

func chrX3493(pos int, cp encoding.Encoding) string {
	if cp != AsaX34_1963 {
		return ""
	}
	const us, end = 31, 128
	if pos >= end || pos == 125 {
		return " "
	}
	if x := mapX3493(pos); x != "" {
		return x
	}
	if pos <= us {
		return " "
	}
	if pos >= 96 && pos <= 123 {
		return " "
	}
	return ""
}

func mapX3493(i int) string {
	m := map[int]string{
		0:   "␀",
		4:   "␄",
		7:   "␇",
		9:   "␉",
		10:  "␊",
		11:  "␋",
		12:  "␌",
		13:  "␍",
		14:  "␎",
		15:  "␏",
		17:  "␑",
		18:  "␒",
		19:  "␓",
		20:  "␔",
		94:  "↑",
		95:  "←",
		124: "␆",
		126: "␛",
		127: "␡",
	}
	return m[i]
}

func chrX3495(pos int, cp encoding.Encoding) string {
	if cp != AsaX34_1965 {
		return ""
	}
	const sub, grave, tilde, at, not, bar, end = 26, 64, 92, 96, 124, 126, 128
	if pos >= end {
		return " "
	}
	switch pos {
	case sub:
		return " "
	case grave:
		return "`"
	case tilde:
		return "~"
	case at:
		return "@"
	case not:
		return "¬"
	case bar:
		return "|"
	}
	return ""
}

func chrX3497(pos int, cp encoding.Encoding) string {
	if cp != AnsiX34_1967 {
		return ""
	}
	const end = 128
	if pos >= end {
		return " "
	}
	return ""
}

func chrISO11(pos int, cp encoding.Encoding) string {
	// ISO-8859-11 is not included in Go so a user defined charmap is used.
	iso8859_11 := charmap.XUserDefined
	if cp != iso8859_11 {
		return ""
	}
	const pad, nbsp = 128, 160
	if pos >= pad && pos < nbsp {
		return " "
	}
	return ""
}

// character converts rune to an encoded string.
func character(pos int, r rune, cp encoding.Encoding) string {
	if s := chrX3493(pos, cp); s != "" {
		return s
	}
	if s := chrX3495(pos, cp); s != "" {
		return s
	}
	if s := chrX3497(pos, cp); s != "" {
		return s
	}
	if s := chrISO11(pos, cp); s != "" {
		return s
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
		case pos >= pad && pos < nbsp:
			return " "
		case pos >= nbsp:
			return string(rune(pos))
		}
	}
	// rune to string
	return string(r)
}

// charmapAlias humanizes encodings.
func charmapAlias(cp encoding.Encoding) string { //nolint:cyclop
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
