package cmd

import (
	"github.com/bengarrett/retrotxtgo/pkg/table"
	"github.com/spf13/cobra"
)

func Language() *cobra.Command {
	s := "List the legacy codepage target languages"
	l := "List the legacy codepage target languages."
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

//nolint:gochecknoinits
func init() {
	Cmd.AddCommand(Language())
}
