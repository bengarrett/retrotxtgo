package cmd

import (
	"fmt"

	cfg "github.com/bengarrett/retrotxtgo/cmd/internal/config"
	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func configCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   fmt.Sprintf("%s configuration and defaults", meta.Name),
		Long:    fmt.Sprintf("%s settings, setup and default configurations.", meta.Name),
		Example: example.Config.Print(),
		Run: func(cmd *cobra.Command, args []string) {
			if err := flag.PrintUsage(cmd, args...); err != nil {
				logs.Fatal(err)
			}
			logs.FatalCmd("config", args...)
		},
	}
}

// init is always called by the Cobra library to be used for global flags and commands.
//nolint:gochecknoinits
func init() {
	const highColor, basicColor = "COLORTERM", "TERM"
	if str.Term(str.GetEnv(highColor), str.GetEnv(basicColor)) == "none" {
		// disable all color output
		color.Enable = false
	}
	cc := configCommand()
	rootCmd.AddCommand(cc)
	cc.AddCommand(cfg.Create.Command())
	cc.AddCommand(cfg.Delete.Command())
	cc.AddCommand(cfg.Edit.Command())
	cc.AddCommand(cfg.Info.Command())
	cc.AddCommand(cfg.Set.Command())
	cc.AddCommand(cfg.Setup.Command())
	// create
	cfg.Create.Command().Flags().BoolVarP(&flag.Config.Ow, "overwrite", "y", false,
		"overwrite and reset the existing config file")
	// info
	cfg.Info.Command().Flags().BoolVarP(&flag.Config.Configs, "configs", "c", false,
		"list all the available configuration setting names")
	cfg.Info.Command().Flags().StringVarP(&flag.Config.Style, "style", "s", "",
		"choose a syntax highligher")
	cfg.Info.Command().Flags().BoolVar(&flag.Config.Styles, "styles", false,
		"list and preview the available syntax highlighers")
	// set
	cfg.Set.Command().Flags().BoolVarP(&flag.Config.Configs, "list", "l", false,
		"list all the available setting names")
}
