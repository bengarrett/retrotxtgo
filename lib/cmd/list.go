package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/ianaindex"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Available codepages and tabled datasets",
	Example: "  retrotxt list codepages\n  retrotxt list table cp437  \n  retrotxt list tables",
	Run: func(cmd *cobra.Command, args []string) {
		checkUse(cmd, args)
	},
}

var listCmdCodepages = &cobra.Command{
	Use:     "codepages",
	Aliases: []string{"c", "cp"},
	Short:   "list available legacy codepages that RetroTxt can convert into UTF-8",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(convert.List())
	},
}

var listCmdExamples = &cobra.Command{
	Use:     "examples",
	Aliases: []string{"e"},
	Short:   "list pre-packaged text files for use with the create, info and view commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(examples())
	},
}

var listCmdTable = &cobra.Command{
	Use:     "table [codepage names or aliases]",
	Aliases: []string{"t"},
	Short:   "display one or more tables showing the codepage and all their characters",
	Example: "  retrotxt table cp437\n  retrotxt table cp437 latin1 windows-1252\n  retrotxt table iso-8859-15",
	Run: func(cmd *cobra.Command, args []string) {
		checkUse(cmd, args)
		for _, arg := range args {
			table, err := convert.Table(arg)
			if ok := logs.CheckCont("list.table", err); ok {
				fmt.Println(table.String())
			}
		}
	},
}

var listCmdTables = &cobra.Command{
	Use:   "tables",
	Short: "display tables showing known codepages and characters",
	Run: func(cmd *cobra.Command, args []string) {
		for _, e := range convert.Encodings() {
			name, err := ianaindex.MIME.Name(e)
			if ok := logs.CheckCont("list.tables.ianaindex", err); !ok {
				logs.LogCont(err)
				continue
			}
			// keep 0F,1F controls. blank other ?
			// tables -> Macintosh to list alt. names Mac OS Roman
			// Windows 874 is not showing different chars from ISO-11
			// https://en.wikipedia.org/wiki/ISO/IEC_8859-11#Vendor_extensions
			// japanese needs fixing
			table, err := convert.Table(name)
			if ok := logs.CheckCont("list.tables", err); ok {
				fmt.Println(table.String())
			}
		}
	},
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
