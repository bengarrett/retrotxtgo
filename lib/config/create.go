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
}

func configExist(name string, suffix string) {
	cmd := strings.TrimSuffix(name, suffix)
	fmt.Printf("A config file exists: %s\n",
		logs.Cf(viper.ConfigFileUsed()))
	fmt.Printf("to edit it:\t%s\n", logs.Cp(cmd+" edit"))
	fmt.Printf("   delete:\t%s\n", logs.Cp(cmd+" delete"))
	fmt.Printf("   reset:\t%s\n", logs.Cp(cmd+" create --overwrite"))
	os.Exit(exit + 1)
}
