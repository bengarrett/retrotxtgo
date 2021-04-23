package logs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"retrotxt.com/retrotxt/lib/str"
)

// Fatal prints a generic error and exits.
func Fatal(issue, arg string, msg error) {
	Println(issue, arg, msg)
	os.Exit(1)
}

// Println prints a generic error.
func Println(issue, arg string, err error) {
	var g = Argument{
		Issue: issue,
		Arg:   arg,
		Err:   err,
	}
	fmt.Println(g.String())
}

// Generic is the standard error type used to apply color to errors.
type Argument struct {
	Issue string // Issue is a summary of the problem
	Arg   string // Arg is the argument, flag or item that triggered the error // Value
	Err   error  // Err is required and is the actual error generated
}

func (g Argument) String() string {
	if g.Err == nil {
		return ""
	}
	a, c := str.Alert(), str.Cf(fmt.Sprintf("%v", g.Err))

	if s := g.unWrap(); s != "" {
		return fmt.Sprintf("%s %s", a, s)
	}

	if g.Issue == "" && g.Arg == "" {
		return fmt.Sprintf("%s %s", a, c) // alert and err
	}
	var b string
	if g.Arg == "" {
		b = str.Ci(fmt.Sprintf("%s,", g.Issue)) // alert, issue and err
	} else {
		b = str.Ci(fmt.Sprintf("%s %s,", g.Issue, g.Arg)) // alert, issue, arg, err
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
		return fmt.Sprintf("cannot open file %q, %s", g.Arg, str.Cf("is there a typo?"))
	}

	fmt.Printf("\n%T %+v\n", uw, uw)
	return ""
}

// Fatal prints a generic error and exits.
func (g Argument) Fatal() {
	fmt.Println(g.String())
	os.Exit(1)
}

// Hint is a standard error type that also offers the user a command hint.
type Hint struct {
	Error Argument
	Hint  string // Hint is an optional solution such as a retrotxt command
}

func (h Hint) String() string {
	err := h.Error
	return err.String() + fmt.Sprintf("\n         run %s", str.Example("retrotxt "+h.Hint))
}
