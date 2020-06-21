package config

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

var setupMode = false

// Setup walks through all the settings and saves them to the configuration file.
func Setup() {
	setupMode = true
	prompt.SetupMode = true
	keys := Keys()
	PrintLocation()
	var w uint = 42
	for i, key := range keys {
		if i == 0 {
			fmt.Printf("\n%s\n\n", str.Cb(enterKey()))
		}
		h := fmt.Sprintf("  %d/%d. RetroTxt Setup - %s",
			i+1, len(keys), key)
		if i == 0 {
			fmt.Println(hr(w))
		}
		fmt.Println(h)
		Update(key)
		fmt.Println(hr(w))
	}
	fmt.Println(Info(viper.GetString("style.info")))
}

func enterKey() string {
	if runtime.GOOS == "darwin" {
		return "Press ↩ return to skip the question or ⌃ control-c to quit"
	}
	return "Press ⏎ enter to skip the question or Ctrl-c to quit"
}

func hr(w uint) string {
	return str.Cb(strings.Repeat("-", int(w)))
}
