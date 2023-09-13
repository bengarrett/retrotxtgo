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
	s := "Information on a text file"
	l := "Discover details and information about any text or text art file."
	expl := strings.Builder{}
	example.Info.String(&expl)
	return &cobra.Command{
		Use:     fmt.Sprintf("info %s", example.Filenames),
		Aliases: []string{"i"},
		GroupID: IDfile,
		Short:   s,
		Long:    l,
		Example: expl.String(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return info.Run(cmd.OutOrStdout(), cmd, args...)
		},
	}
}

func InfoInit() *cobra.Command {
	infoc := InfoCommand()
	infos := format.Format().Info
	s := &strings.Builder{}
	term.Options(s, "print format or syntax", true, true, infos[:]...)
	infoc.Flags().StringVarP(&flag.Info.Format, "format", "f", "color", s.String())
	return infoc
}

//nolint:gochecknoinits
func init() {
	Cmd.AddCommand(InfoInit())
}
