package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/upd"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

// New creates a new configuration file and prints the results.
func New(overwrite bool) error {
	if err := Create(viper.ConfigFileUsed(), overwrite); err != nil {
		return err
	}
	fmt.Printf("%s%s\n%s%s\n",
		str.Info(), "A new file was created.",
		"Config file: ", str.Path(viper.ConfigFileUsed()))
	return nil
}

// Create a named configuration file with the option to overwrite any existing files.
func Create(name string, ow bool) error {
	if name == "" {
		return fmt.Errorf("create configuration: %w", logs.ErrNameNil)
	}
	_, err := os.Stat(name)
	switch {
	case !ow && !os.IsNotExist(err):
		configDoesExist(cmdPath(), "create")
		os.Exit(1)
	case os.IsNotExist(err):
		return createNew(name)
	case err != nil:
		return fmt.Errorf("%s: %q: %w", errMsg("access"), name, err)
	}
	return createNew(name)
}

// createNew creates a new config file.
func createNew(name string) error {
	path, err := filesystem.Touch(name)
	if err != nil {
		return fmt.Errorf("%s: %q: %w", errMsg("create"), name, err)
	}
	InitDefaults()
	err = upd.UpdateConfig(path, false)
	if err != nil {
		return fmt.Errorf("%s: %q: %w", errMsg("update"), name, err)
	}
	return nil
}

// configDoesExist prints a how-to message when a config file already exists.
func configDoesExist(name, suffix string) {
	example := func(s string) string {
		x := fmt.Sprintf("%s %s", strings.TrimSuffix(name, suffix), s)
		return str.Example(x)
	}
	fmt.Printf("%sA config file already exists: %s\n%s\n",
		str.Info(), str.ColFuz(viper.ConfigFileUsed()),
		"Use the following commands to modify it.")
	fmt.Printf("Edit:\t%s\n", example("edit"))
	fmt.Printf("Delete:\t%s\n", example("delete"))
	fmt.Printf("Reset:\t%s\n", example("create --overwrite"))
}

func errMsg(verb string) string {
	return fmt.Sprintf("could not %s the configuration file", verb)
}
