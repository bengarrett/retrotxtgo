package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/config"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

type configFlags struct {
	configs bool
	ow      bool
	shell   string
	style   string
	styles  bool
}

var configFlag configFlags

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   "Configure and save settings for RetroTxt",
	Run: func(cmd *cobra.Command, args []string) {
		checkUse(cmd, args)
		logs.ArgFatal(args)
	},
}

var configCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create a new or reset the config file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Create(viper.ConfigFileUsed(), configFlag.ow); err != nil {
			logs.Fatal("config", "create", err)
		}
	},
}

var configDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d", "del", "rm"},
	Short:   "Remove the config file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Delete(); err.Err != nil {
			err.Fatal()
		}
	},
}

// note: Previously I inserted the results of config.Editor() into
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
		str.Example("retrotxt config set --name=editor"),
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Edit(); err.Err != nil {
			err.Fatal()
		}
	},
}

var configInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i"},
	Short:   "View all the settings configured in the config file",
	Run: func(cmd *cobra.Command, args []string) {
		if configFlag.configs {
			if err := config.List(); err != nil {
				logs.Fatal("config info", "list", err)
			}
			os.Exit(0)
		}
		if configFlag.styles {
			str.JSONStyles("retrotxt info --style")
			os.Exit(0)
		}
		if err := config.Info(configFlag.style); err.Err != nil {
			err.Fatal()
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:     "set [setting names]",
	Aliases: []string{"s"},
	Short:   "Change individual Retrotxt settings",
	Example: str.Example("  retrotxt config set html.meta.description") +
		" # to change the meta description setting\n" +
		str.Example("  retrotxt config set style.info style.html") +
		"   # to set the color styles",
	Run: func(cmd *cobra.Command, args []string) {
		if configFlag.configs {
			if err := config.List(); err != nil {
				logs.Fatal("config", "list", err)
			}
			os.Exit(0)
		}
		checkUse(cmd, args)
		for _, arg := range args {
			config.Set(arg)
		}
	},
}

var configSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup all the available Retrotxt settings",
	Run: func(cmd *cobra.Command, args []string) {
		config.Setup()
	},
}

var configShellCmd = &cobra.Command{
	Use:     "shell",
	Aliases: []string{"sh"},
	Short:   "Apply autocompletion a terminal shell",
	Example: str.Example("  retrotxt config shell --interpreter string [flags]") +
		str.Example("\n  retrotxt config shell -i=bash >> ~/.bash_profile") +
		str.Example("\n  retrotxt config shell -i=zsh >> ~/.zshrc"),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			buf bytes.Buffer
			err error
		)
		lexer, style := "", viper.GetString("style.html")
		switch configFlag.shell {
		case "bash", "bsh", "b":
			lexer = "bash"
			if err = cmd.GenBashCompletion(&buf); err != nil {
				logs.Fatal("shell", "bash", err)
			}
		case "powershell", "posh", "ps", "p":
			lexer = "powershell"
			if err = cmd.GenPowerShellCompletion(&buf); err != nil {
				logs.Fatal("shell", "powershell", err)
			}
		case "zsh", "z":
			lexer = "bash"
			if err = cmd.GenZshCompletion(&buf); err != nil {
				logs.Fatal("shell", "zsh", err)
			}
		default:
			logs.Fatal("the interpreter is not supported:",
				configFlag.shell,
				fmt.Errorf("options: %s", config.Format.String("shell")))
		}
		if err := str.Highlight(buf.String(), lexer, style); err != nil {
			logs.Fatal("config", "shell", err)
		}
	},
}

func init() {
	if str.Term() == "none" {
		color.Enable = false
	}
	var err error
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configCreateCmd)
	configCmd.AddCommand(configDeleteCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configInfoCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configSetupCmd)
	configCmd.AddCommand(configShellCmd)
	// create
	configCreateCmd.Flags().BoolVarP(&configFlag.ow, "overwrite", "y", false,
		"overwrite and reset the existing config file")
	// info
	configInfoCmd.Flags().BoolVarP(&configFlag.configs, "configs", "c", false,
		"list all the available configuration setting names")
	configInfoCmd.Flags().StringVarP(&configFlag.style, "style", "s", "dracula",
		"choose a syntax highligher")
	configInfoCmd.Flags().BoolVar(&configFlag.styles, "styles", false,
		"list and preview the available syntax highlighers")
	// set
	configSetCmd.Flags().BoolVarP(&configFlag.configs, "list", "l", false,
		"list all the available setting names")
	// shell
	configShellCmd.Flags().StringVarP(&configFlag.shell, "interpreter", "i", "",
		str.Required("user shell to receive retrotxt auto-completions")+
			str.Options("", config.Format.Shell, true))
	if err = configShellCmd.MarkFlagRequired("interpreter"); err != nil {
		logs.Fatal("interpreter flag", "", err)
	}
	configShellCmd.SilenceErrors = true
}
