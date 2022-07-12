package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/view"
	"github.com/spf13/cobra"
)

func ViewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     fmt.Sprintf("view %s", example.Filenames),
		Aliases: []string{"v"},
		Short:   "Print a text file to the terminal using standard output",
		Long:    "Print a text file to the terminal using standard output.",
		Example: example.View.Print(),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := view.Run(cmd, args...)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), b)
			return nil
		},
	}
}

func ViewInit() *cobra.Command {
	vc := ViewCommand()
	flag.Encode(&flag.ViewFlag.Encode, vc)
	flag.Controls(&flag.ViewFlag.Controls, vc)
	flag.Runes(&flag.ViewFlag.Swap, vc)
	flag.To(&flag.ViewFlag.To, vc)
	flag.Width(&flag.ViewFlag.Width, vc)
	vc.Flags().SortFlags = false
	return vc
}

func init() { //nolint:gochecknoinits
	Cmd.AddCommand(ViewInit())
}
