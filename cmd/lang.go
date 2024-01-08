package cmd

import (
	"github.com/bengarrett/retrotxtgo/table"
	"github.com/spf13/cobra"
)

func Language() *cobra.Command {
	s := "List the natural languages of legacy code pages"
	l := "List the natural languages and writing alphabets of legacy code pages."
	return &cobra.Command{
		Use:     "lang",
		Aliases: []string{"la", "language"},
		Short:   s,
		Long:    l,
		GroupID: IDcodepage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return table.ListLanguage(cmd.OutOrStdout())
		},
	}
}

func init() {
	Cmd.AddCommand(Language())
}
