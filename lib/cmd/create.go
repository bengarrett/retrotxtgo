package cmd

import (
	"bytes"
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/createcmd"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/rootcmd"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
)

func createCommand() *cobra.Command {
	return &cobra.Command{
		Use:     fmt.Sprintf("create %s", example.Filenames),
		Aliases: []string{"c", "html"},
		Short:   "Create a HTML document from text files",
		Long:    "Create a HTML document from text documents and text art files.",
		Example: example.Create.Print(),
		Run: func(cmd *cobra.Command, args []string) {
			if err := createcmd.Run(cmd, args); err != nil {
				logs.Fatal(err)
			}
		},
	}
}

func init() {
	createCmd := createCommand()
	rootCmd.AddCommand(createCmd)
	// root config must be initialized before getting saved default values
	rootcmd.Init()
	// output flags
	deflts := flag.CreateDefaults()
	flag.Encode(&deflts.Encode, createCmd)
	flag.Controls(&deflts.Controls, createCmd)
	flag.Runes(&deflts.Swap, createCmd)
	dir := createcmd.SaveDir()
	createCmd.Flags().BoolVarP(&flag.HTML.Save.AsFiles, "save", "s", false,
		"save HTML and static files to a the save directory\nor ignore to print (save directory: "+dir+")")
	createCmd.Flags().BoolVarP(&flag.HTML.Save.Compress, "compress", "z", false,
		"store and compress all files into an archive when saving")
	createCmd.Flags().BoolVarP(&flag.HTML.Save.OW, "overwrite", "o", false,
		"overwrite any existing files when saving")
	// meta and html related flags.
	flags := flag.Init()
	keys := flag.Sort(flags)
	for _, i := range keys {
		c := flags[i]
		var buf bytes.Buffer
		buf = c.Body(buf)
		buf = c.Init(createCmd, buf)
	}
	createCmd.Flags().BoolVarP(&flag.HTML.SauceData.Use, "sauce", "", true,
		"use any found SAUCE metadata as HTML meta tags")
	if err := createCmd.Flags().MarkHidden("body"); err != nil {
		logs.FatalMark("body", ErrHide, err)
	}
	if err := createCmd.Flags().MarkHidden("cache"); err != nil {
		logs.FatalMark("cache", ErrHide, err)
	}
	createCmd.Flags().SortFlags = false
}
