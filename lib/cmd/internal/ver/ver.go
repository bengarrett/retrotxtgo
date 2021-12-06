package ver

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/release"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
	"github.com/mattn/go-isatty"
	"golang.org/x/term"
)

// printVer
func Print() string {
	const tabWidth, copyright, years = 8, "\u00A9", "2020-21"
	exe, err := self()
	if err != nil {
		exe = err.Error()
	}
	newVer, v := release.Check()
	appDate := ""
	if meta.App.Date != meta.Placeholder {
		appDate = fmt.Sprintf(" (%s)", meta.App.Date)
	}
	var b bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&b, 0, tabWidth, 0, '\t', 0)
	fmt.Fprintf(w, "%s %s\n", meta.Name, meta.Print())
	fmt.Fprintf(w, "%s %s Ben Garrett\n", copyright, years)
	fmt.Fprintln(w, color.Primary.Sprint(meta.URL))
	fmt.Fprintf(w, "\n%s\t%s %s%s\n", color.Secondary.Sprint("build:"), runtime.Compiler, meta.App.BuiltBy, appDate)
	fmt.Fprintf(w, "%s\t%s/%s\n", color.Secondary.Sprint("platform:"), runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(w, "%s\t%s\n", color.Secondary.Sprint("terminal:"), terminal())
	fmt.Fprintf(w, "%s\t%s\n", color.Secondary.Sprint("go:"), strings.Replace(runtime.Version(), "go", "v", 1))
	fmt.Fprintf(w, "%s\t%s\n", color.Secondary.Sprint("path:"), exe)
	if newVer {
		fmt.Fprintf(w, "\n%s\n", release.Print(meta.App.Version, v))
	}
	w.Flush()
	return b.String()
}

// Self returns the path to this dupers executable file.
func self() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("self error: %w", err)
	}
	return exe, nil
}

func terminal() string {
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
	defer term.Restore(int(os.Stdin.Fd()), oldState)
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
		s := "VT100 compatible"
		if err != nil {
			return s
		}
		return fmt.Sprintf("%s (%dx%d)", s, w, h)
	}
	return unknown()
}
