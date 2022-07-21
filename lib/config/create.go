package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

var ErrExist = errors.New("a config file already exists")

// New creates a new configuration file and prints the results.
func New(w io.Writer, overwrite bool) error {
	if err := Create(w, viper.ConfigFileUsed(), overwrite); err != nil {
		return err
	}
	fmt.Fprintf(w, "%s%s\n%s%s\n",
		str.Info(), "A new file was created.",
		"Config file: ", str.Path(viper.ConfigFileUsed()))
	return nil
}

// Create a named configuration file with the option to overwrite any existing files.
func Create(w io.Writer, name string, ow bool) error {
	if name == "" {
		return fmt.Errorf("create configuration: %w", logs.ErrNameNil)
	}
	_, err := os.Stat(name)
	switch {
	case !ow && !os.IsNotExist(err):
		return ErrExist
	case os.IsNotExist(err):
		return createNew(w, name)
	case err != nil:
		return fmt.Errorf("%s: %q: %w", errMsg("access"), name, err)
	}
	return createNew(w, name)
}

// createNew creates a new config file.
func createNew(w io.Writer, name string) error {
	path, err := filesystem.Touch(name)
	if err != nil {
		return fmt.Errorf("%s: %q: %w", errMsg("create"), name, err)
	}
	InitDefaults()
	err = Save(nil, path)
	if err != nil {
		return fmt.Errorf("%s: %q: %w", errMsg("update"), name, err)
	}
	return nil
}

// DoesExist prints a how-to message when a config file already exists.
func DoesExist(w io.Writer, name, suffix string) {
	example := func(s string) string {
		x := fmt.Sprintf("%s %s", strings.TrimSuffix(name, suffix), s)
		return str.Example(x)
	}
	fmt.Fprintf(w, "%sA config file already exists: %s\n%s\n",
		str.Info(), str.ColFuz(viper.ConfigFileUsed()),
		"Use the following commands to modify it.")
	fmt.Fprintf(w, "Edit:\t%s\n", example("edit"))
	fmt.Fprintf(w, "Delete:\t%s\n", example("delete"))
	fmt.Fprintf(w, "Reset:\t%s\n", example("create --overwrite"))
}

func errMsg(verb string) string {
	return fmt.Sprintf("could not %s the configuration file", verb)
}
