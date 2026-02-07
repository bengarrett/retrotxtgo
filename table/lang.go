package table

import (
	"fmt"
	"io"

	"github.com/bengarrett/retrotxtgo/xud"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

// Lang describes the common natural language uses of the encoding.
type Lang map[encoding.Encoding]string

const (
	formalName = "Formal name"
	namedVal   = "Named value"
	langRegion = "Language, script, or region"
)

// Languages returns a list of code page encodings and their target natural languages.
// These are displayed in the order listed.
func Languages() *Lang {
	const (
		arb = "Arabic"
		eur = " with the € symbol"
		heb = "Hebrew"
		weu = "Western Europe"
		usa = "English, US"
	)
	lang := Lang{
		unicode.UTF8:            "Unicode, all major languages",
		charmap.CodePage037:     usa,
		charmap.CodePage437:     usa,
		charmap.CodePage850:     weu,
		charmap.CodePage852:     "Central Europe Latin script",
		charmap.CodePage855:     "Central Europe Cyrillic script",
		charmap.CodePage858:     weu + eur,
		charmap.CodePage860:     "Portuguese",
		charmap.CodePage862:     heb,
		charmap.CodePage863:     "French Canadian",
		charmap.CodePage865:     "Danish, Norwegian",
		charmap.CodePage866:     "USSR Cyrillic script",
		charmap.CodePage1047:    weu,
		charmap.CodePage1140:    usa,
		charmap.ISO8859_1:       weu,
		charmap.ISO8859_2:       "Central Europe Latin script",
		charmap.ISO8859_3:       "Esperanto, Maltese, Turkish",
		charmap.ISO8859_4:       "Estonian, Latvian, Lithuanian, Greenlandic, Sámi",
		charmap.ISO8859_5:       "Russian Cyrillic script",
		charmap.ISO8859_6:       arb,
		charmap.ISO8859_6E:      arb,
		charmap.ISO8859_6I:      arb,
		charmap.ISO8859_7:       "Greek",
		charmap.ISO8859_8:       heb,
		charmap.ISO8859_8E:      heb,
		charmap.ISO8859_8I:      heb,
		charmap.ISO8859_9:       "Turkish",
		charmap.ISO8859_10:      "Nordic languages",
		xud.XUserDefinedISO11:   "Thai", // ISO-8859-11
		charmap.ISO8859_13:      "Baltic languages",
		charmap.ISO8859_14:      "Celtic languages",
		charmap.ISO8859_15:      weu + eur,
		charmap.ISO8859_16:      "Gaj's Latin alphabet for European languages",
		charmap.KOI8R:           "Russian, Bulgarian",
		charmap.KOI8U:           "Ukrainian",
		charmap.Macintosh:       weu,
		charmap.Windows874:      "Thai",
		charmap.Windows1250:     "Central Europe Latin script",
		charmap.Windows1251:     "Cyrillic script",
		charmap.Windows1252:     "English, " + weu,
		charmap.Windows1253:     "Greek",
		charmap.Windows1254:     "Turkish",
		charmap.Windows1255:     heb,
		charmap.Windows1256:     arb,
		charmap.Windows1257:     "Estonian, Latvian, Lithuanian",
		charmap.Windows1258:     "Vietnamese",
		japanese.ShiftJIS:       "Japanese",
		traditionalchinese.Big5: "Traditional Chinese",
		xud.XUserDefined1963:    usa,
		xud.XUserDefined1965:    usa,
		xud.XUserDefined1967:    usa,
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

	// Create rows for the language table
	rows := make([]LanguageRow, 0)

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
			rows = append(rows, LanguageRow{
				Name:     c.Name,
				Value:    c.Value,
				Language: Language(e),
			})
			// intentionally insert ISO-8895-11 after 10.
			x := xud.XUserDefinedISO11
			rows = append(rows, LanguageRow{
				Name:     xud.Name(x),
				Value:    xud.Name(x),
				Language: Language(x),
			})
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
		rows = append(rows, LanguageRow{
			Name:     c.Name,
			Value:    c.Value,
			Language: Language(e),
		})
	}

	// Use lipgloss table rendering
	return LipglossLanguageTable(wr, rows)
}

// LanguageRow represents a row in the language table.
type LanguageRow struct {
	Name     string
	Value    string
	Language string
}

// LipglossLanguageTable renders a language table using lipgloss styling.
func LipglossLanguageTable(wr io.Writer, rows []LanguageRow) error {
	if wr == nil {
		wr = io.Discard
	}

	// Create lipgloss styles
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("231")).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	// Calculate column widths
	colWidths := calculateLanguageColumnWidths(rows)

	// Create header
	header := createLanguageHeader(headerStyle, colWidths)

	// Create rows
	rowStrings := make([]string, 0, len(rows))
	for _, row := range rows {
		rowString := createLanguageRow(&row, cellStyle, colWidths)
		rowStrings = append(rowStrings, rowString)
	}

	// Build the table
	table := borderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			lipgloss.JoinVertical(lipgloss.Left, rowStrings...),
		),
	)

	// Write the table
	fmt.Fprintln(wr, table)
	return nil
}

func calculateLanguageColumnWidths(rows []LanguageRow) [3]int {
	var widths [3]int

	// Header widths
	headers := []string{formalName, namedVal, langRegion}
	for i, header := range headers {
		if len(header) > widths[i] {
			widths[i] = len(header)
		}
	}

	// Data widths
	for _, row := range rows {
		if len(row.Name) > widths[0] {
			widths[0] = len(row.Name)
		}
		if len(row.Value) > widths[1] {
			widths[1] = len(row.Value)
		}
		if len(row.Language) > widths[2] {
			widths[2] = len(row.Language)
		}
	}

	// Add some padding
	for i := range widths {
		widths[i] += 2
	}

	return widths
}

func createLanguageHeader(style lipgloss.Style, widths [3]int) string {
	headerCells := []string{
		style.Render(FitString(formalName, widths[0])),
		style.Render(FitString(namedVal, widths[1])),
		style.Render(FitString(langRegion, widths[2])),
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, headerCells...)
}

func createLanguageRow(row *LanguageRow, cellStyle lipgloss.Style, widths [3]int) string {
	cells := []string{
		cellStyle.Render(FitString(row.Name, widths[0])),
		cellStyle.Render(FitString(row.Value, widths[1])),
		cellStyle.Render(FitString(row.Language, widths[2])),
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, cells...)
}
