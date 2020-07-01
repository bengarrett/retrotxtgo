package config

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/internal/pack"
	"retrotxt.com/retrotxt/lib/prompt"
	"retrotxt.com/retrotxt/lib/str"
)

var setupMode = false

// Setup walks through all the settings and saves them to the configuration file.
func Setup() {
	setupMode, prompt.SetupMode = true, true
	keys := Keys()
	logo()
	PrintLocation()
	var w uint = 80
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

func logo() {
	n := "text/retrotxt.utf8ans"
	b := pack.Get(n)
	if b == nil {
		log.Fatal(errors.New("pkg.name is unknown: " + n))
	}
	// convert and print
	fmt.Println(string(b))
}
