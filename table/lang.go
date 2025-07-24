package table

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/term"
	"github.com/bengarrett/retrotxtgo/xud"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

// Lang describes the common natural language uses of the encoding.
type Lang map[encoding.Encoding]string

// Languages returns a list of code page encodings and their target natural languages.
func Languages() *Lang {
	lang := Lang{
		charmap.CodePage037:   "US English",
		charmap.CodePage437:   "US English",
		charmap.CodePage850:   "West Europe",
		charmap.CodePage852:   "Central Europe Latin script",
		charmap.CodePage855:   "Central Europe Cyrillic script",
		charmap.CodePage858:   "West Europe, includes the € symbol",
		charmap.CodePage860:   "Portuguese",
		charmap.CodePage862:   "Hebrew",
		charmap.CodePage863:   "French Canadian",
		charmap.CodePage865:   "Danish, Norwegian",
		charmap.CodePage866:   "USSR Cyrillic script",
		charmap.CodePage1047:  "West Europe",
		charmap.CodePage1140:  "US English",
		charmap.ISO8859_1:     "West Europe",
		charmap.ISO8859_2:     "Central Europe Latin script",
		charmap.ISO8859_3:     "Esperanto, Maltese, Turkish",
		charmap.ISO8859_4:     "Estonian, Latvian, Lithuanian, Greenlandic, Sámi",
		charmap.ISO8859_5:     "Russian Cyrillic script",
		charmap.ISO8859_6:     "Arabic",
		charmap.ISO8859_6E:    "Arabic",
		charmap.ISO8859_6I:    "Arabic",
		charmap.ISO8859_7:     "Greek",
		charmap.ISO8859_8:     "Hebrew",
		charmap.ISO8859_8E:    "Hebrew",
		charmap.ISO8859_8I:    "Hebrew",
		charmap.ISO8859_9:     "Turkish",
		charmap.ISO8859_10:    "Nordic languages",
		xud.XUserDefinedISO11: "Thai", // ISO-8859-11
		charmap.ISO8859_13:    "Baltic languages",
		charmap.ISO8859_14:    "Celtic languages",
		charmap.ISO8859_15:    "West Europe, includes the € symbol",
		charmap.ISO8859_16:    "Gaj's Latin alphabet for European languages",
		charmap.KOI8R:         "Russian, Bulgarian",
		charmap.KOI8U:         "Ukrainian",
		charmap.Macintosh:     "West Europe",
		charmap.Windows874:    "Thai",
		charmap.Windows1250:   "Central Europe Latin script",
		charmap.Windows1251:   "Cyrillic script",
		charmap.Windows1252:   "English and West Europe",
		charmap.Windows1253:   "Greek",
		charmap.Windows1254:   "Turkish",
		charmap.Windows1255:   "Hebrew",
		charmap.Windows1256:   "Arabic",
		charmap.Windows1257:   "Estonian, Latvian, Lithuanian",
		charmap.Windows1258:   "Vietnamese",
		japanese.ShiftJIS:     "Japanese",
		unicode.UTF8:          "Unicode, all major languages",
		xud.XUserDefined1963:  "US English",
		xud.XUserDefined1965:  "US English",
		xud.XUserDefined1967:  "US English",
	}
	return &lang
}

// Language returns the natural language usage of the encoding.
func Language(e encoding.Encoding) string {
	l := *Languages()
	return l[e]
}

// ListLanguage writes a tabled list of supported IANA character set encodings
// and the languages they target.
func ListLanguage(wr io.Writer) error {
	if wr == nil {
		wr = io.Discard
	}
	const header, title = " Formal name\t Named value\t Language\t",
		" Known legacy code pages and their target languages, regions, or alphabets "
	const padding, width = 2, 76
	w := tabwriter.NewWriter(wr, 0, 0, padding, ' ', 0)
	term.Head(wr, width, title)
	fmt.Fprintf(w, "\n%s\n", header)
	x := Charmaps()
	x = append(x,
		xud.XUserDefined1963,
		xud.XUserDefined1965,
		xud.XUserDefined1967)
	for _, e := range x {
		switch e {
		case charmap.XUserDefined:
			continue
		case charmap.ISO8859_10:
			c, err := Rows(e)
			if err != nil {
				return fmt.Errorf("%q: %w", e, err)
			}
			fmt.Fprintf(w, " %s\t %s\t %s\t\n",
				c.Name, c.Value, Language(e))
			// intentionally insert ISO-8895-11 after 10.
			x := xud.XUserDefinedISO11
			fmt.Fprintf(w, " %s\t %s\t %s\t\n",
				x, xud.Name(x), Language(x))
			continue
		case
			unicode.UTF16(unicode.BigEndian, unicode.UseBOM),
			unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
			unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
			utf32.UTF32(utf32.BigEndian, utf32.UseBOM),
			utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM),
			utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM):
			continue
		}
		c, err := Rows(e)
		if err != nil {
			return fmt.Errorf("%q: %w", e, err)
		}
		// do not use ANSI colors in cells as it will break the table layout
		fmt.Fprintf(w, " %s\t %s\t %s\t\n",
			c.Name, c.Value, Language(e))
	}
	return w.Flush()
}
