package cmd

/*
fixes:
1:  config create --config=honk.yml
*/

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type configFlags struct {
	configs bool
	list    bool
	ow      bool
	set     string
	shell   string
	style   string
}

var configArgs configFlags

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure and save settings for RetroTxt",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Usage()
		logs.Check("config usage:", err)
		if len(args) != 0 || cmd.Flags().NFlag() != 0 {
			logs.CheckArg("config", args)
		}
	},
}

var configCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new or reset the config file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Create(viper.ConfigFileUsed(), configArgs.ow); err != nil {
			logs.Check("config create", err)
		}
	},
}

var configDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove the config file",
	Run: func(cmd *cobra.Command, args []string) {
		if e := config.Delete(); e.Err != nil {
			e.Exit(1)
		}
	},
}

var configEditor = func() string {
	e := config.Editor()
	if e == "" {
		return ""
	}
	return fmt.Sprintf(" using %s", e)
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: fmt.Sprintf("Edit the config file%s", configEditor()),
	Long: fmt.Sprintf("Edit the config file%s", configEditor()) +
		"\n\nTo switch editors either:" +
		"\n  Set one by creating or changing the " +
		logs.Example("$EDITOR") +
		" environment variable in your shell configuration." +
		"\n  Set an editor in the configuration file, " +
		logs.Example("retrotxt config set --name=editor"),
	Run: func(cmd *cobra.Command, args []string) {
		if e := config.Edit(); e.Err != nil {
			e.Exit(1)
		}
	},
}

var configInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "View all the settings configured in the config file",
	Example: logs.Example("  retrotxt config info --syntax-style=\"\"") +
		" # disable the syntax highligher",
	Run: func(cmd *cobra.Command, args []string) {
		if configArgs.configs {
			config.List()
		} else if configArgs.list {
			logs.YamlStyles("retrotxt info --style")
		} else if e := config.Info(configArgs.style); e.Err != nil {
			e.Exit(1)
		}
	},
}

var configSetExample = func() string {
	return logs.Example("  retrotxt config set --name create.meta.description") +
		" # to change the meta description setting\n" +
		logs.Example("  retrotxt config set --name version.format") +
		"          # to set the version command output format"
}

var configSetCmd = &cobra.Command{
	Use:     "set",
	Short:   "Change a Retrotxt setting",
	Example: configSetExample(),
	Run: func(cmd *cobra.Command, args []string) {
		config.Set(configArgs.set)
	},
}

var configShellCmd = &cobra.Command{
	Use:     "shell",
	Short:   "Apply autocompletion a terminal shell",
	Example: logs.Example("  retrotxt config shell --interpreter string [flags]"),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			buf   bytes.Buffer
			err   error
			lexer string
		)
		switch configArgs.shell {
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
				Arg: configArgs.shell,
				Msg: fmt.Errorf("options: %s", config.Format.String("shell"))})
		}
		if err := quick.Highlight(os.Stdout, buf.String(), lexer, "terminal256", "monokai"); err != nil {
			fmt.Println(buf.String())
		}
	},
}

func init() {
	var err error
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configCreateCmd)
	configCmd.AddCommand(configDeleteCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configInfoCmd)
	configCmd.AddCommand(configShellCmd)
	configCmd.AddCommand(configSetCmd)
	// create
	configCreateCmd.Flags().BoolVarP(&configArgs.ow, "overwrite", "y", false,
		"overwrite and reset the existing config file")
	// info
	configInfoCmd.Flags().BoolVarP(&configArgs.configs, "configs", "c", false,
		"list all the available configuration settings")
	configInfoCmd.Flags().StringVarP(&configArgs.style, "style", "s", "monokai",
		"choose a syntax highligher")
	configInfoCmd.Flags().BoolVarP(&configArgs.list, "list", "l", false,
		"list and preview the available syntax highlighers")
	// set
	configSetCmd.Flags().StringVarP(&configArgs.set, "name", "n", "",
		fmt.Sprintf("the setting name in dot syntax%s", logs.Required())+
			fmt.Sprintf("\nrun %s", logs.Example("retrotxt config info"))+
			" to see a list of names")
	err = configSetCmd.MarkFlagRequired("name")
	logs.Check("name flag", err)
	// shell
	configShellCmd.Flags().StringVarP(&configArgs.shell, "interpreter", "i", "",
		"user shell to receive retrotxtgo auto-completions\nchoices: "+
			logs.Ci(strings.Join(config.Format.Shell, ", ")))
	err = configShellCmd.MarkFlagRequired("interpreter")
	logs.Check("interpreter flag", err)
	configShellCmd.SilenceErrors = true
}
