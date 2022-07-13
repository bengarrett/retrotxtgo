package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func ConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   fmt.Sprintf("%s configuration and defaults", meta.Name),
		Long:    fmt.Sprintf("%s settings, setup and default configurations.", meta.Name),
		Example: example.Config.Print(),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := flag.PrintUsage(cmd, args...); err != nil {
				return err
			}
			if len(args) > 0 {
				logs.FatalCmd("config", args...) // TODO: move funcs into \cmd
			}
			return nil
		},
	}
}

func ConfigInit() *cobra.Command {
	cc := ConfigCommand()
	cc.AddCommand(Create.Command())
	cc.AddCommand(Delete.Command())
	cc.AddCommand(Edit.Command())
	cc.AddCommand(Info.Command())
	cc.AddCommand(Set.Command())
	cc.AddCommand(Setup.Command())
	// create
	Create.Command().Flags().BoolVarP(&flag.Config.Ow, "overwrite", "y", false,
		"overwrite and reset the existing config file")
	// info
	Info.Command().Flags().BoolVarP(&flag.Config.Configs, "configs", "c", false,
		"list all the available configuration setting names")
	Info.Command().Flags().StringVarP(&flag.Config.Style, "style", "s", "",
		"choose a syntax highligher")
	Info.Command().Flags().BoolVar(&flag.Config.Styles, "styles", false,
		"list and preview the available syntax highlighers")
	// set
	Set.Command().Flags().BoolVarP(&flag.Config.Configs, "list", "l", false,
		"list all the available setting names")
	// hidden test flag
	cc.PersistentFlags().BoolVar(&flag.Config.Test, "test", false,
		"hidden flag to use an alternative config for config testing")
	if err := cc.PersistentFlags().MarkHidden("test"); err != nil {
		logs.FatalMark("test", ErrHide, err)
	}
	return cc
}

// init is always called by the Cobra library to be used for global flags and commands.
//nolint:gochecknoinits
func init() {
	const highColor, basicColor = "COLORTERM", "TERM"
	if str.Term(str.GetEnv(highColor), str.GetEnv(basicColor)) == "none" {
		// disable all color output
		color.Enable = false
	}
	Cmd.AddCommand(ConfigInit())
}
