package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/viper"
)

// Create a named configuration file.
func Create(name string, ow bool) (err error) {
	if name == "" {
		return errors.New("config create: name value cannot be empty")
	}
	if _, err := os.Stat(name); !os.IsNotExist(err) && !ow {
		configExist(cmdPath, "create")
		//return errors.New("config create: named file exists but overwrite (ow) is false")
	}
	path, err := filesystem.Touch(name)
	if err != nil {
		return err
	}
	err = UpdateConfig(path, true)
	return err
	//dir := filepath.Dir(name)
	// if cfg := viper.ConfigFileUsed(); cfg != "" && !ow {
	// 	if _, err := os.Stat(cfg); !os.IsNotExist(err) {
	// 		configExist(cmdPath, "create")
	// 	}
	// 	p := filepath.Dir(cfg)
	// 	if _, err := os.Stat(p); os.IsNotExist(err) {
	// 		fmt.Println(p)
	// 		if err := os.MkdirAll(p, permDir); err != nil {
	// 			logs.Check("", err)
	// 			os.Exit(exit + 2)
	// 		}
	// 	}
	// }
	//writeConfig(false) // TODO: return err?
}

func configExist(name string, suffix string) {
	cmd := strings.TrimSuffix(name, suffix)
	fmt.Printf("A config file already is in use: %s\n",
		logs.Cf(viper.ConfigFileUsed()))
	fmt.Printf("To edit it: %s\n", logs.Cp(cmd+"edit"))
	fmt.Printf("To delete:  %s\n", logs.Cp(cmd+"delete"))
	os.Exit(exit + 1)
}
