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

// infoCmd represents the info command.
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
			if err = info.Stdin(infoFlag.format, b...); err != nil {
				logs.Fatal("info", "parse stdin", err)
			}
			os.Exit(0)
		}
		checkUse(cmd, args...)
		var n info.Names
		n.Length = len(args)
		for i, arg := range args {
			n.Index = i + 1
			if err := n.Info(arg, infoFlag.format); err.Err != nil {
				if err.Err == info.ErrNoFile {
					err.Fatal()
				}
				if err := cmd.Usage(); err != nil {
					logs.Println("command", "usage", err)
				}
				err.Fatal()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	i := config.Format().Info
	infoCmd.Flags().StringVarP(&infoFlag.format, "format", "f", "color",
		str.Options("output format", true, i[:]...))
}
