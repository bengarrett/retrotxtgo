package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Info prints the content of a configuration file.
func Info(style string) (err logs.IssueErr) {
	fmt.Println(str.Cp("RetroTxt default configurations when no flags are given."))
	PrintLocation()
	if m := missing(); len(m) > 0 {
		fmt.Printf("These settings are missing should be configured: %s\n", strings.Join(m, ", "))
	}

	// ignore unknown or invalid settings in the config file
	// TODO: break into a func and create a count func
	// if config is missing settings. list those settings and their values at the bottom of config info
	var used = make(map[string]interface{})
	for _, key := range viper.AllKeys() {
		if d := Defaults[key]; d != nil {
			used[key] = viper.Get(key)
		}
	}
	out, e := yaml.Marshal(used)
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
		str.Highlight(string(out), "yaml", style)
	}
	return logs.IssueErr{}
}

func missing() (list []string) {
	d, l := len(Defaults), len(viper.AllSettings())
	if d == l {
		return list
	}
	for key := range Defaults {
		if !viper.IsSet(key) {
			list = append(list, key)
		}
	}
	return list
}
