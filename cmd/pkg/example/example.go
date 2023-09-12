// Package example provides help usage examples for the cmd package.
package example

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/gookit/color"
)

// Filenames is the placeholder for the filenames in the help usage examples.
const Filenames = "[filenames]"

// Example is the type for the help usage examples.
type Example int

const (
	Cmd      Example = iota // Cmd is the example for the root command.
	Examples                // Examples are the examples for the list examples command.
	Table                   // Table are the examples for the list tables command.
	Info                    // Info is the example for the info command.
	View                    // View is the example for the view command.
)

// Print returns help usage examples.
func (e Example) String(w io.Writer) {
	if w == nil {
		w = io.Discard
	}
	b := &bytes.Buffer{}
	// change example operating system path separator
	t := template.Must(template.New("example").Parse(e.result()))
	err := t.Execute(b, string(os.PathSeparator))
	if err != nil {
		log.Fatal(err)
	}
	// color the example text except text following
	// the last hash #, which is treated as a comment
	const cmmt, sentence = "#", 2
	scanner := bufio.NewScanner(b)
	cnt := 0
	rows := len(strings.Split(b.String(), "\n"))
	for scanner.Scan() {
		cnt++
		s := strings.Split(scanner.Text(), cmmt)
		l := len(s)
		if l < sentence {
			fmt.Fprint(w, term.Info(scanner.Text()))
			if cnt < rows {
				fmt.Fprintln(w)
			}
			continue
		}
		// do not the last hash as a comment
		ex := strings.Join(s[:l-1], cmmt)
		fmt.Fprint(w, term.Info(ex))
		fmt.Fprintf(w, "%s%s", color.Secondary.Sprint(cmmt), s[l-1])
		if cnt < rows {
			fmt.Fprintln(w)
		}
	}
}

func (e Example) result() string {
	switch e {
	case Cmd:
		return cmd()
	case Examples:
		return examples()
	case Table:
		return table()
	case Info:
		return info()
	case View:
		return view()
	}
	return ""
}

func cmd() string {
	s := &strings.Builder{}
	fmt.Fprintf(s, "  %s info %s\n", meta.Bin, Filenames)
	fmt.Fprintf(s, "  %s view %s\n", meta.Bin, Filenames)
	fmt.Fprintf(s, "  %s example\n\n", meta.Bin)
	fmt.Fprintf(s, "  %s list\n", meta.Bin)
	fmt.Fprintf(s, "  %s tables\n", meta.Bin)
	fmt.Fprintf(s, "  %s table [codepage names or aliases]", meta.Bin)
	return s.String()
}

func examples() string {
	s := &strings.Builder{}
	fmt.Fprintf(s, "  %s examples\t\t# list the builtin examples\n", meta.Bin)
	fmt.Fprintf(s, "  %s info ascii\t\t# information on the buildin ascii sample\n", meta.Bin)
	fmt.Fprintf(s, "  %s info ascii -f json\t# information in json format on the ascii sample\n", meta.Bin)
	fmt.Fprintf(s, "  %s view ascii\t\t# view the ascii example\n", meta.Bin)
	fmt.Fprintf(s, "  %s info ansi.rgb\t# information on the 24-bit color ansi example\n", meta.Bin)
	fmt.Fprintf(s, "  %s view ansi.rgb\t# view the 24-bit color ansi example", meta.Bin)
	return s.String()
}

func table() string {
	s := &strings.Builder{}
	fmt.Fprintf(s, "  %s table cp437\n", meta.Bin)
	fmt.Fprintf(s, "  %s table cp437 latin1 windows-1252\n", meta.Bin)
	fmt.Fprintf(s, "  %s table iso-8859-15\n", meta.Bin)
	fmt.Fprintf(s, "  %s list\t# list the supported codepage tables\n", meta.Bin)
	return s.String()
}

func info() string {
	s := &strings.Builder{}
	fmt.Fprintf(s, "  %s info text.asc logo.jpg      # print the information of multiple files\n", meta.Bin)
	fmt.Fprintf(s, "  %s info file.txt --format=json # print the information using a structured syntax\n", meta.Bin)
	return s.String()
}

func view() string {
	s := &strings.Builder{}
	fmt.Fprintf(s, "  %s view file.txt -e latin1\n", meta.Bin)
	fmt.Fprintf(s, "  %s view file1.txt file2.txt --encode=\"iso-8859-1\"\n", meta.Bin)
	fmt.Fprintf(s, "  cat file.txt | %s view", meta.Bin)
	return s.String()
}
