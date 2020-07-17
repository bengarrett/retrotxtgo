package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gookit/color"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

// Info prints the content of a configuration file.
func Info(style string) (err logs.Generic) {
	fmt.Println(str.Cp("RetroTxt default configurations when no flags are given."))
	PrintLocation()
	out, e := json.MarshalIndent(Enabled(), "", " ")
	if e != nil {
		return logs.Generic{
			Issue: "failed to read configuration in yaml syntax",
			Err:   e,
		}
	}
	switch style {
	case "none", "":
		fmt.Println(string(out))
	default:
		if !str.Valid(style) {
			fmt.Printf("unknown style %q, so using none\n", style)
			fmt.Println(string(out))
			break
		}
		e = str.Highlight(string(out), "json", style)
		if e != nil {
			return logs.Generic{
				Issue: "failed to run highligher",
				Err:   e,
			}
		}
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
			fmt.Print(str.Example("retrotxt config set setup\n"))
		}
	}
	return logs.Generic{}
}
