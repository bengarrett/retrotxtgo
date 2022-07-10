package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/infocmd"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/cobra"
)

func infoCommand() *cobra.Command {
	return &cobra.Command{
		Use:     fmt.Sprintf("info %s", example.Filenames),
		Aliases: []string{"i"},
		Short:   "Information on a text file",
		Long:    "Discover details and information about any text or text art file.",
		Example: example.Info.Print(),
		Run: func(cmd *cobra.Command, args []string) {
			if err := infocmd.Run(cmd, args); err != nil {
				logs.Fatal(err)
			}
		},
	}
}

//nolint:gochecknoinits
func init() {
	infoCmd := infoCommand()
	rootCmd.AddCommand(infoCmd)
	infos := config.Format().Info
	infoCmd.Flags().StringVarP(&flag.InfoFlag.Format, "format", "f", "color",
		str.Options("print format or syntax", true, true, infos[:]...))
}
