package cmd

import (
	"strings"

	"github.com/bengarrett/retrotxtgo/cmd/example"
	"github.com/bengarrett/retrotxtgo/cmd/list"
	"github.com/spf13/cobra"
)

func ListExample() *cobra.Command {
	s := "Browse and view built-in sample files"
	l := "Browse and view built-in sample text art and documents."
	expl := strings.Builder{}
	example.Examples.String(&expl)
	return &cobra.Command{
		Use:     "example",
		Aliases: []string{"e", "sample", "s"},
		GroupID: IDsample,
		Short:   s,
		Long:    l,
		Example: expl.String(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return list.Examples(cmd.OutOrStdout())
		},
	}
}

func init() {
	Cmd.AddCommand(ListExample())
}
