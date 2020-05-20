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
func Info() (err logs.IssueErr) {
	fmt.Println(logs.Cp("RetroTxt default configurations when no flags are given."))
	PrintLocation()
	sets, e := yaml.Marshal(viper.AllSettings())
	if e != nil {
		return logs.IssueErr{
			Issue: "failed to read configuration in yaml syntax",
			Err:   e,
		}
	}
	if err := quick.Highlight(os.Stdout, string(sets), "yaml", "terminal256", infoStyles); err != nil {
		fmt.Println(string(sets))
	}
	return logs.IssueErr{}
}
