// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/configcmd"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   fmt.Sprintf("%s configuration and defaults", meta.Name),
	Long:    fmt.Sprintf("%s settings, setup and default configurations.", meta.Name),
	Example: exampleCmd(example.Config),
	Run: func(cmd *cobra.Command, args []string) {
		if err := flag.PrintUsage(cmd, args...); err != nil {
			logs.Fatal(err)
		}
		logs.FatalCmd("config", args...)
	},
}

// init is always called by the Cobra library to be used for global flags and commands.
func init() {
	const highColor, basicColor = "COLORTERM", "TERM"
	if str.Term(str.GetEnv(highColor), str.GetEnv(basicColor)) == "none" {
		// disable all color output
		color.Enable = false
	}
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configcmd.Create)
	configCmd.AddCommand(configcmd.Delete)
	configCmd.AddCommand(configcmd.Edit)
	configCmd.AddCommand(configcmd.Info)
	configCmd.AddCommand(configcmd.Set)
	configCmd.AddCommand(configcmd.Setup)
	// create
	configcmd.Create.Flags().BoolVarP(&flag.Config.Ow, "overwrite", "y", false,
		"overwrite and reset the existing config file")
	// info
	configcmd.Info.Flags().BoolVarP(&flag.Config.Configs, "configs", "c", false,
		"list all the available configuration setting names")
	configcmd.Info.Flags().StringVarP(&flag.Config.Style, "style", "s", "",
		"choose a syntax highligher")
	configcmd.Info.Flags().BoolVar(&flag.Config.Styles, "styles", false,
		"list and preview the available syntax highlighers")
	// set
	configcmd.Set.Flags().BoolVarP(&flag.Config.Configs, "list", "l", false,
		"list all the available setting names")
}
