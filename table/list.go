package table

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/term"
	"github.com/bengarrett/retrotxtgo/xud"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

const (
	latin            = "isolatin"
	encodingCapacity = 100 // Estimated capacity for all supported encodings
)

// Row is an item for the list of code pages.
type Row struct {
	Name    string // Name is the formal name of the character encoding.
	Value   string // Value is the short name of the character encoding.
	Numeric string // Numeric is an optional, shorter numeric value of the character encoding.
	Alias   string // Alias is an optional, informal but common use value of the character encoding.
}

// Charmaps returns all the supported legacy text encodings.
func Charmaps() []encoding.Encoding {
	// Preallocate with estimated capacity for all encodings.
	e := make([]encoding.Encoding, 0, encodingCapacity)
	// Create a collection of all the encodings.
	a := charmap.All
	a = append(a, japanese.All...)
	a = append(a, traditionalchinese.All...)
	a = append(a, unicode.All...)
	a = append(a, utf32.All...)
	// Iterate the collection and skip the unwanted and duplicate encodings.
	for _, m := range a {
		switch m {
		case japanese.EUCJP,
			japanese.ISO2022JP,
			charmap.MacintoshCyrillic:
			continue
		}
		e = append(e, m)
	}
	return e
}

// List returns a tabled list of supported IANA character set encodings.
func List(wr io.Writer) error {
	return listWithLipgloss(wr)
}

// listWithLipgloss uses lipgloss for modern table formatting.
func listWithLipgloss(wr io.Writer) error {
	rows, err := createTableRows()
	if err != nil {
		return err
	}

	if err := LipglossTable(wr, rows); err != nil {
		return err
	}

	printLegendAndFooter(wr)
	return nil
}

// createTableRows creates the rows for the table.
func createTableRows() ([]Row, error) {
	// Preallocate rows with estimated capacity
	const extraEncodings = 3 // For the extra XUserDefined encodings
	rows := make([]Row, 0, len(Charmaps())+extraEncodings)
	x := Charmaps()
	x = append(x, xud.XUserDefined1963, xud.XUserDefined1965, xud.XUserDefined1967)

	for _, e := range x {
		if e == charmap.XUserDefined {
			continue
		}
		c, err := Rows(e)
		if err != nil {
			return nil, err
		}

		rows = append(rows, processRow(e, c)...)
	}

	return rows, nil
}

// processRow processes a single row based on the encoding type.
func processRow(e encoding.Encoding, c Row) []Row {
	switch e {
	case charmap.ISO8859_10:
		const iso885910Rows = 2 // ISO-8859-10 + ISO-8859-11
		rows := make([]Row, 0, iso885910Rows)
		rows = append(rows, c)
		// intentionally insert ISO-8895-11 after 10.
		x11 := Row{
			Name:    fmt.Sprint(xud.XUserDefinedISO11),
			Value:   xud.Name(xud.XUserDefinedISO11),
			Numeric: xud.Numeric(xud.XUserDefinedISO11),
			Alias:   xud.Alias(xud.XUserDefinedISO11),
		}
		rows = append(rows, x11)
		return rows
	case charmap.CodePage037, charmap.CodePage1047, charmap.CodePage1140:
		c.Name = "* " + c.Name
		return []Row{c}
	case
		traditionalchinese.Big5,
		unicode.UTF16(unicode.BigEndian, unicode.UseBOM),
		unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
		unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
		utf32.UTF32(utf32.BigEndian, utf32.UseBOM),
		utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM),
		utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM):
		c.Name = "† " + c.Name
		return []Row{c}
	case xud.XUserDefined1963, xud.XUserDefined1965, xud.XUserDefined1967:
		c.Name = "⁑ " + c.Name
		return []Row{c}
	}
	return []Row{c}
}

// printLegendAndFooter prints the legend and footer information.
func printLegendAndFooter(wr io.Writer) {
	// Use lipgloss styles to match the table colors
	specialStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	nonTableStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	tableOnlyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("141"))

	fmt.Fprintln(wr, "\n "+specialStyle.Render("Yellow text")+
		" indicates EBCDIC encodings found on IBM mainframes (not ASCII compatible).")
	fmt.Fprintln(wr, " "+nonTableStyle.Render("Black text")+
		" indicates encodings not usable with the "+term.Example("table")+" command.")
	fmt.Fprintln(wr, " "+tableOnlyStyle.Render("Green text")+
		" indicates encodings only usable with the "+term.Example("table")+" command."+
		"\n You can use the \""+term.Example("table ascii")+"\" command to list all three X3.4 tables.")
	fmt.Fprintln(wr, "\nNamed, numeric, or alias values are all valid code page arguments.")
	fmt.Fprintln(wr, "These values all match ISO 8859-1:")
	cmds := meta.Bin + " table "
	fmt.Fprintf(wr, "  %s%s  %s\n",
		term.Example(cmds), term.Comment("iso-8859-1"), term.Fuzzy("# named"))
	fmt.Fprintf(wr, "  %s%s           %s\n",
		term.Example(cmds), term.Comment("1"), term.Fuzzy("# numeric"))
	fmt.Fprintf(wr, "  %s%s      %s\n",
		term.Example(cmds), term.Comment("latin1"), term.Fuzzy("# alias"))
	fmt.Fprintf(wr, "\n  IBM Code Page 437 (%s) is commonly used on MS-DOS and for ANSI art.\n",
		term.Comment("cp437"))
	fmt.Fprintf(wr, "  ISO 8859-1 (%s) is found on historic Unix, Amiga, and the early Internet.\n",
		term.Comment("latin1"))
	fmt.Fprintf(wr, "  Windows 1252 (%s) is found on Windows ME/98 and earlier systems.\n",
		term.Comment("cp1252"))
	fmt.Fprintf(wr, "  Macintosh (%s) is found on Mac OS 9 and earlier systems.\n",
		term.Comment("macintosh"))
	fmt.Fprintf(wr, "\nModern systems, including %s, PCs, and the web, use Unicode UTF-8 today.\n", meta.Name)
	fmt.Fprintln(wr, "As a subset, UTF-8 is backwards compatible with US-ASCII. For example, the")
	fmt.Fprintln(wr, "capital letter A is represented by the same byte value in both encodings.")
}

// Rows return character encoding details for use in a text table.
func Rows(e encoding.Encoding) (Row, error) {
	if e == nil {
		return Row{}, ErrNil
	}
	r := Row{
		Name: fmt.Sprint(e),
	}
	switch e {
	case xud.XUserDefined1963, xud.XUserDefined1965, xud.XUserDefined1967:
		r.Value = xud.Name(e)
		r.Numeric = xud.Numeric(e)
		r.Alias = xud.Alias(e)
		return r, nil
	}
	var err error
	if r.Value, err = htmlindex.Name(e); err != nil {
		r.Value, err = ianaindex.MIME.Name(e)
		if err != nil {
			return Row{}, err
		}
	} else {
		r.Alias, err = ianaindex.MIME.Name(e)
		if err != nil {
			return Row{}, err
		}
	}
	r.Value = strings.ToLower(Uniform(r.Value))
	if i := Numeric(r.Name); i > -1 {
		r.Numeric = strconv.Itoa(i)
	}
	r.Alias, err = Alias(r.Alias, r.Value, e)
	if err != nil {
		return Row{}, err
	}
	return r, nil
}

// Numeric returns a numeric alias for a character encoding.
// A -1 int is returned whenever an alias could not be generated.
// Unicode based encodings always return -1.
func Numeric(name string) int {
	name = strings.ToLower(name)
	if strings.Contains(name, "utf") {
		return -1
	}
	s1, s2 := strings.Split(name, " "), strings.Split(name, "-")
	if len(s1) < 1 || len(s2) < 1 {
		return -1
	}
	if i, err := strconv.Atoi(s1[len(s1)-1]); err == nil {
		return i
	}
	if i, err := strconv.Atoi(s2[len(s2)-1]); err == nil {
		return i
	}
	return -1
}

// Alias returns an alias for a encoding.
// Only the alias argument is required.
func Alias(alias, value string, e encoding.Encoding) (string, error) {
	a := strings.ToLower(alias)
	if a == value {
		a = ""
	}
	if a != "" {
		return a, nil
	}
	if s := customValues(value); s != "" {
		return s, nil
	}
	a, err := ianaindex.MIB.Name(e)
	if err != nil {
		return "", err
	}
	a = strings.ToLower(a)
	if a == value {
		return "", nil
	}
	if len(a) > 2 && a[:2] == "pc" {
		return "", nil
	}
	if len(a) == 9 && a[:8] == latin {
		return "latin" + a[8:], nil
	}
	if len(a) > 9 && a[:8] == latin {
		return a[8:], nil
	}
	return a, nil
}

// customvalues for aliases can be added here.
func customValues(value string) string {
	switch value {
	case "cp437":
		return "msdos"
	case "cp850":
		return "latinI"
	case "cp852":
		return "latinII"
	case "macintosh":
		return "mac"
	case "big5":
		return "big-5"
	default:
		return ""
	}
}

// Uniform formats MIME values.
func Uniform(mime string) string {
	const limit = 1
	s := mime
	s = strings.Replace(s, "IBM00", "CP", limit)
	s = strings.Replace(s, "IBM01", "CP1", limit)
	s = strings.Replace(s, "IBM", "CP", limit)
	s = strings.Replace(s, "windows-", "CP", limit)
	return s
}
