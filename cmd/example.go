package cmd

import (
	"fmt"
	"strings"

	"github.com/bengarrett/retrotxtgo/cmd/pkg/example"
	"github.com/bengarrett/retrotxtgo/cmd/pkg/list"
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/spf13/cobra"
)

func ListExample() *cobra.Command {
	s := fmt.Sprintf("List the included sample text files available for use with the %s and %s commands",
		term.Example("info"), term.Example("view"))
	l := fmt.Sprintf("List the included sample text art and documents available for use with the %s and %s commands.",
		term.Example("info"), term.Example("view"))
	expl := strings.Builder{}
	example.ListExamples.String(&expl)
	return &cobra.Command{
		Use:     "example",
		Aliases: []string{"e", "sample", "s"},
		GroupID: IDsample,
		Short:   s,
		Long:    l,
		Example: expl.String(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return list.Examples(cmd.OutOrStdout())
		},
	}
}

func init() {
	Cmd.AddCommand(ListExample())
}
