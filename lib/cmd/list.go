package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"

	"retrotxt.com/retrotxt/lib/convert"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Available built-in examples, codepages and tabled datasets",
	Example: "  retrotxt list codepages\n  retrotxt list examples\n  retrotxt list table cp437 cp1252 \n  retrotxt list tables",
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

var listCmdTable = &cobra.Command{
	Use:     "table [codepage names or aliases]",
	Aliases: []string{"t"},
	Short:   "Display one or more tables showing the codepage and all their characters",
	Example: "  retrotxt table cp437\n  retrotxt table cp437 latin1 windows-1252\n  retrotxt table iso-8859-15",
	Run: func(cmd *cobra.Command, args []string) {
		checkUse(cmd, args...)
		for _, arg := range args {
			table, err := convert.Table(arg)
			if err != nil {
				logs.Println("list.table", "", err)
				continue
			}
			fmt.Println(table.String())
		}
	},
}

var listCmdTables = &cobra.Command{
	Use:   "tables",
	Short: "Display tables showing known codepages and characters",
	Run: func(cmd *cobra.Command, args []string) {
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
			switch strings.ToLower(name) {
			case "utf-16", "utf-16be", "utf-16le", "utf-32", "utf-32be", "utf-32le":
				continue
			}
			table, err := convert.Table(name)
			if err != nil {
				logs.Println("list.tables", "", err)
				continue
			}
			fmt.Println(table.String())
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
