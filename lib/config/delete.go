package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/prompt"
)

// Delete a configuration file.
func Delete() (err logs.Generic) {
	cfg := viper.ConfigFileUsed()
	if cfg == "" {
		configMissing(cmdPath, "delete")
	}
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		configMissing(cmdPath, "delete")
	} else if err != nil {
		return logs.Generic{
			Issue: "failed to stat the config file",
			Err:   err,
		}
	}
	PrintLocation()
	if prompt.YesNo("Confirm the configuration file deletion", false) {
		if err := os.Remove(cfg); err != nil {
			return logs.Generic{
				Issue: "failed to remove config file",
				Err:   err,
			}
		}
		fmt.Println("The config is gone")
	}
	return logs.Generic{}
}
