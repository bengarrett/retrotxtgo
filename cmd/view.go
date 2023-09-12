package cmd

import (
	"fmt"
	"log"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/view"
	"github.com/bengarrett/retrotxtgo/cmd/pkg/example"
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
			return view.Run(cmd.OutOrStdout(), cmd, args...)
		},
	}
}

func ViewInit() *cobra.Command {
	vc := ViewCommand()
	f := flag.View()
	flag.Encode(&f.Encode, vc)
	flag.Controls(&f.Controls, vc)
	flag.SwapChars(&f.Swap, vc)
	if err := flag.HiddenTo(&f.To, vc); err != nil {
		log.Fatal(err)
	}
	flag.Width(&f.Width, vc)
	vc.Flags().SortFlags = false
	return vc
}

func init() { //nolint:gochecknoinits
	Cmd.AddCommand(ViewInit())
}
