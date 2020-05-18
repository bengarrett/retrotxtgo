package config

import (
	"fmt"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Info prints the content of a configuration file.
func Info() {
	println(logs.Cp("These are the default configurations used by the commands of RetroTxt when no flags are given.\n"))
	sets, err := yaml.Marshal(viper.AllSettings())
	logs.Check("read configuration yaml", err)
	if err := quick.Highlight(os.Stdout, string(sets), "yaml", "terminal256", infoStyles); err != nil {
		fmt.Println(string(sets))
	}
	println()
}
