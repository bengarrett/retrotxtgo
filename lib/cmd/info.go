package cmd

import (
	"os"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/info"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/cobra"
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
			logs.Check("piped.info", err)
			info.Stdin(b, infoFlag.format)
			logs.Check("piped.info", err)
			os.Exit(0)
		}
		checkUse(cmd, args)
		for _, arg := range args {
			if e := info.Info(arg, infoFlag.format); e.Msg != nil {
				cmd.Usage()
				logs.ChkErr(e)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&infoFlag.format, "format", "f", "color",
		str.Options("output format", config.Format.Info, true))
}
