package cmd

import (
	"fmt"
	"log"

	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/view"
	"github.com/spf13/cobra"
)

func ViewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     fmt.Sprintf("view %s", example.Filenames),
		Aliases: []string{"v"},
		GroupID: "fileCmds",
		Short:   "Print a text file to the terminal using standard output",
		Long:    "Print a text file to the terminal using standard output.",
		Example: fmt.Sprint(example.View),
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
	flag.Encode(&flag.View.Encode, vc)
	flag.Controls(&flag.View.Controls, vc)
	flag.SwapChars(&flag.View.Swap, vc)
	if err := flag.HiddenTo(&flag.View.To, vc); err != nil {
		log.Fatal(err)
	}
	flag.Width(&flag.View.Width, vc)
	vc.Flags().SortFlags = false
	return vc
}

func init() { //nolint:gochecknoinits
	Cmd.AddCommand(ViewInit())
}
