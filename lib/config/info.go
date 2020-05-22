package config

import (
	"bytes"
	"fmt"

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
	out = bytes.ReplaceAll(out, []byte("    "), []byte("  "))
	switch style {
	case "none", "":
		fmt.Println(string(out))
	default:
		logs.Highlight(string(out), "yaml", style)
	}
	return logs.IssueErr{}
}
