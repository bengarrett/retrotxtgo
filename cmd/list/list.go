// Package list provides the list command run function.
package list

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"sort"
	"strings"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/sample"
	"github.com/bengarrett/retrotxtgo/table"
	"github.com/bengarrett/retrotxtgo/xud"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/traditionalchinese"
)

var (
	ErrTable = errors.New("could not display the table")
	ErrIANA  = errors.New("could not work out the IANA index or MIME type")
)

const width = 80

var (
	tabletp = "╭" + strings.Repeat("─", width-1)
	tableln = "│" + strings.Repeat(" ", width-1)
	tablebm = "╰" + strings.Repeat("─", width-1)
)

func nameStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)
}

func descStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))
}

func usageStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("244"))
}

func commandStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)
}

func ansiExample(wr io.Writer, keys []string) io.Writer {
	const asm = "ansi.setmodes"
	ansiExamples := []string{
		"ansi\t", "ansi.aix", "ansi.blank", "ansi.cp",
		"ansi.cpf", "ansi.hvp", "ansi.proof", "ansi.rgb", asm,
	}
	if wr == nil || !hasExamples(keys, ansiExamples) {
		return wr
	}
	m := sample.Map()
	fmt.Fprintln(wr, tabletp)
	fmt.Fprintln(wr, "│ ANSI art and color examples")
	fmt.Fprintln(wr, tableln)
	for _, k := range keys {
		if slices.Contains(ansiExamples, k) {
			name := nameStyle().Render(k)
			desc := descStyle().Render(m[k].Description)
			if k == asm {
				fmt.Fprintf(wr, "│  %s %s\n", name, desc)
				continue
			}
			fmt.Fprintf(wr, "│  %s\t %s\n", name, desc)
		}
	}
	fmt.Fprintln(wr, tablebm)
	return wr
}

func encodingExample(wr io.Writer, keys []string) io.Writer {
	// Encoding Test Files
	const cr437, lf437 = "437.cr", "437.lf"
	encodingExamples := []string{"037", "1252", "437", cr437, lf437, "865"}
	if wr == nil || !hasExamples(keys, encodingExamples) {
		return wr
	}
	m := sample.Map()
	fmt.Fprintln(wr, tabletp)
	fmt.Fprintln(wr, "│ Encoding test files")
	fmt.Fprintln(wr, tableln)
	for _, k := range keys {
		if slices.Contains(encodingExamples, k) {
			name := nameStyle().Render(k)
			desc := descStyle().Render(m[k].Description)
			switch k {
			case cr437, lf437:
				fmt.Fprintf(wr, "│  %s\t %s\n", name, desc)
			default:
				fmt.Fprintf(wr, "│  %s\t\t %s\n", name, desc)
			}
		}
	}
	fmt.Fprintln(wr, tablebm)
	return wr
}

func plainExample(wr io.Writer, keys []string) io.Writer {
	textExamples := []string{"ascii", "iso-1", "iso-15", "us-ascii"}
	if wr == nil || !hasExamples(keys, textExamples) {
		return wr
	}
	m := sample.Map()
	fmt.Fprintln(wr, tabletp)
	fmt.Fprintln(wr, "│ Plain text files")
	fmt.Fprintln(wr, tableln)
	for _, k := range keys {
		if slices.Contains(textExamples, k) {
			name := nameStyle().Render(k)
			desc := descStyle().Render(m[k].Description)
			fmt.Fprintf(wr, "│  %s\t %s\n", name, desc)
		}
	}
	fmt.Fprintln(wr, tablebm)
	return wr
}

func unicodeExample(wr io.Writer, keys []string) io.Writer {
	unicodeExamples := []string{"utf8\t", "utf8.bom", "utf16", "utf16.be", "utf16.le", "utf32", "utf32.be", "utf32.le"}
	if wr == nil || !hasExamples(keys, unicodeExamples) {
		return wr
	}
	m := sample.Map()
	fmt.Fprintln(wr, tabletp)
	fmt.Fprintln(wr, "│ Unicode test files")
	fmt.Fprintln(wr, tableln)
	for _, k := range keys {
		if slices.Contains(unicodeExamples, k) {
			name := nameStyle().Render(k)
			desc := descStyle().Render(m[k].Description)
			fmt.Fprintf(wr, "│  %s\t %s\n", name, desc)
		}
	}
	fmt.Fprintln(wr, tablebm)

	return wr
}

func specialExample(wr io.Writer, keys []string) io.Writer {
	specialExamples := []string{"sauce", "shiftjis"}
	if wr == nil || !hasExamples(keys, specialExamples) {
		return wr
	}
	m := sample.Map()
	fmt.Fprintln(wr, tabletp)
	fmt.Fprintln(wr, "│ Special format files")
	fmt.Fprintln(wr, tableln)
	for _, k := range keys {
		if slices.Contains(specialExamples, k) {
			name := nameStyle().Render(k)
			desc := descStyle().Render(m[k].Description)
			fmt.Fprintf(wr, "│  %s\t %s\n", name, desc)
		}
	}
	fmt.Fprintln(wr, tablebm)
	return wr
}

// Examples writes the list command examples.
func Examples(wr io.Writer) error {
	if wr == nil {
		wr = io.Discard
	}
	m := sample.Map()
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	bin := fmt.Sprintf("  %s ", meta.Bin)
	wr = ansiExample(wr, keys)
	wr = encodingExample(wr, keys)
	wr = plainExample(wr, keys)
	wr = unicodeExample(wr, keys)
	wr = specialExample(wr, keys)

	// Usage section as a complete table
	fmt.Fprintln(wr)
	fmt.Fprintln(wr, tabletp)
	fmt.Fprintln(wr, "│ Usage commands, all examples work with info and view")
	fmt.Fprintln(wr, tableln)
	fmt.Fprintf(wr, "│  %s\n", usageStyle().Render(
		"Print the Windows-1252 English test"))
	fmt.Fprintf(wr, "│  %s\n", commandStyle().Render(bin+"view 1252"))
	fmt.Fprintf(wr, "│  %s\n", usageStyle().Render(
		"Convert and save the Windows-1252 English test as UTF-8"))
	fmt.Fprintf(wr, "│  %s\n", commandStyle().Render(bin+"view 1252 > file.txt"))
	fmt.Fprintf(wr, "│  %s\n", usageStyle().Render(
		"Save the Windows-1252 English test"))
	fmt.Fprintf(wr, "│  %s\n", commandStyle().Render(bin+"view 1252 --original > file.txt"))
	fmt.Fprintf(wr, "│  %s\n", usageStyle().Render(
		"Display statistics and information from a piped source"))
	fmt.Fprintf(wr, "│  %s\n", commandStyle().Render(
		fmt.Sprintf("%sview 1252 --original | %s info", bin, meta.Bin)))
	fmt.Fprintf(wr, "│  %s\n", usageStyle().Render(
		"Display statistics and information from the Windows-1252 English test"))
	fmt.Fprintf(wr, "│  %s\n", commandStyle().Render(bin+"info 1252"))
	fmt.Fprintf(wr, "│  %s\n", usageStyle().Render(
		"Display statistics, information and SAUCE metadata from the SAUCE test"))
	fmt.Fprintf(wr, "│  %s\n", commandStyle().Render(bin+"info sauce"))
	fmt.Fprintf(wr, "│  %s\n", usageStyle().Render(
		"Print multiple examples"))
	fmt.Fprintf(wr, "│  %s\n", commandStyle().Render(bin+"view ansi ascii ansi.rgb"))
	fmt.Fprintln(wr, tablebm)

	return nil
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
		if i > len(names) {
			break
		}
		names[i] = "ascii-67"
		names = append(names[:i+1], names[i:]...)
		names[i] = "ascii-65"
		names = append(names[:i+1], names[i:]...)
		names[i] = "ascii-63"
	}
	// iterate through the tables
	for _, name := range names {
		if err := table.WithLipgloss(w, name); err != nil {
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
		case
			traditionalchinese.Big5,
			charmap.XUserDefined:
			// do not display these encodings
			continue
		case xud.XUserDefinedISO11:
			name = fmt.Sprint(e)
		case
			xud.XUserDefined1963,
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
		if err := table.WithLipgloss(w, name); err != nil {
			return fmt.Errorf("table %s, %w, %w", e, ErrTable, err)
		}
	}
	return nil
}

// hasExamples reports whether any of the examples exist in the keys.
func hasExamples(keys []string, examples []string) bool {
	for _, ex := range examples {
		if slices.Contains(keys, ex) {
			return true
		}
	}
	return false
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
