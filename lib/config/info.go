package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
)

// Info prints the content of a configuration file.
func Info(style string) error {
	fmt.Println(str.Cp(fmt.Sprintf("%s default configurations when no flags are given.", meta.Name)))
	PrintLocation()
	out, err := json.MarshalIndent(Enabled(), "", " ")
	if err != nil {
		return fmt.Errorf("failed to read configuration in yaml syntax: %w", err)
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
		err = str.Highlight(string(out), "json", style, true)
		if err != nil {
			return fmt.Errorf("failed to run highlighter: %w", err)
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
		const tries = 5
		if len(m) < tries {
			fmt.Print(str.Example(fmt.Sprintf("%s config set %s\n", meta.Bin, strings.Join(m, " "))))
		} else {
			fmt.Print(str.Example(fmt.Sprintf("%s config set setup\n", meta.Bin)))
		}
	}
	return nil
}
