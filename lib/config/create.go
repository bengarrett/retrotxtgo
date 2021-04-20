package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/str"
)

// Create a named configuration file.
func Create(name string, ow bool) (err error) {
	if name == "" {
		return fmt.Errorf("create configuration: %w", ErrNoFName)
	}
	_, err = os.Stat(name)
	switch {
	case !os.IsNotExist(err) && !ow:
		configDoesExist(cmdPath, "create")
		os.Exit(1)
	case os.IsNotExist(err):
		// a missing named file is okay
	case err != nil:
		return fmt.Errorf("could not access the configuration file: %q: %w", name, err)
	}
	// create a new config file
	path, err := filesystem.Touch(name)
	if err != nil {
		return fmt.Errorf("could not create the configuration file: %q: %w", name, err)
	}
	InitDefaults()
	err = UpdateConfig(path, false)
	if err != nil {
		return fmt.Errorf("could not update the configuration file: %q: %w", name, err)
	}
	return nil
}

func configDoesExist(name, suffix string) {
	cmd := strings.TrimSuffix(name, suffix)
	fmt.Printf("%s a config file exists: %s\n",
		str.Info(), str.Cf(viper.ConfigFileUsed()))
	fmt.Printf(" edit it: %s\n", str.Cp(cmd+" edit"))
	fmt.Printf("  delete: %s\n", str.Cp(cmd+" delete"))
	fmt.Printf("   reset: %s\n", str.Cp(cmd+" create --overwrite"))
}
