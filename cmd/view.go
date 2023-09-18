package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/view"
	"github.com/bengarrett/retrotxtgo/cmd/pkg/example"
	"github.com/spf13/cobra"
)

func ViewCommand() *cobra.Command {
	s := "Print a text file to the terminal using standard output"
	l := "Print a text file to the terminal using standard output."
	expl := strings.Builder{}
	example.View.String(&expl)
	return &cobra.Command{
		Use:     fmt.Sprintf("view %s", example.Filenames),
		Aliases: []string{"v"},
		GroupID: IDfile,
		Short:   s,
		Long:    l,
		Example: expl.String(),
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

func init() {
	Cmd.AddCommand(ViewInit())
}
