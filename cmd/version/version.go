// Package version provides the stdout template for the version flag.
package version

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/bengarrett/retrotxtgo/cmd/update"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
	"github.com/mattn/go-isatty"
	"golang.org/x/term"
)

const (
	TabWidth = 8 // Width of tab characters.
)

// Template writes the application version, copyright and build variables.
func Template(wr io.Writer) error {
	if wr == nil {
		wr = io.Discard
	}
	const c = "\u00A9"
	exe, err := Self()
	if err != nil {
		exe = err.Error()
	}
	tag, err := update.Check()
	if err != nil {
		return err
	}
	appDate := ""
	if meta.App.Date != meta.Placeholder {
		appDate = fmt.Sprintf(" (%s)", meta.App.Date)
	}
	w := tabwriter.NewWriter(wr, 0, TabWidth, 0, '\t', 0)
	fmt.Fprintf(w, "%s %s\n", meta.Name, meta.String())
	fmt.Fprintf(w, "%s %s Ben Garrett\n", meta.Copyright, c)
	fmt.Fprintln(w, color.Primary.Sprint(meta.URL))
	fmt.Fprintf(w, "\n%s\t%s %s%s\n", color.Secondary.Sprint("build:"), runtime.Compiler, meta.App.BuiltBy, appDate)
	fmt.Fprintf(w, "%s\t%s/%s\n", color.Secondary.Sprint("platform:"), runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(w, "%s\t%s\n", color.Secondary.Sprint("terminal:"), Terminal())
	fmt.Fprintf(w, "%s\t%s\n", color.Secondary.Sprint("go:"), strings.Replace(runtime.Version(), "go", "v", 1))
	fmt.Fprintf(w, "%s\t%s\n", color.Secondary.Sprint("path:"), exe)
	if tag != "" {
		fmt.Fprintln(w)
		update.Notice(w, meta.App.Version, tag)
		fmt.Fprintln(w)
	}
	return w.Flush()
}

// Self returns the path to the executable (this) program.
func Self() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("self error: %w", err)
	}
	return exe, nil
}

// Terminal attempts to determine the host shell environment.
func Terminal() string {
	const win = "windows"
	unknown := func() string {
		if runtime.GOOS == win {
			return "PowerShell"
		}
		return "unknown"
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return unknown()
	}
	defer restoreTerm(oldState)
	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	// code source: https://gist.github.com/mattn/00cf5b7e38f4cceaf7077f527479870c
	if os.Getenv("WT_SESSION") != "" {
		const s = "Windows Terminal"
		if err != nil {
			return s
		}
		return fmt.Sprintf("%s (%dx%d)", s, w, h)
	}
	if isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		const s = "Cygwin"
		if err != nil {
			return s
		}
		return fmt.Sprintf("%s (%dx%d)", s, w, h)
	}
	_, err = os.Stdout.Write([]byte("\x1b[c"))
	if err != nil {
		return unknown()
	}
	defer os.Stdout.SetReadDeadline(time.Time{})
	const timeout = 10 * time.Millisecond
	time.Sleep(timeout)

	var b [100]byte
	n, err := os.Stdout.Read(b[:])
	if err != nil {
		return unknown()
	}
	if n > 0 {
		return fmt.Sprintf("VT100 compatible (%dx%d)", w, h)
	}
	return unknown()
}

func restoreTerm(oldState *term.State) {
	if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
		log.Fatal(err)
	}
}
