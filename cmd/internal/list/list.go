package list

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/sample"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
)

var (
	ErrTable = errors.New("could not display the table")
	ErrIANA  = errors.New("could not work out the IANA index or MIME type")
)

func Examples() *bytes.Buffer {
	m := sample.Map()
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	const padding = 2
	w := tabwriter.NewWriter(&buf, 0, 0, padding, ' ', 0)
	bin := fmt.Sprintf("  %s ", meta.Bin)
	fmt.Fprintf(w, "%s\n",
		str.Head(0, fmt.Sprintf("Packaged example text and ANSI files to test and play with %s", meta.Name)))
	for _, k := range keys {
		fmt.Fprintf(w, "%s\t%s\t\n", k, m[k].Description)
	}
	fmt.Fprintf(w, "\nAny of these packaged examples will work with the %s, %s and %s commands.\n",
		str.Example("create"), str.Example("info"), str.Example("view"))
	fmt.Fprintf(w, "\nPrint the Windows-1252 English test to the terminal.\n%s\n",
		str.Example(bin+"view 1252"))
	fmt.Fprintf(w, "\nConvert the Windows-1252 English test to UTF-8 encoding and save it to a file.\n%s\n",
		str.Example(bin+"view 1252 > file.txt"))
	fmt.Fprintf(w, "\nSave the Windows-1252 English test with its original encoding.\n%s\n",
		str.Example(bin+"view --to=cp1252 1252 > file.txt"))
	fmt.Fprintf(w, "\nDisplay statistics and information from a piped source.\n%s\n",
		str.Example(fmt.Sprintf("%sview --to=cp1252 1252 | %s info", bin, meta.Bin)))
	fmt.Fprintf(w, "\nDisplay statistics and information from the Windows-1252 English test.\n%s\n",
		str.Example(bin+"info 1252"))
	fmt.Fprintf(w, "\nDisplay statistics, information and SAUCE metadata from the SAUCE test.\n%s\n",
		str.Example(bin+"info sauce"))
	fmt.Fprintf(w, "\nCreate and display a HTML document from the Windows-1252 English test.\n%s\n",
		str.Example(bin+"create 1252"))
	fmt.Fprintf(w, "\nCreate and save the HTML and assets from the Windows-1252 English test.\n%s\n",
		str.Example(bin+"create 1252 --save"))
	fmt.Fprintf(w, "\nServe the Windows-1252 English test over a local web server.\n%s\n",
		str.Example(bin+"create 1252 -p0"))
	fmt.Fprintf(w, "\nMultiple examples used together are supported.\n%s\n",
		str.Example(bin+"view ansi ascii ansi.rgb"))
	if err := w.Flush(); err != nil {
		logs.FatalWrap(logs.ErrTabFlush, err)
	}
	return &buf
}

// listTable returns one or more named encodings in a tabled format.
func Table(names ...string) (string, error) {
	// custom ascii shortcut
	ns := names
	for i, name := range ns {
		if name != "ascii" {
			continue
		}
		names[i] = "ascii-67"
		names = append(names[:i+1], names[i:]...)
		names[i] = "ascii-65"
		names = append(names[:i+1], names[i:]...)
		names[i] = "ascii-63"
	}
	s := ""
	// iterate through the tables
	for _, name := range names {
		table, err := convert.Table(name)
		if err != nil {
			return "", err
		}
		s = fmt.Sprintf("%s%s", s, table.String())
	}
	return s, nil
}

// listTbls returns all the supported encodings in a tabled format.
func Tables() (s string) {
	enc := convert.Encodings()
	var tables []encoding.Encoding
	// reorder tables to position X-User-Defined after ISO-8859-10
	for _, tbl := range enc {
		switch tbl {
		case charmap.ISO8859_10:
			tables = append(tables, charmap.ISO8859_10)
			tables = append(tables, charmap.XUserDefined)
			continue
		case charmap.XUserDefined:
			continue
		}
		tables = append(tables, tbl)
	}
	// print tables
	for _, tbl := range tables {
		var (
			err  error
			name string
		)
		if tbl == charmap.XUserDefined {
			name = "iso-8859-11"
		}
		if name == "" {
			name, err = ianaindex.MIME.Name(tbl)
			if err != nil {
				fmt.Println(logs.SprintMark(fmt.Sprint(tbl), ErrIANA, err))
				continue
			}
		}
		if !Printable(name) {
			continue
		}
		table, err := convert.Table(name)
		if err != nil {
			fmt.Println(logs.SprintMark(name, ErrTable, err))
			continue
		}
		s = fmt.Sprintf("%s%s", s, table.String())
	}
	return s
}

// Printable returns true if the named encoding be shown in an 8-bit table.
func Printable(name string) bool {
	switch strings.ToLower(name) {
	case "", "utf-16", "utf-16be", "utf-16le", "utf-32", "utf-32be", "utf-32le":
		return false
	}
	return true
}
