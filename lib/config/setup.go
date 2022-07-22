package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/static"
)

var ErrStart = errors.New("setup start argument is out of range")

const (
	ClearScreen = "\033c"
	ResetScreen = "\033[0m"

	logoname = "ansi/retrotxt.utf8ans"
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
// The start argument can being the walkthrough at the provided question number or it can be left at 0.
func Setup(w io.Writer, start int) error {
	if w == nil {
		return ErrWriter
	}
	const width = 80
	keys := SortKeys()
	if start < 0 || start > len(keys)-1 {
		return fmt.Errorf("%w: %d (permitted values 0-%d)", ErrStart, start, len(keys)-1)
	}
	if err := Logo(w); err != nil {
		return err
	}
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
			return err
		}
		fmt.Fprintln(w, str.HRPad(width))
	}
	fmt.Fprintf(w, "The %s setup and configuration is complete.\n", meta.Name)
	return nil
}

// enterKey returns the appropriate Setup instructions based on the host platform.
func enterKey() string {
	if runtime.GOOS == "darwin" {
		return "To quit setup, press ↩ return to skip the question or ⌃ control-C."
	}
	return "To quit setup, press ⏎ return to skip the question or Ctrl-C."
}

// Logo writes a custom 256-color, ANSI logo for RetroTxt.
func Logo(w io.Writer) error {
	if w == nil {
		return ErrWriter
	}
	b, err := static.ANSI.ReadFile(logoname)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrLogo, logoname)
	}
	// the terminal screen needs to be cleared if the logo is to display correctly
	fmt.Fprintf(w, "%s%s%s\n", ClearScreen, string(b), ResetScreen)
	return nil
}
