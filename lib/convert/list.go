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
func Encodings() (e []encoding.Encoding) {
	a := charmap.All
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
	const header, title = "Formal name\tArgument value\tNumeric value\tAlias value\t",
		" Supported legacy code pages and character encodings "
	var buf bytes.Buffer
	flags := tabwriter.Debug // tabwriter.AlignRight | tabwriter.Debug
	const padding = 2
	w := tabwriter.NewWriter(&buf, 0, 0, padding, ' ', flags)
	fmt.Fprintln(w, "\n"+str.Cp(title))
	fmt.Fprintln(w, strings.Repeat("-", len(title)))
	fmt.Fprintln(w, header)
	fmt.Println(Encodings())
	for _, e := range Encodings() {
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
			fmt.Fprintf(w, " %s*\t %s\t %s\t %s\t\n", c.name, c.value, c.numeric, c.alias)
			continue
		}
		// do not use ANSI colors in cells as it will break the table layout
		fmt.Fprintf(w, " %s\t %s\t %s\t %s\t\n", c.name, c.value, c.numeric, c.alias)
	}
	fmt.Fprintln(w, "\n"+str.Cinf("*")+" EBCDIC data encoding is used on IBM Mainframe OS and is not ASCII compatible.")
	fmt.Fprintln(w, "\nEither argument, numeric or alias values are valid codepage arguments.")
	fmt.Fprintln(w, "  These example codepage arguments all match ISO 8859-1.")
	cmds := fmt.Sprintf("%s list table ", meta.Bin)
	fmt.Fprintf(w, "  %s%s  %s\n",
		str.Example(cmds), str.Cc("iso-8859-1"), str.Cf("# argument"))
	fmt.Fprintf(w, "  %s%s           %s\n",
		str.Example(cmds), str.Cc("1"), str.Cf("# numeric"))
	fmt.Fprintf(w, "  %s%s      %s\n",
		str.Example(cmds), str.Cc("latin1"), str.Cf("# alias"))
	fmt.Fprintf(w, "\n  IBM Code Page 437 (%s) is commonly used on MS-DOS and ANSI art.\n",
		str.Cc("cp437"))
	fmt.Fprintf(w, "  ISO 8859-1 (%s) is found on legacy Unix, Amiga and the early Internet.\n",
		str.Cc("latin1"))
	fmt.Fprintf(w, "  Windows 1252 (%s) is found on legacy Windows 9x and earlier systems.\n",
		str.Cc("cp1252"))
	fmt.Fprintf(w, "  Macintosh (%s) is found on Mac OS 9 and earlier systems.\n",
		str.Cc("macintosh"))
	fmt.Fprintf(w, "\n%s, PCs and the web today use Unicode UTF-8. It is a subset of ISO 8895-1,\n", meta.Name)
	fmt.Fprintln(w, "which allows UTF-8 to be backwards compatible both with it and US-ASCII.")
	if err := w.Flush(); err != nil {
		logs.ProblemFatal(logs.ErrTabFlush, err)
	}
	return &buf
}

// Cells return character encoding details for use in a text table.
func cells(e encoding.Encoding) (c cell) {
	if e == nil {
		return cell{}
	}
	c.name = fmt.Sprint(e)
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

// Alias return character encoding aliases.
func alias(s, val string, e encoding.Encoding) string {
	a := strings.ToLower(s)
	if val == a {
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

// Uniform formats MIME values.
func uniform(mime string) (s string) {
	s = mime
	s = strings.Replace(s, "IBM00", "CP", 1)
	s = strings.Replace(s, "IBM01", "CP1", 1)
	s = strings.Replace(s, "IBM", "CP", 1)
	s = strings.Replace(s, "windows-", "CP", 1)
	return s
}
