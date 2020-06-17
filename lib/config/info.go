package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/gookit/color"
	"github.com/spf13/viper"
)

// Info prints the content of a configuration file.
func Info(style string) (err logs.IssueErr) {
	fmt.Println(str.Cp("RetroTxt default configurations when no flags are given."))
	PrintLocation()
	out, e := json.MarshalIndent(Enabled(), "", " ")
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
		str.Highlight(string(out), "json", style)
	}
	if m := Missing(); len(m) > 0 {
		s := "These settings are missing and should be configured"
		if len(m) == 1 {
			s = "This setting is missing and should be configured"
		}
		fmt.Printf("\n\n%s:\n%s\n\n%s:\n", color.Warn.Sprint(s),
			strings.Join(m, ", "), color.Warn.Sprint("To repair"))
		if len(m) < 5 {
			fmt.Print(str.Example(fmt.Sprintf("retrotxt config set %s\n", strings.Join(m, " "))))
		} else {
			fmt.Print(str.Example(fmt.Sprintf("retrotxt config set setup\n")))
		}
	}
	return logs.IssueErr{}
}

// Enabled returns all the Viper keys holding a value that are used by RetroTxt.
// This will hide all unrecognised manual edits to the configuration file.
func Enabled() map[string]interface{} {
	var sets = make(map[string]interface{})
	for _, key := range viper.AllKeys() {
		if d := Defaults[key]; d != nil {
			sets[key] = viper.Get(key)
		}
	}
	return sets
}

// Missing lists the settings that are not found in the configuration file.
// This could be due to new features being added after the file was generated
// or because of manual edits.
func Missing() (list []string) {
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
