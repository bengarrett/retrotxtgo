package table

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/pkg/asa"
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

var ErrNilEncoding = errors.New("character encoding cannot be a nil value")

const latin = "isolatin"

// Row is an item for the list of code pages.
type Row struct {
	Name    string // Name is the formal name of the character encoding.
	Value   string // Value is the short name of the character encoding.
	Numeric string // Numeric is an optional, shorter numeric value of the character encoding.
	Alias   string // Alias is an optional, informal but common use value of the character encoding.
}

// Charmaps returns all the supported legacy text encodings.
func Charmaps() []encoding.Encoding {
	e := []encoding.Encoding{}
	// create a collection of all the encodings
	a := charmap.All
	a = append(a, japanese.All...)
	a = append(a, unicode.All...)
	a = append(a, utf32.All...)
	// iterate the collection and skip the unwanted and duplicate encodings
	for _, m := range a {
		switch m {
		case japanese.EUCJP,
			japanese.ISO2022JP,
			charmap.MacintoshCyrillic:
			// charmap.XUserDefined: // XUserDefined creates a duplicate of Windows 874.
			continue
		}
		e = append(e, m)
	}
	return e
}

// List returns a tabled list of supported IANA character set encodings.
func List(wr io.Writer) error { //nolint:funlen
	if wr == nil {
		wr = io.Discard
	}
	const header, title = " Formal name\t Named value\t Numeric value\t Alias value\t",
		" Known legacy code pages and character encodings "
	const verticalBars = tabwriter.Debug
	const padding, width = 2, 76
	w := tabwriter.NewWriter(wr, 0, 0, padding, ' ', verticalBars)
	if _, err := term.Head(w, width, title); err != nil {
		return err
	}
	fmt.Fprintf(w, "\n%s\n", header)
	x := Charmaps()
	x = append(x, asa.XUserDefined1963, asa.XUserDefined1965, asa.XUserDefined1967)
	for _, e := range x {
		if e == charmap.XUserDefined {
			continue
		}
		c, err := Rows(e)
		if err != nil {
			return err
		}
		switch e {
		case charmap.ISO8859_10:
			fmt.Fprintf(w, " %s\t %s\t %s\t %s\t\n",
				c.Name, c.Value, c.Numeric, c.Alias)
			// intentionally insert ISO-8895-11 after 10.
			fmt.Fprintf(w, " %s\t %s\t %s\t %s\t\n",
				"ISO 8895-11", "iso-8895-11", "11", "iso889511")
			continue
		case charmap.CodePage037, charmap.CodePage1047, charmap.CodePage1140:
			fmt.Fprintf(w, " * %s\t %s\t %s\t %s\t\n",
				c.Name, c.Value, c.Numeric, c.Alias)
			continue
		case
			unicode.UTF16(unicode.BigEndian, unicode.UseBOM),
			unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
			unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
			utf32.UTF32(utf32.BigEndian, utf32.UseBOM),
			utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM),
			utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM):
			fmt.Fprintf(w, " † %s\t %s\t %s\t %s\t\n",
				c.Name, c.Value, c.Numeric, c.Alias)
			continue
		case asa.XUserDefined1963, asa.XUserDefined1965, asa.XUserDefined1967:
			fmt.Fprintf(w, " ⁑ %s\t %s\t %s\t %s\t\n",
				c.Name, c.Value, c.Numeric, c.Alias)
			continue
		}
		// do not use ANSI colors in cells as it will break the table layout
		fmt.Fprintf(w, " %s\t %s\t %s\t %s\t\n",
			c.Name, c.Value, c.Numeric, c.Alias)
	}
	fmt.Fprintln(w, "\n "+term.Info("*")+
		" A EBCDIC encoding in use on IBM mainframes but is not ASCII compatible.")
	fmt.Fprintln(w, " "+term.Info("†")+
		" UTF-16/32 encodings are NOT usable with the "+term.Example("table")+" command.")
	fmt.Fprintln(w, " "+term.Info("⁑")+
		" ANSI X3.4 encodings are only usable with the "+term.Example("table")+" command."+
		"\n   You can use the "+term.Example("table ascii")+" command to list all three X3.4 tables.")
	fmt.Fprintln(w, "\nEither named, numeric or alias values are valid codepage arguments.")
	fmt.Fprintln(w, "  These values all match ISO 8859-1.")
	cmds := fmt.Sprintf("%s table ", meta.Bin)
	fmt.Fprintf(w, "  %s%s  %s\n",
		term.Example(cmds), term.Comment("iso-8859-1"), term.Fuzzy("# named"))
	fmt.Fprintf(w, "  %s%s           %s\n",
		term.Example(cmds), term.Comment("1"), term.Fuzzy("# numeric"))
	fmt.Fprintf(w, "  %s%s      %s\n",
		term.Example(cmds), term.Comment("latin1"), term.Fuzzy("# alias"))
	fmt.Fprintf(w, "\n  IBM Code Page 437 (%s) is commonly used on MS-DOS and ANSI art.\n",
		term.Comment("cp437"))
	fmt.Fprintf(w, "  ISO 8859-1 (%s) is found on historic Unix, Amiga and the early Internet.\n",
		term.Comment("latin1"))
	fmt.Fprintf(w, "  Windows 1252 (%s) is found on Windows ME/98 and earlier systems.\n",
		term.Comment("cp1252"))
	fmt.Fprintf(w, "  Macintosh (%s) is found on Mac OS 9 and earlier systems.\n",
		term.Comment("macintosh"))
	fmt.Fprintf(w, "\n%s, PCs and the web today use Unicode UTF-8. As a subset,\n", meta.Name)
	fmt.Fprintln(w, "UTF-8 is backwards compatible with both ISO 8895-1 and US-ASCII.")
	return w.Flush()
}

// Rows return character encoding details for use in a text table.
func Rows(e encoding.Encoding) (Row, error) {
	if e == nil {
		return Row{}, ErrNilEncoding
	}
	r := Row{
		Name: fmt.Sprint(e),
	}
	switch e {
	case asa.XUserDefined1963, asa.XUserDefined1965, asa.XUserDefined1967:
		r.Value = asa.Name(e)
		r.Numeric = asa.Numeric(e)
		r.Alias = asa.Alias(e)
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
		r.Numeric = fmt.Sprint(i)
	}
	r.Alias, err = Alias(e, r.Alias, r.Value)
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
	if i, err := strconv.Atoi(s1[len(s1)-1]); err == nil {
		return i
	}
	if i, err := strconv.Atoi(s2[len(s2)-1]); err == nil {
		return i
	}
	return -1
}

// Alias return character encoding aliases.
func Alias(e encoding.Encoding, alias, value string) (string, error) {
	a := strings.ToLower(alias)
	if a == value {
		a = ""
	}
	if a != "" {
		return a, nil
	}
	switch value {
	case "cp437":
		return "msdos", nil
	case "cp850":
		return "latinI", nil
	case "cp852":
		return "latinII", nil
	case "macintosh":
		return "mac", nil
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
