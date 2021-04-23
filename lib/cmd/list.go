// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
	"retrotxt.com/retrotxt/lib/convert"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/sample"
	"retrotxt.com/retrotxt/lib/str"
)

const listExample = "  retrotxt list codepages\n  retrotxt list examples\n  retrotxt list table cp437 cp1252 \n  retrotxt list tables"

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Available built-in examples, codepages and tabled datasets",
	Example: exampleCmd(listExample),
	Run: func(cmd *cobra.Command, args []string) {
		checkUse(cmd, args...)
		logs.ArgFatal(args...)
	},
}

var listCmdCodepages = &cobra.Command{
	Use:     "codepages",
	Aliases: []string{"c", "cp"},
	Short:   "List available legacy codepages that RetroTxt can convert into UTF-8",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(convert.List())
	},
}

var listCmdExamples = &cobra.Command{
	Use:     "examples",
	Aliases: []string{"e"},
	Short: "List pre-packaged text files for use with the " +
		str.Example("create") + ", " + str.Example("save") + ", " + str.Example("info") + " and " + str.Example("view") + " commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(examples())
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
	var flags uint = 0 // tabwriter.AlignRight | tabwriter.Debug
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', flags)
	const title = "\n Packaged example text and ANSI files to test and play with RetroTxt "
	fmt.Fprintln(w, str.Cp(title))
	fmt.Fprintln(w, strings.Repeat("-", len(title)))
	for _, k := range keys {
		fmt.Fprintf(w, "%s\t%s\t\n", k, m[k].Description)
	}
	fmt.Fprintln(w, "\nAny of these packaged examples will work with the",
		str.Example("create")+",", str.Example("info"), "and", str.Example("view"), "commands.")
	fmt.Fprintln(w, "\nPrint the Windows-1252 English test to the terminal.\n"+str.Example("  retrotxt view 1252"))
	fmt.Fprintln(w, "\nConvert the Windows-1252 English test to UTF-8 encoding and save it to a file.\n"+
		str.Example("  retrotxt view 1252 > file.txt"))
	fmt.Fprintln(w, "\nSave the Windows-1252 English test with its original encoding.\n"+
		str.Example("  retrotxt view --to=cp1252 1252 > file.txt"))
	fmt.Fprintln(w, "\nDisplay statistics and information from a piped source.\n"+
		str.Example(" retrotxt view --to=cp1252 1252 | retrotxt info"))
	fmt.Fprintln(w, "\nDisplay statistics and information from the Windows-1252 English test.\n"+str.Example("  retrotxt info 1252"))
	fmt.Fprintln(w, "\nDisplay statistics, information and SAUCE metadata from the SAUCE test.\n"+str.Example("  retrotxt info sauce"))
	fmt.Fprintln(w, "\nCreate and display a HTML document from the Windows-1252 English test.\n"+str.Example("  retrotxt create 1252"))
	fmt.Fprintln(w, "\nCreate and save the HTML and assets from the Windows-1252 English test.\n"+str.Example("  retrotxt create 1252 --save"))
	fmt.Fprintln(w, "\nServe the Windows-1252 English test over a local web server.\n"+str.Example("  retrotxt create 1252 -p0"))
	fmt.Fprintln(w, "\nMultiple examples used together are supported.")
	fmt.Fprintln(w, str.Example("  retrotxt view ansi ascii ansi.rgb"))
	if err := w.Flush(); err != nil {
		logs.ProblemFatal(logs.ErrTabFlush, err)
	}
	return &buf
}

const listTableExample = "  retrotxt table cp437\n  retrotxt table cp437 latin1 windows-1252\n  retrotxt table iso-8859-15"

var listCmdTable = &cobra.Command{
	Use:     "table [codepage names or aliases]",
	Aliases: []string{"t"},
	Short:   "Display one or more tables showing the codepage and all their characters",
	Example: exampleCmd(listTableExample),
	Run: func(cmd *cobra.Command, args []string) {
		checkUse(cmd, args...)
		fmt.Println(listTable(args))
	},
}

func listTable(args []string) (s string) {
	for _, arg := range args {
		table, err := convert.Table(arg)
		if err != nil {
			logs.Println("list.table", "", err)
			continue
		}
		s += fmt.Sprintln(table.String())
	}
	return s
}

var listCmdTables = &cobra.Command{
	Use:   "tables",
	Short: "Display tables showing known codepages and characters",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(listAllTables())
	},
}

func listAllTables() (s string) {
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
				logs.Println("list.tables.ianaindex", "", err)
				continue
			}
		}
		if skipTable(name) {
			continue
		}
		table, err := convert.Table(name)
		if err != nil {
			logs.Println("list.tables", "", err)
			continue
		}
		s += fmt.Sprintln(table.String())
	}
	return s
}

func init() {
	// list cmd
	rootCmd.AddCommand(listCmd)
	// codepages cmd
	listCmd.AddCommand(listCmdCodepages)
	// examples cmd
	listCmd.AddCommand(listCmdExamples)
	// table cmd
	listCmd.AddCommand(listCmdTable)
	// tables cmd
	listCmd.AddCommand(listCmdTables)
}

// skipTable ignores encodings that cannot be correctly shown in an 8-bit table.
func skipTable(name string) bool {
	switch strings.ToLower(name) {
	case "utf-16", "utf-16be", "utf-16le", "utf-32", "utf-32be", "utf-32le":
		return true
	}
	return false
}
