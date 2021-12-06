// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/infocmd"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command.
var infoCmd = &cobra.Command{
	Use:     fmt.Sprintf("info %s", example.Filenames),
	Aliases: []string{"i"},
	Short:   "Information on a text file",
	Long:    "Discover details and information about any text or text art file.",
	Example: example.Print(example.Info),
	Run: func(cmd *cobra.Command, args []string) {
		if err := infocmd.Run(cmd, args); err != nil {
			logs.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infos := config.Format().Info
	infoCmd.Flags().StringVarP(&flag.InfoFlag.Format, "format", "f", "color",
		str.Options("print format or syntax", true, true, infos[:]...))
}
