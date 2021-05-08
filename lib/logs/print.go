package logs

import (
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/str"
)

// Errorf formats and returns the error.
func Errorf(err error) string {
	return fmt.Sprintf("%s %s", str.Alert(), err)
}

// Fatal prints a formatted error and exits.
func Fatal(err error) {
	fmt.Println(Errorf(err))
	os.Exit(1)
}

// Hint formats and returns the error with a usage suggestion or hint.
func Hint(s string, err error) string {
	if err == nil {
		return ""
	}
	if s == "" {
		return Errorf(err)
	}
	return fmt.Sprintf("%s\n         run %s", Errorf(err), str.Example("retrotxt "+s))
}

// ProblemCmd returns the command does not exist.
func ProblemCmd(name string, err error) string {
	alert := str.Alert()
	return fmt.Sprintf("%s the command %s does not exist, %s", alert, name, err)
}

// ProblemCmdFatal prints a problem with the flag and exits.
func ProblemCmdFatal(name, flag string, err error) {
	fmt.Println(ProblemFlag(name, flag, err))
	os.Exit(1)
}

// ProblemFlag prints a problem with the flag.
func ProblemFlag(name, flag string, err error) string {
	alert, toggle := str.Alert(), "--"
	if strings.Contains(flag, "-") {
		toggle = ""
	} else if len(flag) == 1 {
		toggle = "-"
	}
	return fmt.Sprintf("%s with the %s %s%s flag, %s", alert, name, toggle, flag, err)
}

// ProblemMark formats the errors and highlights the value.
func ProblemMark(value string, err, errs error) string {
	v := value
	a := str.Alert()
	n := str.Cf(fmt.Sprintf("%v", err))
	e := str.Cf(fmt.Sprintf("%v", errs))
	return fmt.Sprintf("%s %s %q: %s", a, n, v, e)
}

// ProblemMark formats the errors, highlights the value and exits.
func ProblemMarkFatal(value string, err, errs error) {
	fmt.Println(ProblemMark(value, err, errs))
	os.Exit(1)
}

// Problemf formats the errors.
func Problemf(err, errs error) {
	e := fmt.Errorf("%s: %w", err, errs)
	fmt.Printf("%s%s\n", str.Alert(), e)
}

// Problemf formats the errors and exits.
func ProblemFatal(err, errs error) {
	Problemf(err, errs)
	os.Exit(1)
}
