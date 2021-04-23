// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/config"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

const bash, zsh = "bash", "zsh"

var ErrShellCompletion = errors.New("could not generate completion for")

type configFlags struct {
	configs bool
	ow      bool
	styles  bool
	shell   string
	style   string
}

var configFlag configFlags

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   "Configure and save settings for RetroTxt",
	Run: func(cmd *cobra.Command, args []string) {
		checkUse(cmd, args...)
		logs.ArgFatal(args...)
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
		logs.Println("config", "create", err)
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
		if err := config.Delete(); err.Err != nil {
			err.Fatal()
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
		if configInfo() {
			os.Exit(0)
		}
	},
}

func configInfo() (exit bool) {
	if configFlag.configs {
		if err := config.List(); err != nil {
			logs.CmdProblemFatal("config info", "list", err)
		}
		return true
	}
	if configFlag.styles {
		str.JSONStyles("retrotxt info --style")
		return true
	}
	style := viper.GetString("style.info")
	if configFlag.style != "" {
		style = configFlag.style
	}
	if style == "" {
		style = "dracula"
	}
	if err := config.Info(style); err.Err != nil {
		err.Fatal()
	}
	return false
}

const configSetExample = `  retrotxt config set html.meta.description # to change the meta description setting
retrotxt config set style.info style.html # to set the color styles`

var configSetCmd = &cobra.Command{
	Use:     "set [setting names]",
	Aliases: []string{"s"},
	Short:   "Change individual Retrotxt settings",
	Example: exampleCmd(configSetExample),
	Run: func(cmd *cobra.Command, args []string) {
		if configSet() {
			os.Exit(0)
		}
		checkUse(cmd, args...)
		for _, arg := range args {
			config.Set(arg)
		}
	},
}

func configSet() bool {
	if configFlag.configs {
		if err := config.List(); err != nil {
			logs.CmdProblemFatal("config", "list", err)
		}
		return true
	}
	return false
}

var configSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup all the available Retrotxt settings",
	Run: func(cmd *cobra.Command, args []string) {
		config.Setup()
	},
}

const configShellExample = `  retrotxt config shell --interpreter string [flags]
  retrotxt config shell -i=bash >> ~/.bash_profile
  retrotxt config shell -i=zsh >> ~/.zshrc`

var configShellCmd = &cobra.Command{
	Use:     "shell",
	Aliases: []string{"sh"},
	Short:   "Apply autocompletion a terminal shell",
	Example: exampleCmd(configShellExample),
	Run: func(cmd *cobra.Command, args []string) {
		const ps = "powershell"
		var (
			buf bytes.Buffer
			err error
		)
		lexer, style := "", viper.GetString("style.html")
		switch configFlag.shell {
		case bash, "bsh", "b":
			lexer = bash
			if err = cmd.GenBashCompletion(&buf); err != nil {
				logs.MarkProblemFatal(bash, ErrShellCompletion, err)
			}
		case ps, "posh", "ps", "p":
			lexer = ps
			if err = cmd.GenPowerShellCompletion(&buf); err != nil {
				logs.MarkProblemFatal(ps, ErrShellCompletion, err)
			}
		case zsh, "z":
			lexer = bash
			if err = cmd.GenZshCompletion(&buf); err != nil {
				logs.MarkProblemFatal(zsh, ErrShellCompletion, err)
			}
		default:
			s := config.Format().Shell
			logs.Fatal(fmt.Sprintf("options: %s", s[:]),
				configFlag.shell,
				ErrIntpr)
		}
		if err := str.Highlight(buf.String(), lexer, style, true); err != nil {
			logs.MarkProblemFatal("shell", logs.ErrHighlight, err)
		}
	},
}

func init() {
	if str.Term(str.GetEnv("COLORTERM"), str.GetEnv("TERM")) == "none" {
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
	configInfoCmd.Flags().StringVarP(&configFlag.style, "style", "s", "",
		"choose a syntax highligher")
	configInfoCmd.Flags().BoolVar(&configFlag.styles, "styles", false,
		"list and preview the available syntax highlighers")
	// set
	configSetCmd.Flags().BoolVarP(&configFlag.configs, "list", "l", false,
		"list all the available setting names")
	// shell
	s := config.Format().Shell
	configShellCmd.Flags().StringVarP(&configFlag.shell, "interpreter", "i", "",
		str.Required("user shell to receive retrotxt auto-completions")+
			str.Options("", true, s[:]...))
	if err = configShellCmd.MarkFlagRequired("interpreter"); err != nil {
		logs.MarkProblemFatal("interpreter", logs.ErrMarkRequire, err)
	}
	configShellCmd.SilenceErrors = true
}
