package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/viper"
)

// InitDefaults initialises flag and configuration defaults.
func InitDefaults() {
	viper.SetDefault("create.layout", "standard")
	viper.SetDefault("create.title", "RetroTxt | example")
	viper.SetDefault("create.meta.author", "")
	viper.SetDefault("create.meta.color-scheme", "")
	viper.SetDefault("create.meta.description", "An example")
	viper.SetDefault("create.meta.generator", true)
	viper.SetDefault("create.meta.keywords", "")
	viper.SetDefault("create.meta.referrer", "")
	viper.SetDefault("create.meta.theme-color", "")
	viper.SetDefault("create.save-directory", "")
	viper.SetDefault("create.server-port", 8080)
	viper.SetDefault("info.format", "color")
	viper.SetDefault("version.format", "color")
}

// SetConfig reads and loads a configuration file.
func SetConfig(configFlag string) {
	viper.SetConfigType("yaml")
	if configFlag == "" {
		configFlag = Filepath()
	}
	viper.SetConfigFile(configFlag)

	if err := viper.ReadInConfig(); err != nil {
		var e = logs.ConfigErr{
			FileUsed: viper.ConfigFileUsed(),
			Err:      err,
		}
		if errors.Is(err, os.ErrNotExist) {
			if configFlag != "" {
				// require manual generation for custom config files
				e.Err = errors.New("does not exist\n\t use the command: retrotxt config create --config")
				fmt.Println(e.String())
				os.Exit(1)
			} else {
				// auto-generate a new, default config file
				Create(viper.ConfigFileUsed(), false)
			}
		} else {
			// config fail
			fmt.Println(e.String())
			os.Exit(1)
		}
	} else if configFlag != "" {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
