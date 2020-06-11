package transform

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/japanese"
)

type iana struct {
	mime  string
	index string
	mib   string
	s     []string
}

// Encodings returns all the supported legacy text encodings.
func Encodings() (e []encoding.Encoding) {
	a := append(charmap.All, japanese.All...)
	for _, m := range a {
		switch fmt.Sprintf("%v", m) {
		case
			"EUC-JP",
			"ISO-2022-JP",
			"Macintosh Cyrillic",
			"X-User-Defined":
			continue
		}
		e = append(e, m)
	}
	return e
}

// List returns a tabled list of supported IANA character set encodings
func List() *bytes.Buffer {
	var buf bytes.Buffer
	header := "\tname\tvalue\tcommon alias"
	w := new(tabwriter.Writer)
	w.Init(&buf, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, str.Cp(" Known legacy code pages and character encodings"))
	fmt.Fprintln(w, str.Cf(strings.Repeat("\u2015", 49)))
	fmt.Fprintln(w, header)
	for _, e := range Encodings() {
		n, v, a := cells(e)
		fmt.Fprintf(w, "\t%s\t%s\t%s\n", n, str.Ci(v), a) // name, value, alias
	}
	fmt.Fprintln(w, header)
	fmt.Fprint(w, "\n"+str.Cinf("*")+" IBM Code Page 437 ("+str.Cc("CP437")+") is commonly used by MS-DOS English text and ANSI art")
	fmt.Fprint(w, "\n  ISO 8859-1 ("+str.Cc("Latin1")+") is found in legacy Internet, Unix and Amiga documents")
	fmt.Fprint(w, "\n  Windows 1252 ("+str.Cc("CP1252")+") is found in legacy English language Windows operating systems")
	w.Flush()
	return &buf
}

func cells(e encoding.Encoding) (n, v, a string) {
	if e == nil {
		return n, v, a
	}
	var err error
	if v, err = htmlindex.Name(e); err == nil {
		a, _ = ianaindex.MIME.Name(e)
	} else {
		v, _ = ianaindex.MIME.Name(e)
	}
	v = strings.ToLower(uniform(v))
	a = alias(a, v, e)
	return fmt.Sprint(e), v, a
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
		}
		a, _ = ianaindex.MIB.Name(e)
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
