package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		logs.CheckCmd(args)
	},
}

var configCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create a new or reset the config file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Create(viper.ConfigFileUsed(), configFlag.ow); err != nil {
			logs.Check("config create", err)
		}
	},
}

var configDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d", "del", "rm"},
	Short:   "Remove the config file",
	Run: func(cmd *cobra.Command, args []string) {
		if e := config.Delete(); e.Err != nil {
			e.Exit(1)
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
	Short:   fmt.Sprintf("Edit the config file"),
	Long: fmt.Sprintf("Edit the config file") +
		"\n\nTo switch editors either:" +
		"\n  Set one by creating or changing the " +
		str.Example("$EDITOR") +
		" environment variable in your shell configuration." +
		"\n  Set an editor in the configuration file, " +
		str.Example("retrotxt config set --name=editor"),
	Run: func(cmd *cobra.Command, args []string) {
		if e := config.Edit(); e.Err != nil {
			e.Exit(1)
		}
	},
}

var configInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i"},
	Short:   "View all the settings configured in the config file",
	Example: str.Example("  retrotxt config info --syntax-style=\"\"") +
		" # disable the syntax highligher",
	Run: func(cmd *cobra.Command, args []string) {
		if configFlag.configs {
			config.List()
			os.Exit(0)
		}
		if configFlag.styles {
			str.JSONStyles("retrotxt info --style")
			os.Exit(0)
		}
		if e := config.Info(configFlag.style); e.Err != nil {
			e.Exit(1)
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:     "set [setting names]",
	Aliases: []string{"s"},
	Short:   "Change individual Retrotxt settings",
	Example: str.Example("  retrotxt config set create.meta.description") +
		" # to change the meta description setting\n" +
		str.Example("  retrotxt config set style.info style.html") +
		"   # to set the color styles",
	Run: func(cmd *cobra.Command, args []string) {
		if configFlag.configs {
			config.List()
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
			buf   bytes.Buffer
			err   error
			lexer string
			style string = viper.GetString("style.html")
		)
		switch configFlag.shell {
		case "bash", "bsh", "b":
			lexer = "bash"
			err = cmd.GenBashCompletion(&buf)
			logs.Check("shell bash", err)
		case "powershell", "posh", "ps", "p":
			lexer = "powershell"
			err = cmd.GenPowerShellCompletion(&buf)
			logs.Check("shell powershell", err)
		case "zsh", "z":
			lexer = "bash"
			err = cmd.GenZshCompletion(&buf)
			logs.Check("shell zsh", err)
		default:
			logs.ChkErr(logs.Err{Issue: "the interpreter is not supported:",
				Arg: configFlag.shell,
				Msg: fmt.Errorf("options: %s", config.Format.String("shell"))})
		}
		if err := str.Highlight(buf.String(), lexer, style); err != nil {
			logs.Check("config shell", err)
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
	err = configShellCmd.MarkFlagRequired("interpreter")
	logs.Check("interpreter flag", err)
	configShellCmd.SilenceErrors = true
}
