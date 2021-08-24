// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type configFlags struct {
	configs bool
	ow      bool
	styles  bool
	style   string
}

var (
	configFlag    configFlags
	configExample = fmt.Sprintf("  %s %s %s\n%s %s %s",
		meta.Bin, "config setup", "# Walk through all the settings",
		meta.Bin, "config set --list", "# List all the settings in use")
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   fmt.Sprintf("%s configuration and defaults", meta.Name),
	Long:    fmt.Sprintf("%s settings, setup and default configurations.", meta.Name),
	Example: exampleCmd(configExample),
	Run: func(cmd *cobra.Command, args []string) {
		if !printUsage(cmd, args...) {
			logs.FatalCmd("config", args...)
		}
	},
}

var configCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create or reset the config file",
	Long:    fmt.Sprintf("Create or reset the %s configuration file.", meta.Name),
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.New(configFlag.ow); err != nil {
			logs.FatalWrap(logs.ErrCfgCreate, err)
		}
	},
}

var configDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d", "del", "rm"},
	Short:   "Remove the config file",
	Long:    fmt.Sprintf("Remove the %s configuration file.", meta.Name),
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Delete(); err != nil {
			logs.Fatal(err)
		}
	},
}

var configEditLong = fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n",
	fmt.Sprintf("Edit the %s configuration file.", meta.Name),
	"To change the editor program, either:",
	fmt.Sprintf("  1. Configure one by creating a %s shell environment variable.",
		str.Example("$EDITOR")),
	"  2. Set an editor in the configuration file:",
	str.Example(fmt.Sprintf("     %s config set --name=editor", meta.Bin)),
)

// Note: Previously I inserted the results of config.Editor() into
// the Short and Long fields. This will cause a logic error because
// viper.GetString("editor") is not yet set and the EDITOR env value
// will instead always be used.
var configEditCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"e"},
	Short:   "Edit the config file\n",
	Long:    configEditLong,
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Edit(); err != nil {
			logs.Fatal(err)
		}
	},
}

var configInfoExample = fmt.Sprintf(`  %s config info   # List the default setting values
%s config set -c # List the settings and help hints`, meta.Bin, meta.Bin)

var configInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i"},
	Example: exampleCmd(configInfoExample),
	Short:   "List all the settings in use",
	Long:    fmt.Sprintf("List all the %s settings in use.", meta.Name),
	Run: func(cmd *cobra.Command, args []string) {
		if configInfo() {
			return
		}
	},
}

// configInfo is the "config info" run command.
func configInfo() (exit bool) {
	if configFlag.configs {
		if err := config.List(); err != nil {
			logs.FatalFlag("config info", "list", err)
		}
	}
	if configFlag.styles {
		str.JSONStyles(fmt.Sprintf("%s info --style", meta.Bin))
		return true
	}
	style := viper.GetString("style.info")
	if configFlag.style != "" {
		style = configFlag.style
	}
	if style == "" {
		style = "dracula"
	}
	if err := config.Info(style); err != nil {
		logs.Fatal(err)
	}
	return false
}

var configSetExample = fmt.Sprintf("  %s %s %s\n%s %s %s\n%s %s %s",
	meta.Bin, "config set --list", "# List the available settings",
	meta.Bin, "config set html.meta.description", "# Edit the meta description setting",
	meta.Bin, "config set style.info style.html", fmt.Sprintf("# Edit both the %s color styles", meta.Name),
)

var configSetCmd = &cobra.Command{
	Use:     "set [setting names]",
	Aliases: []string{"s"},
	Short:   "Edit a setting",
	Long:    fmt.Sprintf("Edit a %s setting.", meta.Name),
	Example: exampleCmd(configSetExample),
	Run: func(cmd *cobra.Command, args []string) {
		if configListAll() {
			return
		}
		if !printUsage(cmd, args...) {
			for _, arg := range args {
				config.Set(arg)
			}
		}
	},
}

// configListAll is the "config set --list" command run.
func configListAll() (ok bool) {
	if configFlag.configs {
		if err := config.List(); err != nil {
			logs.FatalFlag("config", "list", err)
		}
		return true
	}
	return false
}

var configSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Walk through all the settings",
	Long:  fmt.Sprintf("Walk through all of the %s settings.", meta.Name),
	Run: func(cmd *cobra.Command, args []string) {
		const startAt = 0
		config.Setup(startAt)
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
	configCmd.AddCommand(configCreateCmd)
	configCmd.AddCommand(configDeleteCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configInfoCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configSetupCmd)
	// create
	configCreateCmd.Flags().BoolVarP(&configFlag.ow, "overwrite", "y", false,
		"overwrite and reset the existing config file")
	// info
	configInfoCmd.Flags().BoolVarP(&configFlag.configs, "configs", "c", false,
		"list all the available configuration setting names")
	configInfoCmd.Flags().StringVarP(&configFlag.style, "style", "s", "",
		"choose a syntax highligher")
	configInfoCmd.Flags().BoolVar(&configFlag.styles, "styles", false,
		"list and preview the available syntax highlighers")
	// set
	configSetCmd.Flags().BoolVarP(&configFlag.configs, "list", "l", false,
		"list all the available setting names")
}
