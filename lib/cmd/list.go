// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/sample"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
)

var listExample = fmt.Sprintf("  %s\n%s\n%s\n%s",
	fmt.Sprintf("%s list codepages", meta.Bin),
	fmt.Sprintf("%s list examples", meta.Bin),
	fmt.Sprintf("%s list table cp437 cp1252", meta.Bin),
	fmt.Sprintf("%s list tables", meta.Bin))

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Available inbuilt examples, codepages and tabled datasets",
	Long:    "List the available inbuilt text art and text documents, codepages and their tabled values.",
	Example: exampleCmd(listExample),
	Run: func(cmd *cobra.Command, args []string) {
		if !printUsage(cmd, args...) {
			logs.FatalCmd("list", args...)
		}
	},
}

var listCmdCodepages = &cobra.Command{
	Use:     "codepages",
	Aliases: []string{"c", "cp"},
	Short: fmt.Sprintf("List the legacy codepages that %s can convert to UTF-8",
		meta.Name),
	Long: fmt.Sprintf("List the available legacy codepages that %s can convert to UTF-8.",
		meta.Name),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(convert.List())
	},
}

var listExpExample = fmt.Sprintf("  %s\n%s\n%s\n%s\n%s",
	fmt.Sprintf("%s list examples # list the builtin examples", meta.Bin),
	fmt.Sprintf("%s info ascii    # information on the buildin ascii example", meta.Bin),
	fmt.Sprintf("%s view ascii    # view the ascii example", meta.Bin),
	fmt.Sprintf("%s create ascii  # create the ascii example", meta.Bin),
	fmt.Sprintf("%s save ascii    # save the ascii example", meta.Bin))

var listCmdExamples = &cobra.Command{
	Use:     "examples",
	Aliases: []string{"e"},
	Short: fmt.Sprintf("List builtin text files available for use with the %s, %s, %s and %s commands",
		str.Example("create"), str.Example("save"), str.Example("info"), str.Example("view")),
	Long: fmt.Sprintf("List builtin text art and documents available for use with the %s, %s, %s and %s commands.",
		str.Example("create"), str.Example("save"), str.Example("info"), str.Example("view")),
	Example: exampleCmd(listExpExample),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(examples())
	},
}

func examples() *bytes.Buffer {
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

var listTableExample = fmt.Sprintf("  %s\n%s\n%s",
	fmt.Sprintf("%s table cp437", meta.Bin),
	fmt.Sprintf("%s table cp437 latin1 windows-1252", meta.Bin),
	fmt.Sprintf("%s table iso-8859-15", meta.Bin))

var listCmdTable = &cobra.Command{
	Use:     "table [codepage names or aliases]",
	Aliases: []string{"t"},
	Short:   "Display one or more codepage tables showing all the characters in use",
	Long:    "Display one or more codepage tables showing all the characters in use.",
	Example: exampleCmd(listTableExample),
	Run: func(cmd *cobra.Command, args []string) {
		if !printUsage(cmd, args...) {
			fmt.Print(listTable(args...))
		}
	},
}

// listTable returns one or more named encodings in a tabled format.
func listTable(names ...string) (s string) {
	for _, name := range names {
		table, err := convert.Table(name)
		if err != nil {
			fmt.Println(logs.PrintfMark(name, ErrTable, err))
			continue
		}
		s = fmt.Sprintf("%s%s", s, table.String())
	}
	return s
}

var listCmdTables = &cobra.Command{
	Use:   "tables",
	Short: "Display the characters of every codepage table inuse",
	Long:  "Display the characters of every codepage table inuse.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(listTbls())
	},
}

// listTbls returns all the supported encodings in a tabled format.
func listTbls() (s string) {
	for _, e := range convert.Encodings() {
		var (
			err  error
			name string
		)
		if e == charmap.XUserDefined {
			name = "iso-8859-11"
		}
		if name == "" {
			name, err = ianaindex.MIME.Name(e)
			if err != nil {
				fmt.Println(logs.PrintfMark(fmt.Sprint(e), ErrIANA, err))
				continue
			}
		}
		if !usableTbl(name) {
			continue
		}
		table, err := convert.Table(name)
		if err != nil {
			fmt.Println(logs.PrintfMark(name, ErrTable, err))
			continue
		}
		s = fmt.Sprintf("%s%s", s, table.String())
	}
	return s
}

// usableTbl returns true if the named encoding be shown in an 8-bit table.
func usableTbl(name string) bool {
	switch strings.ToLower(name) {
	case "", "utf-16", "utf-16be", "utf-16le", "utf-32", "utf-32be", "utf-32le":
		return false
	}
	return true
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listCmdCodepages)
	listCmd.AddCommand(listCmdExamples)
	listCmd.AddCommand(listCmdTable)
	listCmd.AddCommand(listCmdTables)
}
