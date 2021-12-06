// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"bytes"
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/createcmd"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
)

// createCmd represents the create command.
var createCmd = &cobra.Command{
	Use:     fmt.Sprintf("create %s", filenames),
	Aliases: []string{"c", "html"},
	Short:   "Create a HTML document from text files",
	Long:    "Create a HTML document from text documents and text art files.",
	Example: example.Print(example.Create),
	Run: func(cmd *cobra.Command, args []string) {
		f := convert.Flag{
			Controls:  flag.CreateDefaults.Controls,
			SwapChars: flag.CreateDefaults.Swap,
		}
		// handle defaults, use these control codes
		if c := cmd.Flags().Lookup("controls"); !c.Changed {
			f.Controls = []string{"eof", "tab"}
		}
		// handle defaults, swap out these characters with UTF-8 alternatives
		if s := cmd.Flags().Lookup("swap-chars"); !s.Changed {
			f.SwapChars = []string{null, verticalBar}
		}
		// handle the defaults for most other flags
		createcmd.Strings(cmd)
		// handle standard input (stdio)
		if filesystem.IsPipe() {
			createcmd.ParsePipe(cmd, f)
			return
		}
		// handle the hidden --body flag value,
		// used for debugging, it ignores most other flags and
		// overrides the <pre></pre> content before exiting
		if body := cmd.Flags().Lookup("body"); body.Changed {
			createcmd.ParseBody(cmd)
			return
		}
		if err := flag.PrintUsage(cmd, args...); err != nil {
			logs.Fatal(err)
		}
		createcmd.ParseFiles(cmd, f, args...)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	// root config must be initialized before getting saved default values
	initConfig()
	// output flags
	flagEncode(&flag.CreateDefaults.Encode, createCmd)
	flagControls(&flag.CreateDefaults.Controls, createCmd)
	flagRunes(&flag.CreateDefaults.Swap, createCmd)
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
