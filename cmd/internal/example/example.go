package example

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"text/template"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
)

const Filenames = "[filenames]"

type Example int

const (
	Cmd Example = iota
	Config
	ConfigInfo
	Create
	List
	ListExamples
	ListTable
	Info
	Set
	View
)

// Print returns help usage examples.
func (e Example) String() string {
	var b bytes.Buffer
	// change example operating system path separator
	t := template.Must(template.New("example").Parse(e.result()))
	err := t.Execute(&b, string(os.PathSeparator))
	if err != nil {
		log.Fatal(err)
	}
	// color the example text except text following
	// the last hash #, which is treated as a comment
	const cmmt, sentence = "#", 2
	scanner, s := bufio.NewScanner(&b), ""
	for scanner.Scan() {
		ss := strings.Split(scanner.Text(), cmmt)
		l := len(ss)
		if l < sentence {
			s += str.ColInf(scanner.Text()) + "\n  "
			continue
		}
		// do not the last hash as a comment
		ex := strings.Join(ss[:l-1], cmmt)
		s += str.ColInf(ex)
		s += fmt.Sprintf("%s%s\n  ", color.Secondary.Sprint(cmmt), ss[l-1])
	}
	return strings.TrimSpace(s)
}

func (e Example) result() string {
	switch e {
	case Cmd:
		return cmd()
	case Config:
		return config()
	case ConfigInfo:
		return configInfo()
	case Create:
		return create()
	case List:
		return list()
	case ListExamples:
		return listExamples()
	case ListTable:
		return listTable()
	case Info:
		return info()
	case Set:
		return set()
	case View:
		return view()
	}
	return ""
}

func cmd() string {
	return fmt.Sprintf("  %s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
		"# save the text files as webpages",
		fmt.Sprintf("%s create %s", meta.Bin, Filenames),
		"# save the text files as webpages stored in a zip file",
		fmt.Sprintf("%s create %s --compress", meta.Bin, Filenames),
		"# print detailed information about the text files",
		fmt.Sprintf("%s info   %s", meta.Bin, Filenames),
		"# print the text files as Unicode text",
		fmt.Sprintf("%s view   %s", meta.Bin, Filenames),
		fmt.Sprintf("# configure the %s flags and settings", meta.Name),
		fmt.Sprintf("%s config setup", meta.Bin),
	)
}

func config() string {
	return fmt.Sprintf("  %s %s %s\n%s %s %s",
		meta.Bin, "config setup", "# Walk through all the settings",
		meta.Bin, "config set --list", "# List all the settings in use")
}

func configInfo() string {
	return fmt.Sprintf("  %s\n%s",
		fmt.Sprintf("%s config info   # List the default setting values", meta.Bin),
		fmt.Sprintf("%s config set -c # List the settings and help hints", meta.Bin))
}

func create() string {
	return fmt.Sprintf("  %s%s\n%s%s\n%s%s\n%s%s\n%s%s",
		"# print a HTML file created from file.txt\n",
		fmt.Sprintf("%s create file.txt --title \"A text file\" --description \"Some text goes here\"", meta.Bin),
		"# save HTML files created from file1.txt and file2.asc\n",
		fmt.Sprintf("%s create file1.txt file2.asc --save", meta.Bin),
		"# save and compress a HTML file created from file.txt in Downloads.\n",
		fmt.Sprintf("%s create ~{{.}}Downloads{{.}}file.txt --compress", meta.Bin),
		"# host the HTML file created from file.txt\n",
		fmt.Sprintf("%s create file.txt --serve=%d", meta.Bin, meta.WebPort),
		"# pipe a HTML file created from file.txt\n",
		fmt.Sprintf("%s create file.txt | %s", meta.Bin, cat()))
}

// cat returns the os command name to concatenate a file to standard output.
func cat() string {
	if runtime.GOOS == "windows" {
		return "type"
	}
	return "cat"
}

func list() string {
	return fmt.Sprintf("  %s\n%s\n%s\n%s",
		fmt.Sprintf("%s list codepages", meta.Bin),
		fmt.Sprintf("%s list examples", meta.Bin),
		fmt.Sprintf("%s list table cp437 cp1252", meta.Bin),
		fmt.Sprintf("%s list tables", meta.Bin))
}

func listExamples() string {
	return fmt.Sprintf("  %s\n%s\n%s\n%s\n%s",
		fmt.Sprintf("%s list examples # list the builtin examples", meta.Bin),
		fmt.Sprintf("%s info ascii    # information on the buildin ascii example", meta.Bin),
		fmt.Sprintf("%s view ascii    # view the ascii example", meta.Bin),
		fmt.Sprintf("%s create ascii  # create the ascii example", meta.Bin),
		fmt.Sprintf("%s save ascii    # save the ascii example", meta.Bin))
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

func set() string {
	return fmt.Sprintf("  %s %s %s\n%s %s %s\n%s %s %s",
		meta.Bin, "config set --list", "# List the available settings",
		meta.Bin, "config set html.meta.description", "# Edit the meta description setting",
		meta.Bin, "config set style.info style.html", fmt.Sprintf("# Edit both the %s color styles", meta.Name),
	)
}

func view() string {
	return fmt.Sprintf("  %s\n%s\n%s",
		fmt.Sprintf("%s view file.txt -e latin1", meta.Bin),
		fmt.Sprintf("%s view file1.txt file2.txt --encode=\"iso-8859-1\"", meta.Bin),
		fmt.Sprintf("cat file.txt | %s view", meta.Bin))
}
