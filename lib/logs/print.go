package logs

import (
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
)

const errorCode = 1

// Hint formats and returns the error with a usage suggestion or hint.
func Hint(s string, err error) string {
	if err == nil {
		return ""
	}
	if s == "" {
		return Printf(err)
	}
	return fmt.Sprintf("%s\n run %s",
		Printf(err), str.Example(fmt.Sprintf("%s %s", meta.Bin, s)))
}

// Fatal prints a formatted error and exits.
func Fatal(err error) {
	if err != nil {
		fmt.Println(Printf(err))
	}
	os.Exit(errorCode)
}

// Printf formats and returns the error.
func Printf(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%s %s", str.Alert(), err)
}

// ProblemCmd returns the command does not exist.
func ProblemCmd(name string, err error) string {
	if name == "" || err == nil {
		return ""
	}
	return fmt.Sprintf("%s the command %s does not exist, %s",
		str.Alert(), name, err)
}

// ProblemCmdFatal prints a problem with the flag and exits.
func ProblemCmdFatal(name, flag string, err error) {
	fmt.Println(ProblemFlag(name, flag, err))
	os.Exit(errorCode)
}

// ProblemFlag prints a problem with the flag.
func ProblemFlag(name, flag string, err error) string {
	if name == "" || err == nil {
		return ""
	}
	alert, toggle := str.Alert(), "--"
	if strings.Contains(flag, "-") {
		toggle = ""
	} else if len(flag) == 1 {
		toggle = "-"
	}
	return fmt.Sprintf("%s with the %s %s%s flag, %s",
		alert, name, toggle, flag, err)
}

// ProblemMark formats the errors and highlights the value.
func ProblemMark(value string, err, errs error) string {
	if value == "" || err == nil || errs == nil {
		return ""
	}
	return fmt.Sprintf("%s %s %q: %s",
		str.Alert(), str.Cf(fmt.Sprintf("%v", err)), value, str.Cf(fmt.Sprintf("%v", errs)))
}

// ProblemMark formats the errors, highlights the value and exits.
func ProblemMarkFatal(value string, err, errs error) {
	fmt.Println(ProblemMark(value, err, errs))
	os.Exit(errorCode)
}

// Problemf formats the errors.
func Problemf(err, errs error) string {
	if err == nil || errs == nil {
		return ""
	}
	e := fmt.Errorf("%s: %w", err, errs)
	return fmt.Sprintf("%s%s", str.Alert(), e)
}

// Problemf formats the errors and exits.
func ProblemFatal(err, errs error) {
	fmt.Println(Problemf(err, errs))
	os.Exit(errorCode)
}
