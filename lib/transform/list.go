package transform

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
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
	var (
		buf bytes.Buffer
		ii  iana
		err error
	)
	w := new(tabwriter.Writer)
	w.Init(&buf, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, str.Cp(" Known legacy code pages and character encodings"))
	fmt.Fprintln(w, str.Cf(strings.Repeat("\u2015", 49)))
	fmt.Fprintln(w, "\tName\tValue\tAlias")
	// append other supported legacy encodings
	c := Encodings()
	for _, n := range c {
		name := fmt.Sprint(n)
		if name == "X-User-Defined" {
			continue
		}
		if ii.mime, err = ianaindex.MIME.Name(n); err != nil {
			continue
		}
		ii.index, _ = ianaindex.IANA.Name(n)
		ii.mib, _ = ianaindex.MIB.Name(n)
		ii.s = strings.Split(name, " ")
		fmt.Fprintf(w, "\t%s\t%s", name, str.Ci(uniform(ii.mime)))
		fmt.Fprintf(w, "\t%s", ii.mib)
		fmt.Fprint(w, "\n")
	}
	fmt.Fprintln(w, "\tName\tValue\tAlias")
	fmt.Fprint(w, "\n"+str.Cinf("*")+" IBM Code Page 437 ("+str.Cc("CP437")+") is commonly used by MS-DOS English text and ANSI art")
	fmt.Fprint(w, "\n  ISO 8859-1 ("+str.Cc("ISOLatin1")+") is found in legacy Internet, Unix and Amiga documents")
	fmt.Fprint(w, "\n  Windows 1252 ("+str.Cc("CP1252")+") is found in legacy English language Windows operating systems")
	w.Flush()
	return &buf
}

func uniform(mime string) (s string) {
	s = mime
	s = strings.Replace(s, "IBM00", "CP", 1)
	s = strings.Replace(s, "IBM01", "CP1", 1)
	s = strings.Replace(s, "IBM", "CP", 1)
	s = strings.Replace(s, "windows-", "CP", 1)
	return s
}
