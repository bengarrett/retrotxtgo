package cmd

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/ianaindex"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Codepages and tabled datasets available",
	Example: "  retrotxt list codepages\n  retrotxt list table cp437  \n  retrotxt list tables",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Usage()
			logs.Check("list.usage", err)
			os.Exit(0)
		}
	},
}

var listCmdCodepages = &cobra.Command{
	Use:   "codepages",
	Short: "list available legacy codepages that RetroTxt can convert into UTF-8",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(convert.List())
	},
}

var listCmdTable = &cobra.Command{
	Use:     "table",
	Short:   "display a table showing the codepage and all its characters",
	Example: "  retrotxt table cp437",
	Run: func(cmd *cobra.Command, args []string) {
		table, err := convert.Table(viewArgs.cp)
		logs.ChkErr(logs.Err{Issue: "table", Arg: viewArgs.cp, Msg: err})
		fmt.Println(table.String())
	},
}

var listCmdTables = &cobra.Command{
	Use:   "tables",
	Short: "display tables showing known codepages and characters",
	Run: func(cmd *cobra.Command, args []string) {
		for _, e := range convert.Encodings() {
			name, err := ianaindex.MIME.Name(e)
			if err != nil {
				logs.Log(err)
			} else {
				// keep 0F,1F controls. blank other ?
				// tables -> Macintosh to list alt. names Mac OS Roman
				// Windows 874 is not showing different chars from ISO-11
				// https://en.wikipedia.org/wiki/ISO/IEC_8859-11#Vendor_extensions
				// japanese needs fixing
				table, err := convert.Table(name)
				logs.ChkErr(logs.Err{Issue: "tables", Arg: name, Msg: err})
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
	// table cmd
	listCmd.AddCommand(listCmdTable)
	listCmdTable.Flags().StringVarP(&viewArgs.cp, "codepage", "c", "cp437",
		"legacy character encoding table to display")
	// tables cmd
	listCmd.AddCommand(listCmdTables)
}
