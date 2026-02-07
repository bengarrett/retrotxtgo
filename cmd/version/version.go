// Package version provides the stdout template for the version flag.
package version

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bengarrett/retrotxtgo/cmd/update"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/color"
	"github.com/mattn/go-isatty"
	"golang.org/x/term"
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

	// Define lipgloss styles
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	treeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244"))

	// Build the version information
	content := fmt.Sprintf("%s %s\n", meta.Name, meta.String())
	content += fmt.Sprintf("%s %s Ben Garrett\n", meta.Copyright, c)
	content += color.Primary.Sprint(meta.URL) + "\n\n"

	// Create a tree structure for system info
	content += treeStyle.Render("┌ build: ") + fmt.Sprintf("%s %s%s\n", runtime.Compiler, meta.App.BuiltBy, appDate)
	content += treeStyle.Render("├ platform: ") + fmt.Sprintf("%s/%s\n", runtime.GOOS, runtime.GOARCH)
	content += treeStyle.Render("├ terminal: ") + Terminal() + "\n"
	content += treeStyle.Render("├ go: ") + strings.Replace(runtime.Version(), "go", "v", 1) + "\n"
	content += treeStyle.Render("└ path: ") + exe + "\n"

	if tag != "" {
		content += "\n"
		content += update.NoticeString(meta.App.Version, tag)
		content += "\n"
	}

	// Apply border styling
	styled := borderStyle.Render(content)

	_, err = fmt.Fprint(wr, styled)
	return err
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
