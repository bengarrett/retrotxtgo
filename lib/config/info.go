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
	fmt.Printf("%s%s\n%s%s\n\n", str.Info(),
		Location(),
		meta.Name, " default settings in use.",
	)
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
	}
	if s := missing(Missing()...); s != "" {
		fmt.Print(s)
	}
	return nil
}

// missing returns a printed list of missing settings in the config file.
func missing(list ...string) string {
	const tries = 5
	l := len(list)
	if l == 0 {
		return ""
	}
	t := "These settings are missing and should be configured"
	if l == 1 {
		t = "This setting is missing and should be configured"
	}
	var s string
	if l < tries {
		s = str.Example(fmt.Sprintf("%s config set %s\n",
			meta.Bin, strings.Join(list, " ")))
	} else {
		s = str.Example(fmt.Sprintf("%s config set setup\n",
			meta.Bin))
	}
	return fmt.Sprintf("\n\n%s:\n%s\n\n%s:\n%s",
		color.Warn.Sprint(t),
		strings.Join(list, ", "),
		color.Warn.Sprint("To repair"), s)
}
