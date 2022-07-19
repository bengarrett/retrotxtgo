package config

import (
	"fmt"
	"io"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

// Delete a configuration file.
func Delete(w io.Writer, ask bool) error {
	name := viper.ConfigFileUsed()
	if name == "" {
		configMissing(w, CmdPath(), "delete")
		return nil
	}
	f, err := os.Stat(name)
	if os.IsNotExist(err) {
		configMissing(w, CmdPath(), "delete")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to access the config file: %w", err)
	}
	if !ask {
		fmt.Fprint(w, Location())
		fmt.Fprintf(w, "%d bytes\n", f.Size())
	} else {
		fmt.Fprintf(w, "%s%s", str.Confirm(), Location())
		if !prompt.YesNo(w, "Delete this file", false) {
			fmt.Fprintln(w, "Deletion cancelled")
			return nil
		}
	}
	if err := os.Remove(name); err != nil {
		return fmt.Errorf("failed to remove config file: %w", err)
	}
	fmt.Fprintln(w, "File is deleted")
	return nil
}
