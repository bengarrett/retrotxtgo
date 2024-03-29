// Package logs handles the formatting and returning of errors.
package logs

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/term"
)

const OSErr = 1 // OSErr is the operating system exit code for a program error.

// Fatal prints the error to stderr and exits.
func Fatal(err error) {
	if err == nil {
		return
	}
	const panicing = false
	switch panicing {
	case true:
		log.Printf("error type: %T\tmsg: %v\n", err, err)
		log.Panic(err)
	default:
		fmt.Fprintln(os.Stderr, Sprint(err))
		os.Exit(OSErr)
	}
}

// FatalS highlights the string and errors then exits.
func FatalS(err, wrap error, s string) {
	fmt.Fprintln(os.Stderr, SprintS(err, wrap, s))
	os.Exit(OSErr)
}

// Hint returns a formatted error with a usage suggestion or hint.
// If s is empty then just the error is formatted.
func Hint(err error, s string) string {
	if err == nil {
		return ""
	}
	if s == "" {
		return Sprint(err)
	}
	return fmt.Sprintf("%s\n run %s",
		Sprint(err), term.Example(fmt.Sprintf("%s %s", meta.Bin, s)))
}

// Sprint formats and returns the error.
func Sprint(err error) string {
	if err == nil {
		return ""
	}
	elms, seps := strings.Split(err.Error(), ";"), []string{}
	for _, elm := range elms {
		if elm == "" || elm == "<nil>" {
			continue
		}
		seps = append(seps, elm)
	}
	return fmt.Sprintf("%s%s.", term.Alert(), strings.Join(seps, ".\n"))
}

// SprintCmd returns the command does not exist.
func SprintCmd(err error, cmd string) string {
	if cmd == "" || err == nil {
		return ""
	}
	return fmt.Sprintf("%s the command %s does not exist, %s",
		term.Alert(), cmd, err)
}

// SprintFlag returns a problem with the flag.
func SprintFlag(err error, cmd, flag string) string {
	if cmd == "" || err == nil {
		return ""
	}
	alert, toggle := term.Alert(), "--"
	if strings.Contains(flag, "-") {
		toggle = ""
	} else if len(flag) == 1 {
		toggle = "-"
	}
	return fmt.Sprintf("%s with the %s %s%s flag, %s",
		alert, cmd, toggle, flag, err)
}

// SprintS highlights the string and errors then exits.
func SprintS(err, wrap error, s string) string {
	if s == "" || err == nil || wrap == nil {
		return ""
	}
	return fmt.Sprintf("%s %s %q: %s",
		term.Alert(), term.Fuzzy(fmt.Sprintf("%v", err)), s,
		term.Fuzzy(fmt.Sprintf("%v", wrap)))
}
