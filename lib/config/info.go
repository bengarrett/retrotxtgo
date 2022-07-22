package config

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
)

// Info prints the content of a configuration file.
func Info(w io.Writer, style string) error {
	if w == nil {
		return ErrWriter
	}
	fmt.Fprintf(w, "%s%s\n%s%s\n\n", str.Info(),
		Location(),
		meta.Name, " default settings in use.",
	)
	out, err := json.MarshalIndent(Enabled(), "", " ")
	if err != nil {
		return fmt.Errorf("failed to read configuration in yaml syntax: %w", err)
	}
	switch style {
	case "none", "":
		fmt.Fprintln(w, string(out))
	default:
		if !str.Valid(style) {
			fmt.Fprintf(w, "unknown style %q, so using none\n", style)
			fmt.Fprintln(w, string(out))
			break
		}
		err = str.HighlightWriter(w, string(out), "json", style, true)
		if err != nil {
			return fmt.Errorf("failed to run highlighter: %w", err)
		}
	}
	return Alert(w, Missing()...)
}

// Alert returns writes list of missing settings in the config file.
func Alert(w io.Writer, list ...string) error {
	if w == nil {
		return ErrWriter
	}
	const tries = 5
	l := len(list)
	if l == 0 {
		return nil
	}
	t := "These settings are missing and should be configured"
	if l == 1 {
		t = "This setting is missing and should be configured"
	}
	s := str.Example(fmt.Sprintf("%s config set setup\n", meta.Bin))
	if l < tries {
		s = str.Example(fmt.Sprintf("%s config set %s\n",
			meta.Bin, strings.Join(list, " ")))
	}
	fmt.Fprintf(w, "\n\n%s:\n%s\n\n%s:\n%s",
		color.Warn.Sprint(t),
		strings.Join(list, ", "),
		color.Warn.Sprint("To repair"), s)
	return nil
}
