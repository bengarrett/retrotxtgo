package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/info"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/cobra"
)

func InfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:     fmt.Sprintf("info %s", example.Filenames),
		Aliases: []string{"i"},
		Short:   "Information on a text file",
		Long:    "Discover details and information about any text or text art file.",
		Example: fmt.Sprint(example.Info),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := info.Run(cmd, args); err != nil {
				return err
			}
			return nil
		},
	}
}

func InfoInit() *cobra.Command {
	ic := InfoCommand()
	infos := config.Format().Info
	ic.Flags().StringVarP(&flag.InfoFlag.Format, "format", "f", "color",
		str.Options("print format or syntax", true, true, infos[:]...))
	return ic
}

//nolint:gochecknoinits
func init() {
	Cmd.AddCommand(InfoInit())
}
