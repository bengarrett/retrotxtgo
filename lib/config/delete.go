package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

// Delete a configuration file.
func Delete(ask bool) (*bytes.Buffer, error) {
	w := new(bytes.Buffer)
	name := viper.ConfigFileUsed()
	if name == "" {
		w = configMissing(w, CmdPath(), "delete")
		return w, nil
	}
	f, err := os.Stat(name)
	if os.IsNotExist(err) {
		w = configMissing(w, CmdPath(), "delete")
		return w, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to access the config file: %w", err)
	}
	if !ask {
		fmt.Fprint(w, Location())
		fmt.Fprintf(w, "%d bytes\n", f.Size())
	} else {
		fmt.Fprintf(w, "%s%s", str.Confirm(), Location())
		w.WriteTo(os.Stdout)
		if !prompt.YesNo("Delete this file", false) {
			fmt.Fprintln(w, "Deletion cancelled")
			return w, nil
		}
	}
	if err := os.Remove(name); err != nil {
		return nil, fmt.Errorf("failed to remove config file: %w", err)
	}
	fmt.Fprintln(w, "File is deleted")
	return w, nil
}
