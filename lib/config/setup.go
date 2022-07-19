package config

import (
	"fmt"
	"io"
	"log"
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
func CtrlC(w io.Writer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Fprintf(w, "  exited from the %s setup", meta.Name)
		os.Exit(0)
	}()
}

// Setup walks through all the settings and saves them to the configuration file.
// Start the walk through at the provided question number or leave it at 0.
func Setup(w io.Writer, start int) {
	const width = 80
	keys := KeySort()
	fmt.Fprintln(w, logo())
	fmt.Fprintf(w, "Walk through all of the %s settings.\n", meta.Name)
	fmt.Fprintln(w, Location())
	fmt.Fprintln(w, enterKey())
	fmt.Fprintln(w, "\n"+str.HRPad(width))
	CtrlC(w)
	for i, key := range keys {
		if start > i+1 {
			continue
		}
		h := fmt.Sprintf("  %d/%d. %s Setup - %s",
			i+1, len(keys), meta.Name, key)
		fmt.Fprintln(w, h)
		if err := Update(w, key, true); err != nil {
			log.Print(err)
		}
		fmt.Fprintln(w, str.HRPad(width))
	}
	fmt.Fprintf(w, "The %s setup and configuration is complete.\n", meta.Name)
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
		logs.FatalMark(n, logs.ErrSampleName, ErrLogo)
	}
	// the terminal screen needs to be cleared if the logo is to display correctly
	return fmt.Sprintf("%s%s%s",
		clearScreen, string(b), resetScreen)
}
