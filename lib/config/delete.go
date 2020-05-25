package config

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/spf13/viper"
)

// Delete a configuration file.
func Delete() (err logs.IssueErr) {
	cfg := viper.ConfigFileUsed()
	if cfg == "" {
		configMissing(cmdPath, "delete")
	}
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		configMissing(cmdPath, "delete")
	} else if err != nil {
		return logs.IssueErr{
			Issue: "failed to stat the config file",
			Err:   err,
		}
	}
	PrintLocation()
	switch prompt.YesNo("Confirm the configuration file deletion", false) {
	case true:
		if err := os.Remove(cfg); err != nil {
			return logs.IssueErr{
				Issue: "failed to remove config file",
				Err:   err,
			}
		}
		fmt.Println("The config is gone")
	}
	return logs.IssueErr{}
}
