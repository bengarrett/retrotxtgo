package cmd

import (
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/gookit/color"
)

// func ConfigCommand() *cobra.Command {
// 	return &cobra.Command{
// 		Use:     "config",
// 		Aliases: []string{"cfg"},
// 		Short:   fmt.Sprintf("%s configuration and defaults", meta.Name),
// 		Long:    fmt.Sprintf("%s settings, setup and default configurations.", meta.Name),
// 		Example: fmt.Sprint(example.Config),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			if err := flag.Help(cmd, args...); err != nil {
// 				return err
// 			}
// 			if len(args) > 0 {
// 				logs.FatalCmd("config", args...) // TODO: move funcs into \cmd
// 			}
// 			return nil
// 		},
// 	}
// }

// func ConfigInit() *cobra.Command {
// 	// create := Create.Command()
// 	// info := Info.Command()
// 	// sets := Set.Command()
// 	// cc := ConfigCommand()
// 	// cc.AddGroup(&cobra.Group{ID: "configfile", Title: "Config File:"})
// 	// cc.AddGroup(&cobra.Group{ID: "settings", Title: "Settings:"})
// 	// cc.AddCommand(create)
// 	// cc.AddCommand(Delete.Command())
// 	// cc.AddCommand(Edit.Command())
// 	// cc.AddCommand(info)
// 	// cc.AddCommand(sets)
// 	// cc.AddCommand(Setup.Command())
// 	// create
// 	// create.Flags().BoolVarP(&Config.OW, "overwrite", "y", false,
// 	// 	"overwrite and reset the existing config file")
// 	// // info
// 	// info.Flags().BoolVarP(&Config.Configs, "configs", "c", false,
// 	// 	"list all the available configuration setting names")
// 	// info.Flags().StringVarP(&Config.Style, "style", "s", "",
// 	// 	"choose a syntax highligher")
// 	// info.Flags().BoolVar(&Config.Styles, "styles", false,
// 	// 	"list and preview the available syntax highlighers")
// 	// // set
// 	// sets.Flags().BoolVarP(&Config.Configs, "list", "l", false,
// 	// 	"list all the available setting names")
// 	// hidden test flag
// 	// cc.PersistentFlags().BoolVar(&Config.Test, "test", false,
// 	// 	"hidden flag to use an alternative config for config testing")
// 	// if err := cc.PersistentFlags().MarkHidden("test"); err != nil {
// 	// 	logs.FatalMark("test", ErrHide, err)
// 	// }
// 	//return cc
// }

// init is always called by the Cobra library to be used for global flags and commands.
//
//nolint:gochecknoinits
func init() {
	const highColor, basicColor = "COLORTERM", "TERM"
	if str.Term(str.GetEnv(highColor), str.GetEnv(basicColor)) == "none" {
		// disable all color output
		color.Enable = false
	}
	//Cmd.AddCommand(ConfigInit())
}
