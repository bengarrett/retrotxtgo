package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/prompt"
)

// Delete a configuration file.
func Delete() (err logs.Argument) {
	cfg := viper.ConfigFileUsed()
	if cfg == "" {
		configMissing(cmdPath, "delete")
		os.Exit(1)
	}
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		configMissing(cmdPath, "delete")
		os.Exit(1)
	} else if err != nil {
		return logs.Argument{
			Issue: "failed to stat the config file",
			Err:   err,
		}
	}
	PrintLocation()
	if prompt.YesNo("Confirm the configuration file deletion", false) {
		if err := os.Remove(cfg); err != nil {
			return logs.Argument{
				Issue: "failed to remove config file",
				Err:   err,
			}
		}
		fmt.Println("The config is gone")
	}
	return logs.Argument{}
}
