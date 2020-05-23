package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/info"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
)

var infoArgs struct {
	filename string
	format   string
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information on a text file",
	Run: func(cmd *cobra.Command, args []string) {
		if infoArgs.filename != "" {
			if err := info.Info(infoArgs.filename, infoArgs.format); err.Msg != nil {
				if fmt.Sprint(err.Msg) == "format:invalid" {
					err := cmd.Usage()
					fmt.Println()
					logs.Check("info usage:", err)
					logs.CheckFlag("format", infoArgs.format, config.Format.Info)
				}
				logs.ChkErr(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&infoArgs.filename, "name", "n", "",
		logs.Cp("text file to analyse")+" (required)\n")
	infoCmd.Flags().StringVarP(&infoArgs.format, "format", "f", "color",
		"output format \noptions: "+logs.Ci(config.Format.String("info")))
	err := infoCmd.MarkFlagRequired("name")
	logs.Check("name flag", err)
	infoCmd.Flags().SortFlags = false
}
