package config

import (
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

// Setup walks through all the settings and saves them to the configuration file.
func Setup() {
	keys := Keys()
	logo()
	PrintLocation()
	var width uint = 80
	watch()
	for i, key := range keys {
		if i == 0 {
			fmt.Printf("\n\n  %s\n\n", str.Cinf(str.Center(int(width), enterKey())))
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

// EnterKey returns the appropriate Setup instructions based on the user's platform.
func enterKey() string {
	if runtime.GOOS == "darwin" {
		return "Press ↩ return to skip the question or ⌃ control-c to quit"
	}
	return "Press ⏎ enter to skip the question or Ctrl-c to quit"
}

// HR returns a horizontal rule.
func hr(width uint) string {
	return str.Cb(strings.Repeat("-", int(width)))
}

// Logo prints the RetroTxt ANSI logo.
func logo() {
	const clear, reset, n = "\033c", "\033[0m", "text/retrotxt.utf8ans"
	b, err := static.Text.ReadFile(n)
	if err != nil {
		logs.ProblemMarkFatal(n, logs.ErrSampFile, ErrLogo)
	}
	// the terminal screen needs to be cleared if the logo is to display correctly
	fmt.Println(clear + string(b) + reset)
}

// Watch intercepts Ctrl-C key combinations to exit out of the Setup.
func watch() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n\nexited setup\n")
		os.Exit(0)
	}()
}
