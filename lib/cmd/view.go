package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/viewcmd"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
)

func viewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     fmt.Sprintf("view %s", example.Filenames),
		Aliases: []string{"v"},
		Short:   "Print a text file to the terminal using standard output",
		Long:    "Print a text file to the terminal using standard output.",
		Example: example.View.Print(),
		Run: func(cmd *cobra.Command, args []string) {
			b, err := viewcmd.Run(cmd, args...)
			if err != nil {
				logs.Fatal(err)
			}
			fmt.Print(b)
		},
	}
}

func init() { //nolint:gochecknoinits
	viewCmd := viewCommand()
	rootCmd.AddCommand(viewCmd)
	flag.Encode(&flag.ViewFlag.Encode, viewCmd)
	flag.Controls(&flag.ViewFlag.Controls, viewCmd)
	flag.Runes(&flag.ViewFlag.Swap, viewCmd)
	flag.To(&flag.ViewFlag.To, viewCmd)
	flag.Width(&flag.ViewFlag.Width, viewCmd)
	viewCmd.Flags().SortFlags = false
}
