package config

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/static"
	"github.com/spf13/viper"
)

// Setup walks through all the settings and saves them to the configuration file.
func Setup() {
	keys := Keys()
	logo()
	PrintLocation()
	const width = 80
	watch()
	for i, key := range keys {
		if i == 0 {
			fmt.Printf("\n\n  %s\n\n", str.Cinf(str.Center(int(width), enterKey())))
		}
		h := fmt.Sprintf("  %d/%d. %s Setup - %s",
			i+1, len(keys), meta.Name, key)
		if i == 0 {
			fmt.Println(str.HR(width))
			fmt.Println("")
		}
		fmt.Println(h)
		Update(key, true)
		fmt.Println(str.HRPadded(width))
	}
	fmt.Println(Info(viper.GetString("style.info")))
}

// EnterKey returns the appropriate Setup instructions based on the user's platform.
func enterKey() string {
	if runtime.GOOS == "darwin" {
		return "Press ↩ return to skip the question or ⌃ control-c to quit"
	}
	return "Press ⏎ return to skip the question or Ctrl-c to quit"
}

// Logo prints the ANSI logo.
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n\n%s Quit setup\n", str.Info())
		os.Exit(0)
	}()
}
