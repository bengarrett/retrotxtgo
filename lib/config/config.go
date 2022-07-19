// Package config handles the user configations.
package config

import (
	"errors"
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/set"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/viper"
)

var (
	ErrEditorNil = errors.New("no suitable text editor can be found")
	ErrEditorRun = errors.New("editor cannot be run")
	ErrLogo      = errors.New("program logo is missing")
	ErrSaveType  = errors.New("save value type is unsupported")
)

const namedFile = "config.yaml"

func CmdPath() string {
	return fmt.Sprintf("%s config", meta.Bin)
}

// Tip provides a brief help on the config file configurations.
func Tip() get.Hints {
	return get.Tip()
}

// Formats choices for flags.
type Formats struct {
	Info [5]string
}

// Format flag choices for the info command.
func Format() Formats {
	return Formats{
		Info: [5]string{"color", "json", "json.min", "text", "xml"},
	}
}

// Enabled returns all the Viper keys holding a value that are used.
// This will hide all unrecognized manual edits to the configuration file.
func Enabled() map[string]interface{} {
	sets := make(map[string]interface{})
	for _, key := range viper.AllKeys() {
		if d := get.Reset()[key]; d != nil {
			sets[key] = viper.Get(key)
		}
	}
	return sets
}

// KeySort list all the available configuration setting names sorted by hand.
func KeySort() []string {
	all := set.Keys()
	keys := []string{
		get.FontFamily, get.Title, get.LayoutTmpl, get.FontEmbed,
		get.SaveDir, get.Serve, get.Editor, get.Styleh, get.Stylei,
	}
	for _, key := range all {
		found := false
		for _, used := range keys {
			if key == used {
				found = true
				break
			}
		}
		if !found {
			keys = append(keys, key)
		}
	}
	return keys
}

// Missing returns the settings that are not found in the configuration file.
// This could be due to new features being added after the file was generated
// or because of manual file edits.
func Missing() []string {
	list := []string{}
	if len(get.Reset()) == len(viper.AllSettings()) {
		return list
	}
	for key := range get.Reset() {
		if !viper.IsSet(key) {
			list = append(list, key)
		}
	}
	return list
}
