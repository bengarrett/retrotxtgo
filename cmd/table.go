package cmd

import (
	"strings"

	"github.com/bengarrett/retrotxtgo/cmd/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/list"
	"github.com/spf13/cobra"
)

func Table() *cobra.Command {
	s := "Display one or more code page tables showing all the characters in use"
	l := "Display one or more code page tables showing all the characters in use."
	expl := strings.Builder{}
	example.Table.String(&expl)
	return &cobra.Command{
		Use:     "table [code page names or aliases]",
		Aliases: []string{"t"},
		Short:   s,
		Long:    l,
		Example: expl.String(),
		GroupID: IDcodepage,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := flag.Help(cmd, args...); err != nil {
				return err
			}
			return list.Table(cmd.OutOrStdout(), args...)
		},
	}
}

func Tables() *cobra.Command {
	return &cobra.Command{
		Use:     "tables",
		Short:   "Display the characters of every code page table in use",
		Long:    "Display the characters of every code page table in use.",
		GroupID: IDcodepage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list.Tables(cmd.OutOrStdout())
		},
	}
}

func init() {
	Cmd.AddCommand(Table())
	Cmd.AddCommand(Tables())
}
