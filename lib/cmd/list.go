package cmd

import (
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/listcmd"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Available inbuilt examples, codepages and tabled datasets",
	Long:    "List the available inbuilt text art and text documents, codepages and their tabled values.",
	Example: example.List.Print(),
	Run: func(cmd *cobra.Command, args []string) {
		if err := flag.PrintUsage(cmd, args...); err != nil {
			logs.Fatal(err)
		}
		logs.FatalCmd("list", args...)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listcmd.Codepages.Command())
	listCmd.AddCommand(listcmd.Examples.Command())
	listCmd.AddCommand(listcmd.Table.Command())
	listCmd.AddCommand(listcmd.Tables.Command())
}
