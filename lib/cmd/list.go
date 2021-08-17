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

var listExample = fmt.Sprintf("  %s list codepages\n  %s list examples\n  %s list table cp437 cp1252 \n  %s list tables",
	meta.Bin, meta.Bin, meta.Bin, meta.Bin)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Available built-in examples, codepages and tabled datasets",
	Example: exampleCmd(listExample),
	Run: func(cmd *cobra.Command, args []string) {
		if !printUsage(cmd, args...) {
			logs.InvalidCommand("list", args...)
		}
	},
}

var listCmdCodepages = &cobra.Command{
	Use:     "codepages",
	Aliases: []string{"c", "cp"},
	Short:   fmt.Sprintf("List available legacy codepages that %s can convert into UTF-8", meta.Name),
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
	const padding = 2
	w := tabwriter.NewWriter(&buf, 0, 0, padding, ' ', 0)
	bin := fmt.Sprintf("  %s ", meta.Bin)
	title := fmt.Sprintf("\n Packaged example text and ANSI files to test and play with %s ", meta.Name)
	fmt.Fprintln(w, str.Cp(title))
	fmt.Fprintln(w, strings.Repeat("-", len(title)))
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
		logs.ProblemFatal(logs.ErrTabFlush, err)
	}
	return &buf
}

var listTableExample = fmt.Sprintf("  %s table cp437\n  %s table cp437 latin1 windows-1252\n  %s table iso-8859-15",
	meta.Bin, meta.Bin, meta.Bin)

var listCmdTable = &cobra.Command{
	Use:     "table [codepage names or aliases]",
	Aliases: []string{"t"},
	Short:   "Display one or more tables showing the codepage and all their characters",
	Example: exampleCmd(listTableExample),
	Run: func(cmd *cobra.Command, args []string) {
		if !printUsage(cmd, args...) {
			fmt.Println(listTable(args))
		}
	},
}

func listTable(args []string) (s string) {
	for _, arg := range args {
		table, err := convert.Table(arg)
		if err != nil {
			logs.ProblemMark(arg, ErrTable, err)
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
				logs.ProblemMark(fmt.Sprint(e), ErrIANA, err)
				continue
			}
		}
		if skipTable(name) {
			continue
		}
		table, err := convert.Table(name)
		if err != nil {
			logs.ProblemMark(name, ErrTable, err)
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
