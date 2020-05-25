package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/str"
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
	InitDefaults()
	return UpdateConfig(path, false)
}

func configDoesExist(name string, suffix string) {
	cmd := strings.TrimSuffix(name, suffix)
	fmt.Printf("%s a config file exists: %s\n",
		str.Info(), str.Cf(viper.ConfigFileUsed()))
	fmt.Printf(" edit it: %s\n", str.Cp(cmd+" edit"))
	fmt.Printf("  delete: %s\n", str.Cp(cmd+" delete"))
	fmt.Printf("   reset: %s\n", str.Cp(cmd+" create --overwrite"))
	os.Exit(20)
}
