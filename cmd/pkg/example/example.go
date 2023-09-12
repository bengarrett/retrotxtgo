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
	Cmd          Example = iota // Cmd is the example for the root command.
	List                        // List is the example for the list command.
	ListExamples                // ListExamples are the examples for the list examples command.
	ListTable                   // ListTable are the examples for the list tables command.
	Info                        // Info is the example for the info command.
	View                        // View is the example for the view command.
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
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), cmmt)
		l := len(s)
		if l < sentence {
			fmt.Fprintln(w, term.Info(scanner.Text()))
			continue
		}
		// do not the last hash as a comment
		ex := strings.Join(s[:l-1], cmmt)
		fmt.Fprint(w, term.Info(ex))
		fmt.Fprintf(w, "%s%s\n  ", color.Secondary.Sprint(cmmt), s[l-1])
	}
}

func (e Example) result() string {
	switch e {
	case Cmd:
		return cmd()
	case List:
		return list()
	case ListExamples:
		return listExamples()
	case ListTable:
		return listTable()
	case Info:
		return info()
	case View:
		return view()
	}
	return ""
}

func cmd() string {
	return fmt.Sprintf("  %s\n%s\n%s",
		fmt.Sprintf("%s info %s", meta.Bin, Filenames),
		fmt.Sprintf("%s view %s", meta.Bin, Filenames),
		fmt.Sprintf("%s %s      # print text files partial info TODO", meta.Bin, Filenames),
	)
}

func list() string {
	return fmt.Sprintf("  %s\n%s\n%s",
		fmt.Sprintf("%s list codepages", meta.Bin),
		fmt.Sprintf("%s list table cp437 cp1252", meta.Bin),
		fmt.Sprintf("%s list tables", meta.Bin))
}

func listExamples() string {
	return fmt.Sprintf("  %s\n%s\n%s",
		fmt.Sprintf("%s list examples # list the builtin examples", meta.Bin),
		fmt.Sprintf("%s info ascii    # information on the buildin ascii example", meta.Bin),
		fmt.Sprintf("%s view ascii    # view the ascii example", meta.Bin))
}

func listTable() string {
	return fmt.Sprintf("  %s\n%s\n%s",
		fmt.Sprintf("%s table cp437", meta.Bin),
		fmt.Sprintf("%s table cp437 latin1 windows-1252", meta.Bin),
		fmt.Sprintf("%s table iso-8859-15", meta.Bin))
}

func info() string {
	return fmt.Sprintf("  %s %s\n%s %s",
		meta.Bin, "info text.asc logo.jpg # print the information of multiple files",
		meta.Bin, "info file.txt --format=json # print the information using a structured syntax")
}

func view() string {
	return fmt.Sprintf("  %s\n%s\n%s",
		fmt.Sprintf("%s view file.txt -e latin1", meta.Bin),
		fmt.Sprintf("%s view file1.txt file2.txt --encode=\"iso-8859-1\"", meta.Bin),
		fmt.Sprintf("cat file.txt | %s view", meta.Bin))
}
