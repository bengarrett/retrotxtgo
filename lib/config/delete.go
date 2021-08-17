package config

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/spf13/viper"
)

// Delete a configuration file.
func Delete() error {
	cfg := viper.ConfigFileUsed()
	if cfg == "" {
		configMissing(cmdPath(), "delete")
		os.Exit(1)
	}
	var err error
	if _, err = os.Stat(cfg); os.IsNotExist(err) {
		configMissing(cmdPath(), "delete")
		os.Exit(1)
	}
	if err != nil {
		return fmt.Errorf("failed to stat the config file: %w", err)
	}
	PrintLocation()
	if prompt.YesNo("Confirm the configuration file deletion", false) {
		if err := os.Remove(cfg); err != nil {
			return fmt.Errorf("failed to remove config file: %w", err)
		}
		fmt.Println("The config is gone")
	}
	return nil
}
