package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var (
	ErrNoArgs = errors.New("no arguments supplied")
	ErrWriter = errors.New("writer argument cannot be nil")
)

// type Configs int

// const (
// 	Create Configs = iota
// 	Delete
// 	Edit
// 	Info
// 	Set
// 	Setup
// )

// func (c Configs) Command() *cobra.Command {
// 	switch c {
// 	case Create:
// 		return ConfigCreate()
// 	case Delete:
// 		return ConfigDel()
// 	case Edit:
// 		return ConfigEdit()
// 	case Info:
// 		return ConfigInfo()
// 	case Set:
// 		return ConfigSet()
// 	case Setup:
// 		return ConfigSetup()
// 	}
// 	return nil
// }

// type Configer struct {
// 	Configs bool
// 	OW      bool
// 	Styles  bool
// 	Test    bool
// 	Style   string
// }

// var Config Configer

// func ConfigCreate() *cobra.Command {
// 	return &cobra.Command{
// 		Use:     "create",
// 		Aliases: []string{"c"},
// 		Short:   "Create or reset the config file",
// 		Long:    fmt.Sprintf("Create or reset the %s configuration file.", meta.Name),
// 		GroupID: "configfile",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			overwrite := Config.OW
// 			w := cmd.OutOrStdout()
// 			err := config.New(w, overwrite)
// 			if errors.Is(err, config.ErrExist) {
// 				return config.DoesExist(w, config.CmdPath(), "create")
// 			}
// 			if err != nil {
// 				return fmt.Errorf("%w: %s", logs.ErrConfigNew, err)
// 			}
// 			return nil
// 		},
// 	}
// }

// func ConfigDel() *cobra.Command {
// 	return &cobra.Command{
// 		Use:     "delete",
// 		Aliases: []string{"d", "del", "rm"},
// 		Short:   "Remove the config file",
// 		Long:    fmt.Sprintf("Remove the %s configuration file.", meta.Name),
// 		GroupID: "configfile",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			if err := config.Delete(cmd.OutOrStdout(), !flag.Command.Tester); err != nil {
// 				return err
// 			}
// 			return nil
// 		},
// 	}
// }

// Note: Previously I inserted the results of config.Editor() into
// the Short and Long fields. This will cause a logic error because
// viper.GetString("editor") is not yet set and the EDITOR env value
// will instead always be used.
// func ConfigEdit() *cobra.Command {
// 	long := fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n",
// 		fmt.Sprintf("Edit the %s configuration file.", meta.Name),
// 		"To change the editor program, either:",
// 		fmt.Sprintf("  1. Configure one by creating a %s shell environment variable.",
// 			str.Example("$EDITOR")),
// 		"  2. Set an editor in the configuration file:",
// 		str.Example(fmt.Sprintf("     %s config set --name=editor", meta.Bin)),
// 	)
// 	return &cobra.Command{
// 		Use:     "edit",
// 		Aliases: []string{"e"},
// 		Short:   "Edit the config file\n",
// 		Long:    long,
// 		GroupID: "configfile",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			if err := config.Edit(cmd.OutOrStdout()); err != nil {
// 				return err
// 			}
// 			return nil
// 		},
// 	}
// }

// func ConfigInfo() *cobra.Command {
// 	return &cobra.Command{
// 		Use:     "info",
// 		Aliases: []string{"i"},
// 		Example: fmt.Sprint(example.ConfigInfo),
// 		Short:   "List all the settings in use",
// 		Long:    fmt.Sprintf("List all the %s settings in use.", meta.Name),
// 		GroupID: "settings",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			if err := ConfigInfos(cmd.OutOrStdout()); err != nil {
// 				return err
// 			}
// 			return nil
// 		},
// 	}
// }

// func ConfigSet() *cobra.Command {
// 	return &cobra.Command{
// 		Use:     "set [setting names]",
// 		Aliases: []string{"s"},
// 		Short:   "Edit a setting",
// 		Long:    fmt.Sprintf("Edit a %s setting.", meta.Name),
// 		Example: fmt.Sprint(example.Set),
// 		GroupID: "settings",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			if err := ListAll(cmd.OutOrStdout()); err != nil {
// 				return err
// 			}
// 			if err := Usage(cmd, args...); errors.Is(err, ErrNoArgs) {
// 				return nil
// 			} else if err != nil {
// 				return err
// 			}
// 			for _, arg := range args {
// 				err := config.Set(cmd.OutOrStdout(), arg)
// 				if err != nil {
// 					return err
// 				}
// 			}
// 			return nil
// 		},
// 	}
// }

// func ConfigSetup() *cobra.Command {
// 	return &cobra.Command{
// 		Use:     "setup",
// 		Short:   "Walk through all the settings",
// 		Long:    fmt.Sprintf("Walk through all of the %s settings.", meta.Name),
// 		GroupID: "settings",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			const startAt = 0
// 			if err := config.Setup(cmd.OutOrStdout(), startAt); err != nil {
// 				return err
// 			}
// 			return nil
// 		},
// 	}
// }

// ListAll is the "config set --list" command run.
// func ListAll(w io.Writer) error {
// 	if w == nil {
// 		return ErrWriter
// 	}
// 	if !Config.Configs {
// 		return nil
// 	}
// 	if err := config.List(w); err != nil {
// 		return err
// 	}
// 	return nil
// }

///=======================

// Print the usage help and exit; but only when no arguments are given.
// func Usage(cmd *cobra.Command, args ...string) error {
// 	if len(args) == 0 {
// 		if err := cmd.Help(); err != nil {
// 			return err
// 		}
// 		return ErrNoArgs
// 	}
// 	return nil
// }

// Init reads in the config file and ENV variables if set.
// This might be triggered twice due to the Cobra initializer registers.
// func Load() {
// 	w := os.Stdout
// 	// read in environment variables
// 	viper.SetEnvPrefix("env")
// 	viper.AutomaticEnv()
// 	// tester configuration file
// 	if flag.Command.Tester {
// 		if err := LoadTester(w); err != nil {
// 			logs.Fatal(err)
// 		}
// 		return
// 	}
// 	// configuration file
// 	if err := config.Load(w, flag.Command.Config); err != nil {
// 		logs.FatalMark(viper.ConfigFileUsed(), logs.ErrConfigOpen, err)
// 	}
// }

// LoadTester loads an in-memory default configuration file for test purposes.
func LoadTester(w io.Writer) error {
	fmt.Println("The single use, in-memory tester file is in use.")
	fs := afero.NewMemMapFs()
	afs := &afero.Afero{Fs: fs}
	f, err := afs.TempFile("", "ioutil-test")
	if err != nil {
		return err
	}
	if err := config.Create(w, f.Name(), true); err != nil {
		return err
	}
	if err := config.Load(w, f.Name()); err != nil {
		return fmt.Errorf("%w, %s: %s", logs.ErrConfigOpen, err, viper.ConfigFileUsed())
	}
	return nil
}

// func ConfigInfos(w io.Writer) error {
// 	if err := config.Load(w, flag.Command.Config); err != nil {
// 		return fmt.Errorf("%w: %s", logs.ErrConfigOpen, err)
// 	}
// 	// info --configs flag
// 	if Config.Configs {
// 		if err := config.List(w); err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	// info --styles flag
// 	if Config.Styles {
// 		err := str.JSONStyles(w, fmt.Sprintf("%s info --style", meta.Bin))
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	// info --style flag
// 	style := viper.GetString("style.info")
// 	if style == "" {
// 		style = "dracula"
// 	}
// 	if Config.Style != "" {
// 		style = Config.Style
// 	}
// 	// info command
// 	if err := config.Info(w, style); err != nil {
// 		return err
// 	}
// 	return nil
// }
