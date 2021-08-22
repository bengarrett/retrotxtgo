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
)

// CtrlC intercepts Ctrl-C key combinations to exit out of the Setup.
func CtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("  exited from the %s setup", meta.Name)
		os.Exit(0)
	}()
}

// Setup walks through all the settings and saves them to the configuration file.
// Start the walk through at the provided question number or leave it at 0.
func Setup(start int) {
	const width = 80
	keys := KeySort()
	fmt.Printf("%s\n%s\n%s\n%s\n\n",
		logo(),
		fmt.Sprintf("Walk through all of the %s settings.", meta.Name),
		Location(),
		enterKey())
	fmt.Println(str.HRPadded(width))
	CtrlC()
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
	fmt.Printf("The %s setup and configuration is complete.\n", meta.Name)
}

// enterKey returns the appropriate Setup instructions based on the user's platform.
func enterKey() string {
	if runtime.GOOS == "darwin" {
		return "To quit setup, press ↩ return to skip the question or ⌃ control-C."
	}
	return "To quit setup, press ⏎ return to skip the question or Ctrl-C."
}

// logo returns the 256-color, ANSI logo.
func logo() string {
	const clearScreen, resetScreen, n = "\033c", "\033[0m", "text/retrotxt.utf8ans"
	b, err := static.Text.ReadFile(n)
	if err != nil {
		logs.ProblemMarkFatal(n, logs.ErrSampFile, ErrLogo)
	}
	// the terminal screen needs to be cleared if the logo is to display correctly
	return fmt.Sprintf("%s%s%s",
		clearScreen, string(b), resetScreen)
}
