package logs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"retrotxt.com/retrotxt/lib/str"
)

// Argument is used to highlight user supplied argument values in error replies.
type Argument struct {
	Issue string // Issue is a human readable summary of the problem
	Value string // Value of the the argument, flag or item that triggered the error
	Err   error  // Err is required error
}

func CmdProblem(name string, err error) string {
	alert := str.Alert()
	return fmt.Sprintf("%s the command %s does not exist, %s", alert, name, err)
}

func CmdProblemFatal(name, flag string, err error) {
	fmt.Println(FlagProblem(name, flag, err))
	os.Exit(1)
}

// rename to Fatal
func ErrorFatal(err error) {
	fmt.Println(Errorf(err))
	os.Exit(1)
}

func Errorf(err error) string {
	return fmt.Sprintf("%s%s", str.Alert(), err) // TODO: change to string return?
}

func FlagProblem(name, flag string, err error) string {
	alert, toggle := str.Alert(), "--"
	fmt.Println("FLAG:", flag)
	if strings.Contains(flag, "-") {
		toggle = ""
	} else if len(flag) == 1 {
		toggle = "-"
	}
	return fmt.Sprintf("%s with the %s %s%s flag, %s", alert, name, toggle, flag, err)
}

func MarkProblem(value string, new, err error) string {
	v := value
	a := str.Alert()
	n := str.Cf(fmt.Sprintf("%v", new))
	e := str.Cf(fmt.Sprintf("%v", err))

	return fmt.Sprintf("%s %s %q: %s", a, n, v, e)
}

func MarkProblemFatal(value string, new, err error) {
	fmt.Println(MarkProblem(value, new, err))
	os.Exit(1)
}

// ok
func Problemln(new, err error) {
	e := fmt.Errorf("%s: %w", new, err)
	fmt.Printf("%s%s\n", str.Alert(), e)
}

// ok
func ProblemFatal(new, err error) {
	Problemln(new, err)
	os.Exit(1)
}

func SubCmdProblem(name string, err error) string {
	alert := str.Alert()
	return fmt.Sprintf("%s the subcommand %s does not exist, %s", alert, name, err)
}

func (g Argument) XString() string {
	if g.Err == nil {
		return ""
	}
	a, c := str.Alert(), str.Cf(fmt.Sprintf("%v", g.Err))

	if s := g.unWrap(); s != "" {
		return fmt.Sprintf("%s %s", a, s)
	}

	if g.Issue == "" && g.Value == "" {
		return fmt.Sprintf("%s %s", a, c) // alert and err
	}
	var b string
	if g.Value == "" {
		b = str.Ci(fmt.Sprintf("%s,", g.Issue)) // alert, issue and err
	} else {
		b = str.Ci(fmt.Sprintf("%s %s,", g.Issue, g.Value)) // alert, issue, arg, err
	}
	return fmt.Sprintf("%s %s %s", a, b, c)
}

func (g Argument) unWrap() string {
	var fp *fs.PathError
	uw := errors.Unwrap(g.Err)
	if uw == nil {
		return ""
	}
	if errors.As(uw, &fp) {
		return fmt.Sprintf("cannot open file %q, %s", g.Value, str.Cf("is there a typo?"))
	}

	fmt.Printf("\n%T %+v\n", uw, uw)
	return ""
}

func Hint(s string, err error) string {
	return fmt.Sprintf("%s\n         run %s", Errorf(err), str.Example("retrotxt "+s))
}
