// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/viewcmd"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command.
var viewCmd = &cobra.Command{
	Use:     fmt.Sprintf("view %s", filenames),
	Aliases: []string{"v"},
	Short:   "Print a text file to the terminal using standard output",
	Long:    "Print a text file to the terminal using standard output.",
	Example: example.Print(example.View),
	Run: func(cmd *cobra.Command, args []string) {
		b, err := viewcmd.Run(cmd, args...)
		if err != nil {
			logs.Fatal(err)
		}
		fmt.Print(b)
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
	flagEncode(&flag.ViewFlag.Encode, viewCmd)
	flagControls(&flag.ViewFlag.Controls, viewCmd)
	flagRunes(&flag.ViewFlag.Swap, viewCmd)
	flagTo(&flag.ViewFlag.To, viewCmd)
	flagWidth(&flag.ViewFlag.Width, viewCmd)
	viewCmd.Flags().SortFlags = false
}
