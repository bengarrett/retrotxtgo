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
		configExit(cmdPath, "delete")
	}
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		configExit(cmdPath, "delete")
	}
	switch logs.PromptYN("Confirm the file deletion", false) {
	case true:
		if err := os.Remove(cfg); err != nil {
			logs.Save(fmt.Errorf("config delete: could not remove %v %v", cfg, err))
		}
		fmt.Println("The config is gone")
	}
}
