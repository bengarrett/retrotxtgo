package configcmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/long"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/usage"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/cobra"
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
		return create()
	case Delete:
		return delete()
	case Edit:
		return edit()
	case Info:
		return info()
	case Set:
		return set()
	case Setup:
		return setup()
	}
	return nil
}

func create() *cobra.Command {
	return &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create or reset the config file",
		Long:    fmt.Sprintf("Create or reset the %s configuration file.", meta.Name),
		Run: func(cmd *cobra.Command, args []string) {
			if err := config.New(flag.Config.Ow); err != nil {
				logs.FatalWrap(logs.ErrConfigNew, err)
			}
		},
	}
}

func delete() *cobra.Command {
	return &cobra.Command{
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
}

// Note: Previously I inserted the results of config.Editor() into
// the Short and Long fields. This will cause a logic error because
// viper.GetString("editor") is not yet set and the EDITOR env value
// will instead always be used.
func edit() *cobra.Command {
	return &cobra.Command{
		Use:     "edit",
		Aliases: []string{"e"},
		Short:   "Edit the config file\n",
		Long:    long.ConfigEdit,
		Run: func(cmd *cobra.Command, args []string) {
			if err := config.Edit(); err != nil {
				logs.Fatal(err)
			}
		},
	}
}

func info() *cobra.Command {
	return &cobra.Command{
		Use:     "info",
		Aliases: []string{"i"},
		Example: example.CfgInfo.Print(),
		Short:   "List all the settings in use",
		Long:    fmt.Sprintf("List all the %s settings in use.", meta.Name),
		Run: func(cmd *cobra.Command, args []string) {
			if flag.ConfigInfo() {
				return
			}
		},
	}
}

func set() *cobra.Command {
	return &cobra.Command{
		Use:     "set [setting names]",
		Aliases: []string{"s"},
		Short:   "Edit a setting",
		Long:    fmt.Sprintf("Edit a %s setting.", meta.Name),
		Example: example.Set.Print(),
		Run: func(cmd *cobra.Command, args []string) {
			if ListAll() {
				return
			}
			if err := usage.Print(cmd, args...); err != nil {
				logs.Fatal(err)
			}
			for _, arg := range args {
				config.Set(arg)
			}
		},
	}
}

func setup() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Walk through all the settings",
		Long:  fmt.Sprintf("Walk through all of the %s settings.", meta.Name),
		Run: func(cmd *cobra.Command, args []string) {
			const startAt = 0
			config.Setup(startAt)
		},
	}
}

// ListAll is the "config set --list" command run.
func ListAll() (exit bool) {
	if flag.Config.Configs {
		if err := config.List(); err != nil {
			logs.FatalFlag("config", "list", err)
		}
		return true
	}
	return false
}
