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
// Use start to begin the walk through at the question number, or leave it at 0.
func Setup(start int) {
	const width = 80
	keys := Keys()
	fmt.Printf("%s\n%s\n%s\n%s\n\n",
		logo(),
		fmt.Sprintf("Walk through all of the %s settings.", meta.Name),
		Location(),
		enterKey())
	fmt.Println(str.HRPadded(width))
	watch()
	for i, key := range keys {
		if start > i+1 {
			continue
		}
		h := fmt.Sprintf("  %d/%d. %s Setup - %s",
			i+1, len(keys), meta.Name, key)
		fmt.Println(h)
		Update(key, true)
		fmt.Println(str.HRPadded(width))
	}
	fmt.Println(Info(viper.GetString("style.info")))
}

// EnterKey returns the appropriate Setup instructions based on the user's platform.
func enterKey() string {
	if runtime.GOOS == "darwin" {
		return "At any time press ↩ return to skip the question or ⌃ control-c to quit setup."
	}
	return "At any time press ⏎ return to skip the question or Ctrl-c to quit setup."
}

// Logo prints the ANSI logo.
func logo() string {
	const clear, reset, n = "\033c", "\033[0m", "text/retrotxt.utf8ans"
	b, err := static.Text.ReadFile(n)
	if err != nil {
		logs.ProblemMarkFatal(n, logs.ErrSampFile, ErrLogo)
	}
	// the terminal screen needs to be cleared if the logo is to display correctly
	return fmt.Sprint(clear + string(b) + reset)
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
