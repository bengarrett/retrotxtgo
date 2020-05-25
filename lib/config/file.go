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
	home, err := os.UserHomeDir()
	if err != nil {
		logs.LogCont(err)
	}
	viper.SetDefault("create.layout", "standard")
	viper.SetDefault("create.meta.author", "")
	viper.SetDefault("create.meta.color-scheme", "")
	viper.SetDefault("create.meta.description", "An example")
	viper.SetDefault("create.meta.generator", true)
	viper.SetDefault("create.meta.keywords", "")
	viper.SetDefault("create.meta.referrer", "")
	viper.SetDefault("create.meta.theme-color", "")
	viper.SetDefault("create.save-directory", home)
	viper.SetDefault("create.server-port", 8080)
	viper.SetDefault("create.title", "RetroTxt | example")
	viper.SetDefault("editor", "")
	viper.SetDefault("style.html", "lovelace")
	viper.SetDefault("style.yaml", "monokai")
}

// SetConfig reads and loads a configuration file.
func SetConfig(configFlag string) {
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
			if configFlag != "" {
				// require manual generation for custom config files
				e.Err = errors.New("does not exist\n\t use the command: retrotxt config create --config=" + configFlag)
				fmt.Println(e.String())
				os.Exit(1)
			} else if len(os.Args) > 2 {
				switch strings.Join(os.Args[1:3], ".") {
				case "config.create", "config.delete":
					return
				}
				// auto-generate a new, default config file
				Create(viper.ConfigFileUsed(), false)
			}
		} else {
			// config fail
			fmt.Println(e.String())
			os.Exit(1)
		}
	} else if configFlag != "" {
		PrintLocation()
	}
}

// PrintLocation prints the location of the current configuration file.
func PrintLocation() {
	fmt.Println(
		str.Cb(fmt.Sprintf("Config file: %s",
			viper.ConfigFileUsed())))
}
