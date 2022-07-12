package cmd

import (
	"bytes"
	"fmt"

	"github.com/bengarrett/retrotxtgo/cmd/internal/create"
	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/root"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
)

func CreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:     fmt.Sprintf("create %s", example.Filenames),
		Aliases: []string{"c", "html"},
		Short:   "Create a HTML document from text files",
		Long:    "Create a HTML document from text documents and text art files.",
		Example: example.Create.Print(),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := create.Run(cmd, args); err != nil {
				return err
			}
			return nil
		},
	}
}

func CreateInit() *cobra.Command {
	cc := CreateCommand()
	// root config must be initialized before getting saved default values
	root.Init()
	// output flags
	deflts := flag.CreateDefaults()
	flag.Encode(&deflts.Encode, cc)
	flag.Controls(&deflts.Controls, cc)
	flag.Runes(&deflts.Swap, cc)
	dir := create.SaveDir()
	cc.Flags().BoolVarP(&flag.HTML.Save.AsFiles, "save", "s", false,
		"save HTML and static files to a the save directory\nor ignore to print (save directory: "+dir+")")
	cc.Flags().BoolVarP(&flag.HTML.Save.Compress, "compress", "z", false,
		"store and compress all files into an archive when saving")
	cc.Flags().BoolVarP(&flag.HTML.Save.OW, "overwrite", "o", false,
		"overwrite any existing files when saving")
	// meta and html related flags.
	flags := flag.Init()
	keys := flag.Sort(flags)
	for _, i := range keys {
		c := flags[i]
		var buf bytes.Buffer
		buf = c.Body(buf)
		c.Init(cc, buf)
	}
	cc.Flags().BoolVarP(&flag.HTML.SauceData.Use, "sauce", "", true,
		"use any found SAUCE metadata as HTML meta tags")
	if err := cc.Flags().MarkHidden("body"); err != nil {
		logs.FatalMark("body", ErrHide, err)
	}
	if err := cc.Flags().MarkHidden("cache"); err != nil {
		logs.FatalMark("cache", ErrHide, err)
	}
	cc.Flags().SortFlags = false
	return cc
}

//nolint:gochecknoinits
func init() {
	Cmd.AddCommand(CreateCommand())
}
