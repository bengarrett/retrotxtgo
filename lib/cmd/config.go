// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"os"

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
	configExample = fmt.Sprintf(`  %s config setup  # to start the setup walkthrough
%s config set -c # to list all available settings`, meta.Bin, meta.Bin)
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   fmt.Sprintf("%s configuration and save settings", meta.Name),
	Example: exampleCmd(configExample),
	Run: func(cmd *cobra.Command, args []string) {
		if !printUsage(cmd, args...) {
			logs.InvalidCommand("config", args...)
		}
	},
}

var configCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create a new or reset the config file",
	Run: func(cmd *cobra.Command, args []string) {
		if configCreate() {
			os.Exit(1)
		}
	},
}

func configCreate() bool {
	if err := config.Create(viper.ConfigFileUsed(), configFlag.ow); err != nil {
		logs.Problemf(logs.ErrCfgCreate, err)
		return true
	}
	fmt.Println("New config file:", viper.ConfigFileUsed())
	return false
}

var configDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d", "del", "rm"},
	Short:   "Remove the config file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Delete(); err != nil {
			logs.Fatal(err)
		}
	},
}

// Note: Previously I inserted the results of config.Editor() into
// the Short and Long fields. This will cause a logic error because
// viper.GetString("editor") is not yet set and the EDITOR env value
// will instead always be used.
var configEditCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"e"},
	Short:   "Edit the config file",
	Long: `Edit the config file

To switch editors either:
  Set one by creating or changing the ` + str.Example("$EDITOR") +
		` environment variable in your shell configuration.
  Set an editor in the configuration file, ` +
		str.Example(fmt.Sprintf("%s config set --name=editor", meta.Bin)),
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Edit(); err != nil {
			logs.Fatal(err)
		}
	},
}

var configInfoExample = fmt.Sprintf(`  %s config info   # to list the default setting values
%s config set -c # to list the settings and help hints`, meta.Bin, meta.Bin)

var configInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i"},
	Example: exampleCmd(configInfoExample),
	Short:   "View all the settings configured in the config file",
	Run: func(cmd *cobra.Command, args []string) {
		if configInfo() {
			os.Exit(0)
		}
	},
}

func configInfo() (exit bool) {
	if configFlag.configs {
		if err := config.List(); err != nil {
			logs.ProblemCmdFatal("config info", "list", err)
		}
		return true
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

var configSetExample = fmt.Sprintf(`  %s config set html.meta.description # to change the meta description setting
%s config set style.info style.html # to set the color styles`, meta.Bin, meta.Bin)

var configSetCmd = &cobra.Command{
	Use:     "set [setting names]",
	Aliases: []string{"s"},
	Short:   fmt.Sprintf("Change individual %s settings", meta.Name),
	Example: exampleCmd(configSetExample),
	Run: func(cmd *cobra.Command, args []string) {
		if configSet() {
			os.Exit(0)
		}
		if !printUsage(cmd, args...) {
			for _, arg := range args {
				config.Set(arg)
			}
		}
	},
}

func configSet() bool {
	if configFlag.configs {
		if err := config.List(); err != nil {
			logs.ProblemCmdFatal("config", "list", err)
		}
		return true
	}
	return false
}

var configSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: fmt.Sprintf("Setup all the available %s settings", meta.Name),
	Run: func(cmd *cobra.Command, args []string) {
		config.Setup()
	},
}

func init() {
	if str.Term(str.GetEnv("COLORTERM"), str.GetEnv("TERM")) == "none" {
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
