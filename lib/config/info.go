package config

import (
	"fmt"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Info prints the content of a configuration file.
func Info(style string) (err logs.IssueErr) {
	fmt.Println(logs.Cp("RetroTxt default configurations when no flags are given."))
	PrintLocation()
	out, e := yaml.Marshal(viper.AllSettings())
	if e != nil {
		return logs.IssueErr{
			Issue: "failed to read configuration in yaml syntax",
			Err:   e,
		}
	}
	switch style {
	case "none", "":
		fmt.Println(string(out))
	default:
		if err := quick.Highlight(os.Stdout, string(out), "yaml", "terminal256", style); err != nil {
			fmt.Println(string(out))
		}
	}
	return logs.IssueErr{}
}
