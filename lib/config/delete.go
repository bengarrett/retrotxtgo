package config

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/viper"
)

// Delete a configuration file.
func Delete() {
	cfg := viper.ConfigFileUsed()
	if cfg == "" {
		configMissing(cmdPath, "delete")
	}
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		configMissing(cmdPath, "delete")
	}
	PrintLocation()
	switch logs.PromptYN("Confirm the file deletion", false) {
	case true:
		if err := os.Remove(cfg); err != nil {
			logs.Log(fmt.Errorf("config delete: could not remove %v %v", cfg, err))
		}
		fmt.Println("The config is gone")
	}
}
