package convert

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/tabwriter"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"

	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

type cell struct {
	name    string
	value   string
	numeric string
	alias   string
}

// Encodings returns all the supported legacy text encodings.
func Encodings() (e []encoding.Encoding) {
	a := append(charmap.All, japanese.All...)
	a = append(a, unicode.All...)
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
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', flags)
	fmt.Fprintln(w, str.Cp(title))
	fmt.Fprintln(w, strings.Repeat("-", len(title)))
	fmt.Fprintln(w, header)
	fmt.Println(Encodings())
	for _, e := range Encodings() {
		if e == charmap.XUserDefined {
			fmt.Fprintf(w, " %s\t %s\t %s\t %s\t\n", "ISO 8895-11", "iso-8895-11", "11", "iso889511")
			continue
		}
		c := cells(e)
		// do not use ANSI colors in cells as it will break the table layout
		fmt.Fprintf(w, " %s\t %s\t %s\t %s\t\n", c.name, c.value, c.numeric, c.alias)
	}
	fmt.Fprintln(w, "\nEither argument, numeric or alias values are valid codepage arguments")
	fmt.Fprintln(w, "All these codepage arguments will match ISO 8859-1")
	fmt.Fprintln(w, "\n"+str.Example("retrotxt list table iso-8859-1"))
	fmt.Fprintln(w, str.Example("retrotxt list table 1"))
	fmt.Fprintln(w, str.Example("retrotxt list table latin1"))
	fmt.Fprintln(w, "\n"+str.Cinf("*")+" IBM Code Page 437 ("+str.Cc("cp437")+") is commonly used on MS-DOS and with ANSI art")
	fmt.Fprintln(w, "  ISO 8859-1 ("+str.Cc("latin1")+") is found on legacy Unix, Amiga and the early Internet")
	fmt.Fprintln(w, "  Windows 1252 ("+str.Cc("cp1252")+") is found on legacy Windows 9x and earlier systems")
	fmt.Fprintln(w, "  Macintosh ("+str.Cc("macintosh")+") is found on Mac OS 9 and earlier systems")
	fmt.Fprintln(w, "  RetroTxt, modern systems and the web today use UTF-8, a Unicode encoding")
	fmt.Fprintln(w, "  that's a subset of ISO 8859-1 which itself is a subset of US-ASCII")
	if err := w.Flush(); err != nil {
		logs.Fatal("convert list", "flush", err)
	}
	return &buf
}

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
		if len(a) == 9 && a[:8] == "isolatin" {
			return "latin" + a[8:]
		}
		if len(a) > 9 && a[:8] == "isolatin" {
			return a[8:]
		}
	}
	return a
}

func uniform(mime string) (s string) {
	s = mime
	s = strings.Replace(s, "IBM00", "CP", 1)
	s = strings.Replace(s, "IBM01", "CP1", 1)
	s = strings.Replace(s, "IBM", "CP", 1)
	s = strings.Replace(s, "windows-", "CP", 1)
	return s
}
