package cmd

import (
	"github.com/bengarrett/retrotxtgo/table"
	"github.com/spf13/cobra"
)

func Language() *cobra.Command {
	s := "List languages supported by legacy code pages"
	l := "List languages and writing alphabets supported by legacy code pages."
	return &cobra.Command{
		Use:     "lang",
		Aliases: []string{"la", "language"},
		Short:   s,
		Long:    l,
		GroupID: IDcodepage,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return table.ListLanguage(cmd.OutOrStdout())
		},
	}
}

func init() {
	Cmd.AddCommand(Language())
}
