package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/viper"
)

// Create a configuration file.
func Create(ow bool) {
	if cfg := viper.ConfigFileUsed(); cfg != "" && !ow {
		if _, err := os.Stat(cfg); !os.IsNotExist(err) {
			configExists(cmdPath, "create")
		}
		p := filepath.Dir(cfg)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			fmt.Println(p)
			if err := os.MkdirAll(p, permDir); err != nil {
				logs.ChkErr("", err)
				os.Exit(exit + 2)
			}
		}
	}
	writeConfig(false)
}

func configExists(name string, suffix string) {
	cmd := strings.TrimSuffix(name, suffix)
	fmt.Printf("A config file already is in use at: %s\n",
		logs.Cf(viper.ConfigFileUsed()))
	fmt.Printf("To edit it: %s\n", logs.Cp(cmd+"edit"))
	fmt.Printf("To delete:  %s\n", logs.Cp(cmd+"delete"))
	os.Exit(exit + 1)
}
