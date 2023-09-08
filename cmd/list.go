package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/pkg/logs"
	"github.com/spf13/cobra"
)

func ListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		GroupID: "listCmds",
		Short:   "Available codepages and tabled datasets",
		Long:    "List the available codepages and their tabled values.",
		Example: fmt.Sprint(example.List),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := flag.Help(cmd, args...); err != nil {
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
	lc.AddGroup(&cobra.Group{ID: "codepages", Title: "Codepages:"})
	lc.AddGroup(&cobra.Group{ID: "tables", Title: "Codepage Table:"})
	lc.AddCommand(ListCodepages())
	lc.AddCommand(ListTable())
	lc.AddCommand(ListTables())
	return lc
}

//nolint:gochecknoinits
func init() {
	Cmd.AddCommand(ListInit())
	Cmd.AddCommand(ListExamples())
}
