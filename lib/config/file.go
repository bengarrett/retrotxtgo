package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

type configErr struct {
	FileUsed string
	Err      error
}

func (e configErr) String() string {
	return (logs.Err{
		Issue: "config file",
		Arg:   e.FileUsed,
		Msg:   e.Err}).String()
}

// InitDefaults initialises flag and configuration defaults.
func InitDefaults() {
	for key, val := range Defaults {
		viper.SetDefault(key, val)
		viper.Set(key, val)
	}
}

// SetConfig reads and loads a configuration file.
func SetConfig(configFlag string) {
	cfgExit := func(e configErr) {
		// require manual generation for custom config files
		e.Err = errors.New("does not exist\n\t use the command: retrotxt config create --config=" + configFlag)
		fmt.Println(e.String())
		os.Exit(1)
	}
	viper.SetConfigType("yaml")
	var configPath = Path()
	if configFlag != "" {
		configPath = configFlag
	}
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		var e = configErr{
			FileUsed: viper.ConfigFileUsed(),
			Err:      err,
		}
		if errors.Is(err, os.ErrNotExist) {
			// initialise a new, default config file if conditions are met
			switch {
			case len(os.Args) > 2:
				switch strings.Join(os.Args[1:3], ".") {
				case "config.create", "config.delete":
					// never auto-generate a config when these arguments are given
					return
				}
				if configFlag != "" {
					cfgExit(e)
				}
			case configFlag != "":
				cfgExit(e)
			}
			// auto-generate new config, except when --config flag is used
			Create(viper.ConfigFileUsed(), false)
		} else {
			// config fail
			fmt.Println(e.String())
			os.Exit(1)
		}
	} else if configFlag != "" {
		// always print the config location when the --config flag is used
		if len(os.Args) > 0 && os.Args[1] == "config" {
			// except when the config command is in use
			return
		}
		PrintLocation()
	}
}

// PrintLocation prints the location of the current configuration file.
func PrintLocation() {
	fmt.Println(
		str.Cb(fmt.Sprintf("Config file: %s",
			viper.ConfigFileUsed())))
}
