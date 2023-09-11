package cmd

import (
	"bytes"
	"fmt"

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
		RunE:    info.Run,
	}
}

func InfoInit() *cobra.Command {
	ic := InfoCommand()
	infos := format.Format().Info
	b := &bytes.Buffer{}
	_, _ = term.Options(b, "print format or syntax", true, true, infos[:]...)
	ic.Flags().StringVarP(&flag.Info.Format, "format", "f", "color", b.String())
	return ic
}

//nolint:gochecknoinits
func init() {
	Cmd.AddCommand(InfoInit())
}
