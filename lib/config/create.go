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
		configDoesExist(cmdPath, "create")
	} else if os.IsNotExist(err) {
		// a missing named file is okay
	} else if err != nil {
		return err
	}
	// create a new config file
	path, err := filesystem.Touch(name)
	if err != nil {
		return err
	}
	return UpdateConfig(path, false)
}

func configDoesExist(name string, suffix string) {
	cmd := strings.TrimSuffix(name, suffix)
	fmt.Printf("%s a config file exists: %s\n",
		logs.Info(), logs.Cf(viper.ConfigFileUsed()))
	fmt.Printf(" edit it: %s\n", logs.Cp(cmd+" edit"))
	fmt.Printf("  delete: %s\n", logs.Cp(cmd+" delete"))
	fmt.Printf("   reset: %s\n", logs.Cp(cmd+" create --overwrite"))
	os.Exit(exit + 1)
}
