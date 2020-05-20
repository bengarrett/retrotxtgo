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
	ow    bool
	set   string
	shell string
	style string
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
		config.Delete()
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the config file",
	// TODO: add information on configuring an editor
	Run: func(cmd *cobra.Command, args []string) {
		if e := config.Edit(); e.Err != nil {
			e.Exit(1)
		}
	},
}

var configInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "View settings configured by the config",
	Run: func(cmd *cobra.Command, args []string) {
		config.Info()
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Change a configuration",
	Example: `  --name create.meta.description # to change the meta description setting
  --name version.format          # to set the version command output format`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Set(configArgs.set)
	},
}

var configShellCmd = &cobra.Command{
	Use:     "shell",
	Short:   "Apply autocompletion a terminal shell",
	Example: "  retrotxt config shell --interpreter string [flags]",
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
	configInfoCmd.Flags().StringVarP(&configArgs.style, "syntax-style", "c", "monokai",
		"config syntax highligher, use "+logs.Ci("none")+" to disable")
	// set
	configSetCmd.Flags().StringVarP(&configArgs.set, "name", "n", "",
		`the configuration path to edit in dot syntax (see examples)
to see a list of names run: retrotxt config info`)
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
