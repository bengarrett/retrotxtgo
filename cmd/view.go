package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/view"
	"github.com/bengarrett/retrotxtgo/cmd/pkg/example"
	"github.com/spf13/cobra"
)

var viewLong = `Print a text file to the terminal using standard output.

Any texts and documents encoded in UTF-8 or that only use the common 
characters found in the 7-bit (128 characters) ASCII set are printed 
to the terminal as is. However, many old texts and documents are
encoded with legacy 8-bit (256 characters) code pages that first get 
converted to UTF-8 before being printed to the terminal.

For the legacy 8-bit texts, the input encoding will be assumed to be 
Code Page 437 otherwise called OEM-US. But you can change this using
the --input flag.

Common Code Page documents for English texts are:
  Code Page 437 (OEM-US)
  Code Page 850 (OEM Multilingual Latin 1)
  Code Page 858 (OEM Multilingual Latin 1 with the € symbol)
  ISO 8859-1 (Latin 1 commonly found on the web in the 2000s)
  Windows 1252 (Used in consumer Windows of the 1990s)

Otherwise the flags are optional and can be generally ignored 
for most use cases.`

func ViewCommand() *cobra.Command {
	s := "Print a text file to the terminal using standard output"
	l := viewLong
	expl := strings.Builder{}
	example.View.String(&expl)
	return &cobra.Command{
		Use:     fmt.Sprintf("view %s", example.Filenames),
		Aliases: []string{"v"},
		GroupID: IDfile,
		Short:   s,
		Long:    l,
		Example: expl.String(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return view.Run(cmd.OutOrStdout(), cmd, args...)
		},
	}
}

func ViewInit() *cobra.Command {
	vc := ViewCommand()
	f := flag.View()
	flag.Encode(&f.Input, vc)
	flag.Controls(&f.Controls, vc)
	flag.SwapChars(&f.Swap, vc)
	if err := flag.OG(&f.Original, vc); err != nil {
		log.Fatal(err)
	}
	flag.Width(&f.Width, vc)
	vc.Flags().SortFlags = false
	return vc
}

func init() {
	Cmd.AddCommand(ViewInit())
}
