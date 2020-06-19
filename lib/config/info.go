package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/gookit/color"
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
		fmt.Println()
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
