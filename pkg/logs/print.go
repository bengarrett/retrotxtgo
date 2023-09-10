package logs

import (
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/pkg/term"
)

const OSErrCode = 1

// Hint returns a formatted error with a usage suggestion or hint.
func Hint(s string, err error) string {
	if err == nil {
		return ""
	}
	if s == "" {
		return Sprint(err)
	}
	return fmt.Sprintf("%s\n run %s",
		Sprint(err), term.Example(fmt.Sprintf("%s %s", meta.Bin, s)))
}

// FatalFlag prints a problem with the flag and exits.
func FatalFlag(cmd, flag string, err error) {
	fmt.Fprintln(os.Stderr, SprintFlag(cmd, flag, err))
	os.Exit(OSErrCode)
}

// FatalMark formats the errors, highlights the value and exits.
func FatalMark(mark string, err, wrap error) {
	fmt.Fprintln(os.Stderr, SprintMark(mark, err, wrap))
	os.Exit(OSErrCode)
}

// FatalWrap formats the errors and exits.
func FatalWrap(err, wrap error) {
	fmt.Fprintln(os.Stderr, SprintWrap(err, wrap))
	os.Exit(OSErrCode)
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
func SprintCmd(cmd string, err error) string {
	if cmd == "" || err == nil {
		return ""
	}
	return fmt.Sprintf("%s the command %s does not exist, %s",
		term.Alert(), cmd, err)
}

// SprintFlag returns a problem with the flag.
func SprintFlag(cmd, flag string, err error) string {
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

// SprintMark formats and returns the errors and highlights the marked string.
func SprintMark(mark string, err, wrap error) string {
	if mark == "" || err == nil || wrap == nil {
		return ""
	}
	return fmt.Sprintf("%s %s %q: %s",
		term.Alert(), term.Fuzzy(fmt.Sprintf("%v", err)), mark, term.Fuzzy(fmt.Sprintf("%v", wrap)))
}

// SprintWrap returns the formatted errors.
func SprintWrap(err, wrap error) string {
	if err == nil || wrap == nil {
		return ""
	}
	return fmt.Sprintf("%s%s",
		term.Alert(), fmt.Errorf("%w: %w", err, wrap))
}
