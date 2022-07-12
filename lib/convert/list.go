package convert

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

type Cell struct {
	Name    string
	Value   string
	Numeric string
	Alias   string
}

var (
	ErrNilEncoding = errors.New("character encoding cannot be a nil value")
)

const latin = "isolatin"

// Encodings returns all the supported legacy text encodings.
func Encodings() []encoding.Encoding {
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
			//charmap.XUserDefined: // XUserDefined creates a duplicate of Windows 874.
			continue
		}
		e = append(e, m)
	}
	return e
}

// List returns a tabled list of supported IANA character set encodings.
func List() *bytes.Buffer { //nolint:funlen
	const header, title = " Formal name\t Named value\t Numeric value\t Alias value\t",
		" Supported legacy code pages and character encodings "
	var buf bytes.Buffer
	flags := tabwriter.Debug // tabwriter.AlignRight | tabwriter.Debug
	const padding, tblWidth = 2, 73
	w := tabwriter.NewWriter(&buf, 0, 0, padding, ' ', flags)
	fmt.Fprint(w, str.Head(tblWidth, title))
	fmt.Fprintf(w, "\n%s\n", header)
	enc := Encodings()
	enc = append(enc, AsaX34_1963, AsaX34_1965, AnsiX34_1967)
	for _, e := range enc {
		if e == charmap.XUserDefined {
			continue
		}
		c, err := Cells(e)
		if err != nil {
			log.Fatal(err)
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
		case AsaX34_1963, AsaX34_1965, AnsiX34_1967:
			fmt.Fprintf(w, " ⁑ %s\t %s\t %s\t %s\t\n",
				c.Name, c.Value, c.Numeric, c.Alias)
			continue
		}
		// do not use ANSI colors in cells as it will break the table layout
		fmt.Fprintf(w, " %s\t %s\t %s\t %s\t\n",
			c.Name, c.Value, c.Numeric, c.Alias)
	}
	fmt.Fprintln(w, "\n "+str.ColInf("*")+
		" A EBCDIC encoding in use on IBM mainframes but is not ASCII compatible.")
	fmt.Fprintln(w, " "+str.ColInf("†")+
		" UTF-16/32 encodings are NOT usable with the "+str.Example("list table")+" command.")
	fmt.Fprintln(w, " "+str.ColInf("⁑")+
		" ANSI X3.4 encodings are only usable with the "+str.Example("list table")+" command."+
		"\n   You can use the "+str.Example("list table ascii")+" command to list all three X3.4 tables.")
	fmt.Fprintln(w, "\nEither named, numeric or alias values are valid codepage arguments.")
	fmt.Fprintln(w, "  These values all match ISO 8859-1.")
	cmds := fmt.Sprintf("%s list table ", meta.Bin)
	fmt.Fprintf(w, "  %s%s  %s\n",
		str.Example(cmds), str.ColCmt("iso-8859-1"), str.ColFuz("# named"))
	fmt.Fprintf(w, "  %s%s           %s\n",
		str.Example(cmds), str.ColCmt("1"), str.ColFuz("# numeric"))
	fmt.Fprintf(w, "  %s%s      %s\n",
		str.Example(cmds), str.ColCmt("latin1"), str.ColFuz("# alias"))
	fmt.Fprintf(w, "\n  IBM Code Page 437 (%s) is commonly used on MS-DOS and ANSI art.\n",
		str.ColCmt("cp437"))
	fmt.Fprintf(w, "  ISO 8859-1 (%s) is found on legacy Unix, Amiga and the early Internet.\n",
		str.ColCmt("latin1"))
	fmt.Fprintf(w, "  Windows 1252 (%s) is found on legacy Windows 9x and earlier systems.\n",
		str.ColCmt("cp1252"))
	fmt.Fprintf(w, "  Macintosh (%s) is found on Mac OS 9 and earlier systems.\n",
		str.ColCmt("macintosh"))
	fmt.Fprintf(w, "\n%s, PCs and the web today use Unicode UTF-8. As a subset of ISO 8895-1,\n", meta.Name)
	fmt.Fprintln(w, "UTF-8 is backwards compatible with it and US-ASCII.")
	if err := w.Flush(); err != nil {
		logs.FatalWrap(logs.ErrTabFlush, err)
	}
	return &buf
}

// Cells return character encoding details for use in a text table.
func Cells(e encoding.Encoding) (Cell, error) {
	if e == nil {
		return Cell{}, ErrNilEncoding
	}
	c := Cell{
		Name: fmt.Sprint(e),
	}
	switch e {
	case AsaX34_1963, AsaX34_1965, AnsiX34_1967:
		c.Value = AsaX34(e)
		return c, nil
	}

	var err error
	if c.Value, err = htmlindex.Name(e); err == nil {
		c.Alias, err = ianaindex.MIME.Name(e)
		if err != nil {
			return Cell{}, err
		}
	} else {
		c.Value, err = ianaindex.MIME.Name(e)
		if err != nil {
			return Cell{}, err
		}
	}
	c.Value = strings.ToLower(Uniform(c.Value))

	if i := Numeric(c.Name); i > -1 {
		c.Numeric = fmt.Sprint(i)
	}

	c.Alias, err = AliasFmt(c.Alias, c.Value, e)
	if err != nil {
		return Cell{}, err
	}
	return c, nil
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

// AliasFmt return character encoding aliases.
func AliasFmt(alias, value string, e encoding.Encoding) (string, error) {
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
	var err error
	a, err = ianaindex.MIB.Name(e)
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
