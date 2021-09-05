package convert

import (
	"bytes"
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

type cell struct {
	name    string
	value   string
	numeric string
	alias   string
}

const latin = "isolatin"

// Encodings returns all the supported legacy text encodings.
func Encodings() []encoding.Encoding {
	a, e := charmap.All, []encoding.Encoding{}
	a = append(a, japanese.All...)
	a = append(a, unicode.All...)
	a = append(a, utf32.All...)
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
func List() *bytes.Buffer {
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
		c := cells(e)
		switch e {
		case charmap.ISO8859_10:
			fmt.Fprintf(w, " %s\t %s\t %s\t %s\t\n", c.name, c.value, c.numeric, c.alias)
			fmt.Fprintf(w, " %s\t %s\t %s\t %s\t\n", "ISO 8895-11", "iso-8895-11", "11", "iso889511")
			continue
		case charmap.CodePage037, charmap.CodePage1047, charmap.CodePage1140:
			fmt.Fprintf(w, " *%s\t %s\t %s\t %s\t\n", c.name, c.value, c.numeric, c.alias)
			continue
		case AsaX34_1963, AsaX34_1965, AnsiX34_1967:
			fmt.Fprintf(w, " **%s\t %s\t %s\t %s\t\n", c.name, c.value, c.numeric, c.alias)
			continue
		}
		// do not use ANSI colors in cells as it will break the table layout
		fmt.Fprintf(w, " %s\t %s\t %s\t %s\t\n", c.name, c.value, c.numeric, c.alias)
	}
	fmt.Fprintln(w, "\n "+str.ColInf("*")+" EBCDIC encoding, is in use on IBM mainframes and not ASCII compatible.")
	fmt.Fprintln(w, str.ColInf("**")+" ANSI X3.2 encodings are only usable with the "+str.Example("list table")+" command.")
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

// cells return character encoding details for use in a text table.
func cells(e encoding.Encoding) cell {
	if e == nil {
		return cell{}
	}
	c := cell{
		name: fmt.Sprint(e),
	}
	switch e {
	case AsaX34_1963, AsaX34_1965, AnsiX34_1967:
		c.value = x34(e)
		return c
	}
	var err error
	if c.value, err = htmlindex.Name(e); err == nil {
		c.alias, err = ianaindex.MIME.Name(e)
		if err != nil {
			log.Fatal(fmt.Errorf("list cells html index mime name: %w", err))
		}
	} else {
		c.value, err = ianaindex.MIME.Name(e)
		if err != nil {
			log.Fatal(fmt.Errorf("list cells mime name: %w", err))
		}
	}
	c.value = strings.ToLower(uniform(c.value))
	s1, s2 := strings.Split(c.name, " "), strings.Split(c.name, "-")
	if i, err := strconv.Atoi(s1[len(s1)-1]); err == nil {
		c.numeric = fmt.Sprint(i)
	} else if i, err := strconv.Atoi(s2[len(s2)-1]); err == nil {
		c.numeric = fmt.Sprint(i)
	}
	c.alias = alias(c.alias, c.value, e)
	return c
}

// alias return character encoding aliases.
func alias(s, val string, e encoding.Encoding) string {
	a := strings.ToLower(s)
	if a == val {
		a = ""
	}
	if a == "" {
		switch val {
		case "cp437":
			return "msdos"
		case "cp850":
			return "latinI"
		case "cp852":
			return "latinII"
		case "macintosh":
			return "mac"
		}
		var err error
		a, err = ianaindex.MIB.Name(e)
		if err != nil {
			return ""
		}
		a = strings.ToLower(a)
		if a == val {
			return ""
		}
		if len(a) > 2 && a[:2] == "pc" {
			return ""
		}
		if len(a) == 9 && a[:8] == latin {
			return "latin" + a[8:]
		}
		if len(a) > 9 && a[:8] == latin {
			return a[8:]
		}
	}
	return a
}

// uniform formats MIME values.
func uniform(mime string) string {
	const limit = 1
	s := mime
	s = strings.Replace(s, "IBM00", "CP", limit)
	s = strings.Replace(s, "IBM01", "CP1", limit)
	s = strings.Replace(s, "IBM", "CP", limit)
	s = strings.Replace(s, "windows-", "CP", limit)
	return s
}
