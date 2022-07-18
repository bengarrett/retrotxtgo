package config

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

// Delete a configuration file.
func Delete() error {
	name := viper.ConfigFileUsed()
	if name == "" {
		configMissing(CmdPath(), "delete")
		os.Exit(1)
	}
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		configMissing(CmdPath(), "delete")
		os.Exit(1)
	}
	if err != nil {
		return fmt.Errorf("failed to access the config file: %w", err)
	}
	fmt.Printf("%s%s", str.Confirm(), Location())
	if prompt.YesNo("Delete this file", false) {
		if err := os.Remove(name); err != nil {
			return fmt.Errorf("failed to remove config file: %w", err)
		}
		fmt.Println("File is deleted")
	}
	return nil
}
