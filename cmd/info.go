package cmd

import (
	"fmt"
	"strings"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/format"
	"github.com/bengarrett/retrotxtgo/cmd/internal/info"
	"github.com/bengarrett/retrotxtgo/cmd/pkg/example"
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/spf13/cobra"
)

func InfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:     fmt.Sprintf("info %s", example.Filenames),
		Aliases: []string{"i"},
		GroupID: "fileCmds",
		Short:   "Information on a text file",
		Long:    "Discover details and information about any text or text art file.",
		Example: fmt.Sprint(example.Info),
		RunE: func(cmd *cobra.Command, args []string) error {
			return info.Run(cmd.OutOrStdout(), cmd, args...)
		},
	}
}

func InfoInit() *cobra.Command {
	ic := InfoCommand()
	infos := format.Format().Info
	s := &strings.Builder{}
	_, _ = term.Options(s, "print format or syntax", true, true, infos[:]...)
	ic.Flags().StringVarP(&flag.Info.Format, "format", "f", "color", s.String())
	return ic
}

//nolint:gochecknoinits
func init() {
	Cmd.AddCommand(InfoInit())
}
