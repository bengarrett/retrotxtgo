package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"retrotxt.com/retrotxt/lib/config"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/info"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

var infoFlag struct {
	format string
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:     "info [filenames]",
	Aliases: []string{"i"},
	Short:   "Information on a text file",
	Example: "  retrotxt info text.asc logo.jpg\n  retrotxt info file.txt --format=json",
	Run: func(cmd *cobra.Command, args []string) {
		// piped input from other programs
		if filesystem.IsPipe() {
			b, err := filesystem.ReadPipe()
			if err != nil {
				logs.Fatal("info", "read stdin", err)
			}
			if err = info.Stdin(b, infoFlag.format); err != nil {
				logs.Fatal("info", "parse stdin", err)
			}
			os.Exit(0)
		}
		checkUse(cmd, args)
		for _, arg := range args {
			if err := info.Info(arg, infoFlag.format); err.Err != nil {
				cmd.Usage()
				err.Fatal()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&infoFlag.format, "format", "f", "color",
		str.Options("output format", config.Format.Info, true))
}
