package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/pkg/table"
	"github.com/spf13/cobra"
)

func ListCodepage() *cobra.Command {
	s := fmt.Sprintf("List the legacy codepages that %s can convert to UTF-8", meta.Name)
	l := fmt.Sprintf("List the available legacy codepages that %s can convert to UTF-8.", meta.Name)
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "cp", "codepage"},
		Short:   s,
		Long:    l,
		GroupID: IDcodepage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return table.List(cmd.OutOrStdout())
		},
	}
}

func init() {
	Cmd.AddCommand(ListCodepage())
}
