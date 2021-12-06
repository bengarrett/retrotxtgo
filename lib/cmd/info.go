// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/infocmd"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command.
var infoCmd = &cobra.Command{
	Use:     fmt.Sprintf("info %s", filenames),
	Aliases: []string{"i"},
	Short:   "Information on a text file",
	Long:    "Discover details and information about any text or text art file.",
	Example: example.Print(example.Info),
	Run:     infocmd.Run,
	// todo either stick with this or follow view.go where the error handler is placed in here.
}

func init() {
	rootCmd.AddCommand(infoCmd)
	i := config.Format().Info
	infoCmd.Flags().StringVarP(&flag.InfoFlag.Format, "format", "f", "color",
		str.Options("print format or syntax", true, true, i[:]...))
}
