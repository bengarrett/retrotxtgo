package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Configs int

const (
	Create Configs = iota
	Delete
	Edit
	Info
	Set
	Setup
)

func (c Configs) Command() *cobra.Command {
	switch c {
	case Create:
		return ConfigCreate()
	case Delete:
		return ConfigDel()
	case Edit:
		return ConfigEdit()
	case Info:
		return ConfigInfo()
	case Set:
		return ConfigSet()
	case Setup:
		return ConfigSetup()
	}
	return nil
}

type Configer struct {
	Configs bool
	OW      bool
	Styles  bool
	Test    bool
	Style   string
}

var Config Configer

func ConfigCreate() *cobra.Command {
	return &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create or reset the config file",
		Long:    fmt.Sprintf("Create or reset the %s configuration file.", meta.Name),
		RunE: func(cmd *cobra.Command, args []string) error {
			overwrite := Config.OW
			b, err := config.New(overwrite)
			if errors.Is(err, config.ErrExist) {
				b = config.DoesExist(config.CmdPath(), "create")
				fmt.Fprint(cmd.OutOrStdout(), b)
				return nil
			}
			if err != nil {
				return fmt.Errorf("%w: %s", logs.ErrConfigNew, err)
			}
			fmt.Fprint(cmd.OutOrStdout(), b)
			return nil
		},
	}
}

func ConfigDel() *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Aliases: []string{"d", "del", "rm"},
		Short:   "Remove the config file",
		Long:    fmt.Sprintf("Remove the %s configuration file.", meta.Name),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Delete(); err != nil {
				return err
			}
			return nil
		},
	}
}

// Note: Previously I inserted the results of config.Editor() into
// the Short and Long fields. This will cause a logic error because
// viper.GetString("editor") is not yet set and the EDITOR env value
// will instead always be used.
func ConfigEdit() *cobra.Command {
	long := fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n",
		fmt.Sprintf("Edit the %s configuration file.", meta.Name),
		"To change the editor program, either:",
		fmt.Sprintf("  1. Configure one by creating a %s shell environment variable.",
			str.Example("$EDITOR")),
		"  2. Set an editor in the configuration file:",
		str.Example(fmt.Sprintf("     %s config set --name=editor", meta.Bin)),
	)
	return &cobra.Command{
		Use:     "edit",
		Aliases: []string{"e"},
		Short:   "Edit the config file\n",
		Long:    long,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Edit(); err != nil {
				return err
			}
			return nil
		},
	}
}

func ConfigInfo() *cobra.Command {
	return &cobra.Command{
		Use:     "info",
		Aliases: []string{"i"},
		Example: fmt.Sprint(example.ConfigInfo),
		Short:   "List all the settings in use",
		Long:    fmt.Sprintf("List all the %s settings in use.", meta.Name),
		RunE: func(cmd *cobra.Command, args []string) error {
			if ConfigInfoer() {
				return nil
			}
			return nil
		},
	}
}

func ConfigSet() *cobra.Command {
	return &cobra.Command{
		Use:     "set [setting names]",
		Aliases: []string{"s"},
		Short:   "Edit a setting",
		Long:    fmt.Sprintf("Edit a %s setting.", meta.Name),
		Example: fmt.Sprint(example.Set),
		RunE: func(cmd *cobra.Command, args []string) error {
			if ListAll() {
				return nil
			}
			if err := Usage(cmd, args...); err != nil {
				return err
			}
			for _, arg := range args {
				if err := config.Set(arg); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func ConfigSetup() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Walk through all the settings",
		Long:  fmt.Sprintf("Walk through all of the %s settings.", meta.Name),
		RunE: func(cmd *cobra.Command, args []string) error {
			const startAt = 0
			config.Setup(startAt)
			return nil
		},
	}
}

// ListAll is the "config set --list" command run.
func ListAll() (exit bool) {
	if Config.Configs {
		if err := config.List(); err != nil {
			logs.FatalFlag("config", "list", err)
		}
		return true
	}
	return false
}

///=======================

// Print the usage help and exit; but only when no arguments are given.
func Usage(cmd *cobra.Command, args ...string) error {
	if len(args) == 0 {
		if err := cmd.Help(); err != nil {
			return err
		}
		os.Exit(0)
	}
	return nil
}

// Init reads in the config file and ENV variables if set.
// This might be triggered twice due to the Cobra initializer registers.
func Load() {
	// read in environment variables
	viper.SetEnvPrefix("env")
	viper.AutomaticEnv()
	// tester configuration file
	if flag.Command.Tester {
		fmt.Println("The single use, in-memory tester file is in use.")
		fs := afero.NewMemMapFs()
		afs := &afero.Afero{Fs: fs}
		f, err := afs.TempFile("", "ioutil-test")
		if err != nil {
			logs.Fatal(err)
		}
		if err := config.Create(f.Name(), true); err != nil {
			logs.Fatal(err)
		}
		if err := config.SetConfig(f.Name()); err != nil {
			logs.FatalMark(viper.ConfigFileUsed(), logs.ErrConfigOpen, err)
		}
		return
	}
	// configuration file
	if err := config.SetConfig(flag.Command.Config); err != nil {
		logs.FatalMark(viper.ConfigFileUsed(), logs.ErrConfigOpen, err)
	}
}

func ConfigInfoer() (exit bool) {
	if err := config.SetConfig(flag.Command.Config); err != nil {
		logs.FatalMark(viper.ConfigFileUsed(), logs.ErrConfigOpen, err)
	}
	if Config.Configs {
		if err := config.List(); err != nil {
			logs.FatalFlag("config info", "list", err)
		}
	}
	if Config.Styles {
		fmt.Print(str.JSONStyles(fmt.Sprintf("%s info --style", meta.Bin)))
		return true
	}
	style := viper.GetString("style.info")
	if Config.Style != "" {
		style = Config.Style
	}
	if style == "" {
		style = "dracula"
	}
	if err := config.Info(style); err != nil {
		logs.Fatal(err)
	}
	return false
}
