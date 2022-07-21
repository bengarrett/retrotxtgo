// Package prompt inputs for user interactions.
package prompt

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/prompt/internal/key"
	"github.com/bengarrett/retrotxtgo/lib/prompt/internal/port"
	"github.com/bengarrett/retrotxtgo/lib/prompt/internal/read"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/gookit/color"
)

var ErrSkip = errors.New("skipped, no change")

const (
	PortMin uint = port.Min
	PortMax uint = port.Max
)

// CtrlC intercepts any Ctrl-C keyboard input and exits to the shell.
func CtrlC(w io.Writer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Fprintf(w, "\n  %s\n", ErrSkip)
		os.Exit(0)
	}()
}

// IndexStrings asks for a numeric index position and returns a single choice from a string of keys.
func IndexStrings(w io.Writer, options *[]string, setup bool) string {
	if options == nil {
		return ""
	}
	var k key.Keys = *options
	CtrlC(w)
	return k.Prompt(w, os.Stdin, setup)
}

// SkipSet returns a skipped setting string for the setup walk through.
func SkipSet(setup bool) string {
	if !setup {
		return ""
	}
	return str.ColSuc("  skipped setting")
}

// Port asks for and returns a HTTP port value.
func Port(w io.Writer, validate, setup bool) uint {
	return port.Port(w, os.Stdin, validate, setup)
}

// PortValid checks if the network port is within range to serve HTTP.
func PortValid(p uint) bool {
	return port.Valid(p)
}

// ShortStrings asks for and returns a single choice from a string of keys.
// Either the first letter or the full name of the key are accepted.
func ShortStrings(w io.Writer, options *[]string) string {
	if options == nil {
		return ""
	}
	var k key.Keys = *options
	CtrlC(w)
	return k.ShortPrompt(w, os.Stdin)
}

// String asks for and returns a multi-word string.
// Inputting ⏎ the Enter/Return key exits the program,
// or returns an empty value when in SetupMode.
func String(w io.Writer) string {
	CtrlC(w)
	return read.Parse(os.Stdin)
}

// Strings asks for and returns a single choice from a string of keys.
func Strings(w io.Writer, options *[]string, setup bool) string {
	if options == nil {
		return ""
	}
	var k key.Keys = *options
	return k.Prompt(w, os.Stdin, setup)
}

// YesNo asks for a yes or no input.
func YesNo(w io.Writer, ask string, yesDefault bool) bool {
	y, n := "Y", "n"
	if !yesDefault {
		y, n = "y", "N"
	}
	if !strings.HasSuffix(ask, "?") && !strings.HasSuffix(ask, ")") {
		ask = fmt.Sprintf("%s?", ask)
	}
	yes, no := y, n
	if color.Enable {
		yes, _ = str.UnderlineChar(y)
		no, _ = str.UnderlineChar(n)
	}
	fmt.Fprintf(w, "%s [%ses/%so] ", ask, yes, no)
	CtrlC(w)
	input, err := read.Read(os.Stdin)
	if err != nil {
		logs.FatalSave(err)
	}
	return read.ParseYN(input, yesDefault)
}
