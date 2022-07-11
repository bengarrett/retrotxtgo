package root

import (
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/viper"
)

// Init reads in the config file and ENV variables if set.
// This might be triggered twice due to the Cobra initializer registers.
func Init() {
	// read in environment variables
	viper.SetEnvPrefix("env")
	viper.AutomaticEnv()
	// configuration file
	if err := config.SetConfig(flag.RootFlag.Config); err != nil {
		logs.FatalMark(viper.ConfigFileUsed(), logs.ErrConfigOpen, err)
	}
}
