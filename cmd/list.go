package cmd

import (
	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/list"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
)

func ListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "Available inbuilt examples, codepages and tabled datasets",
		Long:    "List the available inbuilt text art and text documents, codepages and their tabled values.",
		Example: example.List.Print(),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := flag.PrintUsage(cmd, args...); err != nil {
				return err
			}
			if len(args) > 0 {
				logs.FatalCmd("list", args...) // TODO replace with a print and error return
			}
			return nil
		},
	}
}

func ListInit() *cobra.Command {
	lc := ListCommand()
	lc.AddCommand(list.Codepages.Command())
	lc.AddCommand(list.Examples.Command())
	lc.AddCommand(list.Table.Command())
	lc.AddCommand(list.Tables.Command())
	return lc
}

//nolint:gochecknoinits
func init() {
	Cmd.AddCommand(ListInit())
}
