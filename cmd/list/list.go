// Package list provides the list command run function.
package list

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/sample"
	"github.com/bengarrett/retrotxtgo/table"
	"github.com/bengarrett/retrotxtgo/term"
	"github.com/bengarrett/retrotxtgo/xud"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
)

var (
	ErrTable = errors.New("could not display the table")
	ErrIANA  = errors.New("could not work out the IANA index or MIME type")
)

// Examples writes the list command examples.
func Examples(wr io.Writer) error {
	if wr == nil {
		wr = io.Discard
	}
	const width = 80
	m := sample.Map()
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	const padding = 2
	w := tabwriter.NewWriter(wr, 0, 0, padding, ' ', 0)
	bin := fmt.Sprintf("  %s ", meta.Bin)
	header := "Packaged example text and ANSI files to test and play with " + meta.Name
	term.Head(w, width, header)
	fmt.Fprintf(w, "\nAny of these packaged examples will work with the %s and %s commands.\n",
		term.Example("info"), term.Example("view"))
	fmt.Fprintln(w)
	for _, k := range keys {
		fmt.Fprintf(w, " %s\t%s\t\n", k, m[k].Description)
	}
	fmt.Fprintln(w)
	term.Head(w, width, "Usage sample commands")
	fmt.Fprintf(w, "\nPrint the Windows-1252 English test to the terminal.\n%s\n",
		term.Example(bin+"view 1252"))
	fmt.Fprintf(w, "\nConvert the Windows-1252 English test to UTF-8 encoding and save it to a file.\n%s\n",
		term.Example(bin+"view 1252 > file.txt"))
	fmt.Fprintf(w, "\nSave the Windows-1252 English test with its original encoding.\n%s\n",
		term.Example(bin+"view 1252 --original > file.txt"))
	fmt.Fprintf(w, "\nDisplay statistics and information from a piped source.\n%s\n",
		term.Example(fmt.Sprintf("%sview 1252 --original | %s info", bin, meta.Bin)))
	fmt.Fprintf(w, "\nDisplay statistics and information from the Windows-1252 English test.\n%s\n",
		term.Example(bin+"info 1252"))
	fmt.Fprintf(w, "\nDisplay statistics, information and SAUCE metadata from the SAUCE test.\n%s\n",
		term.Example(bin+"info sauce"))
	fmt.Fprintf(w, "\nMultiple examples used together are supported.\n%s\n",
		term.Example(bin+"view ansi ascii ansi.rgb"))
	return w.Flush()
}

// Table writes one or more named encodings as a formatted table.
func Table(w io.Writer, names ...string) error {
	if w == nil {
		w = io.Discard
	}
	// custom ascii shortcut
	tables := names
	for i, name := range tables {
		if name != "ascii" {
			continue
		}
		names[i] = "ascii-67"
		names = append(names[:i+1], names[i:]...)
		names[i] = "ascii-65"
		names = append(names[:i+1], names[i:]...)
		names[i] = "ascii-63"
	}
	// iterate through the tables
	for _, name := range names {
		if err := table.Table(w, name); err != nil {
			return err
		}
		fmt.Fprintln(w)
	}
	return nil
}

// Tables writes all the supported encodings as formatted tables.
func Tables(w io.Writer) error {
	if w == nil {
		w = io.Discard
	}
	// use strings builder to reduce memory usage
	// https://yourbasic.org/golang/build-append-concatenate-strings-efficiently/
	tables := make([]encoding.Encoding, 0, len(table.Charmaps()))
	encodings := table.Charmaps()
	encodings = append(encodings,
		xud.XUserDefined1963,
		xud.XUserDefined1965,
		xud.XUserDefined1967)
	// reorder tables to position X-User-Defined after ISO-8859-10
	for _, e := range encodings {
		switch e {
		case charmap.ISO8859_10:
			tables = append(tables, charmap.ISO8859_10)
			tables = append(tables, xud.XUserDefinedISO11)
			continue
		case xud.XUserDefinedISO11:
			continue
		}
		tables = append(tables, e)
	}
	tables = slices.Compact(tables)
	// print tables
	for _, e := range tables {
		var (
			err  error
			name string
		)
		switch e {
		case charmap.XUserDefined:
			continue
		case xud.XUserDefinedISO11:
			name = fmt.Sprint(e)
		case xud.XUserDefined1963,
			xud.XUserDefined1965,
			xud.XUserDefined1967:
			name = xud.Name(e)
		}
		if name == "" {
			name, err = ianaindex.MIME.Name(e)
			if err != nil {
				return fmt.Errorf("table %s, %w, %w", e, ErrIANA, err)
			}
		}
		if !Printable(name) {
			continue
		}
		if err := table.Table(w, name); err != nil {
			return fmt.Errorf("table %s, %w, %w", e, ErrTable, err)
		}
	}
	return nil
}

// Printable reports whether the named encoding can be shown as
// a 256 character table. UTF-16 and UTF-32 are not printable.
func Printable(name string) bool {
	const (
		utf16 = "utf-16"
		l     = len(utf16)
	)
	s := strings.ToLower(name)
	if s == "" {
		return false
	}
	if len(s) < l {
		return true
	}
	if s[:l] == utf16 {
		return false
	}
	if s[:l] == "utf-32" {
		return false
	}
	return true
}
