package config

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
	"retrotxt.com/retrotxt/static"
)

var ErrLogo = errors.New("retrotxt logo is missing")

// Setup walks through all the settings and saves them to the configuration file.
func Setup() {
	keys := Keys()
	logo()
	PrintLocation()
	var width uint = 80
	watch()
	for i, key := range keys {
		if i == 0 {
			fmt.Printf("\n %s\n\n", str.Cinf(enterKey()))
		}
		h := fmt.Sprintf("  %d/%d. RetroTxt Setup - %s",
			i+1, len(keys), key)
		if i == 0 {
			fmt.Println(hr(width))
		}
		fmt.Println(h)
		Update(key, true)
		fmt.Println(hr(width))
	}
	fmt.Println(Info(viper.GetString("style.info")))
}

func enterKey() string {
	if runtime.GOOS == "darwin" {
		return "Press ↩ return to skip the question or ⌃ control-c to quit"
	}
	return "Press ⏎ enter to skip the question or Ctrl-c to quit"
}

func hr(width uint) string {
	return str.Cb(strings.Repeat("-", int(width)))
}

func logo() {
	const clear, reset, n = "\033c", "\033[0m", "text/retrotxt.utf8ans"
	b, err := static.Text.ReadFile(n)
	if err != nil {
		logs.Fatal("unknown pack name", n, ErrLogo)
	}
	// the terminal screen needs to be cleared if the logo is to display correctly
	fmt.Println(clear + string(b) + reset)
}

// watch intercepts Ctrl-C exit key combination.
func watch() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n\nexited setup\n")
		os.Exit(0)
	}()
}
