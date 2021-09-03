package logs

import (
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
)

const OSErrCode = 1

// Hint formats and returns the error with a usage suggestion or hint.
func Hint(s string, err error) string {
	if err == nil {
		return ""
	}
	if s == "" {
		return Sprint(err)
	}
	return fmt.Sprintf("%s\n run %s",
		Sprint(err), str.Example(fmt.Sprintf("%s %s", meta.Bin, s)))
}

// Fatal prints a formatted error and exits.
func Fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, Sprint(err))
	}
	os.Exit(OSErrCode)
}

// FatalFlag prints a problem with the flag and exits.
func FatalFlag(cmd, flag string, err error) {
	fmt.Fprintln(os.Stderr, PrintfFlag(cmd, flag, err))
	os.Exit(OSErrCode)
}

// FatalMark formats the errors, highlights the value and exits.
func FatalMark(mark string, err, wrap error) {
	fmt.Fprintln(os.Stderr, PrintfMark(mark, err, wrap))
	os.Exit(OSErrCode)
}

// FatalWrap formats the errors and exits.
func FatalWrap(err, wrap error) {
	fmt.Fprintln(os.Stderr, PrintfWrap(err, wrap))
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
	return fmt.Sprintf("%s%s.", str.Alert(), strings.Join(seps, ".\n"))
}

// PrintfCmd returns the command does not exist.
func PrintfCmd(cmd string, err error) string {
	if cmd == "" || err == nil {
		return ""
	}
	return fmt.Sprintf("%s the command %s does not exist, %s",
		str.Alert(), cmd, err)
}

// PrintfFlag prints a problem with the flag.
func PrintfFlag(cmd, flag string, err error) string {
	if cmd == "" || err == nil {
		return ""
	}
	alert, toggle := str.Alert(), "--"
	if strings.Contains(flag, "-") {
		toggle = ""
	} else if len(flag) == 1 {
		toggle = "-"
	}
	return fmt.Sprintf("%s with the %s %s%s flag, %s",
		alert, cmd, toggle, flag, err)
}

// PrintfMark formats the errors and highlights the value.
func PrintfMark(mark string, err, wrap error) string {
	if mark == "" || err == nil || wrap == nil {
		return ""
	}
	return fmt.Sprintf("%s %s %q: %s",
		str.Alert(), str.ColFuz(fmt.Sprintf("%v", err)), mark, str.ColFuz(fmt.Sprintf("%v", wrap)))
}

// PrintfWrap formats the errors.
func PrintfWrap(err, wrap error) string {
	if err == nil || wrap == nil {
		return ""
	}
	return fmt.Sprintf("%s%s",
		str.Alert(), fmt.Errorf("%s: %w", err, wrap))
}
