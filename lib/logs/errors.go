package logs

import (
	"fmt"
	"os"

	"retrotxt.com/retrotxt/lib/str"
)

// TODO: replace much of the logs package with a pre-existing sys or structured loggers.
// zerolog: https://github.com/rs/zerolog
// apex: https://github.com/apex/log
// zap: https://github.com/uber-go/zap

// Generic is the standard error type used to apply color to errors
type Generic struct {
	Issue string // Issue is a summary of the problem
	Arg   string // Arg is the argument, flag or item that triggered the error
	Err   error  // Err is the actual error generated
}

func (g Generic) String() string {
	if g.Err == nil {
		return ""
	}
	return fmt.Sprintf("%s %s %s", str.Alert(),
		str.Ci(fmt.Sprintf("%s %s,", g.Issue, g.Arg)), // issue & argument
		str.Cf(fmt.Sprintf("%v", g.Err)))              // error message
}

// Hint is a standard error type that also offers a user hint
type Hint struct {
	Gen  Generic
	Hint string // Hint is an optional solution such as a retrotxt command
}

func (h Hint) String() string {
	gen := h.Gen
	return gen.String() + fmt.Sprintf("\n         run %s", str.Example("retrotxt "+h.Hint))
}

// Fatal prints a generic error and exits.
func (g Generic) Fatal() {
	fmt.Println(g.String())
	os.Exit(1)
}

// Println prints a generic error type.
func Println(issue, arg string, err error) {
	var g = Generic{
		Issue: issue,
		Arg:   arg,
		Err:   err,
	}
	fmt.Println(g.String())
}

// Fatal prints a generic error and exits.
func Fatal(issue, arg string, msg error) {
	Println(issue, arg, msg)
	os.Exit(1)
}

// IssueErr is a generic problem structure.
type IssueErr struct {
	Issue string
	Err   error
}

func (i IssueErr) String() string {
	if i.Err == nil {
		return ""
	}
	if i.Issue == "" {
		return fmt.Sprintf("%s", i.Err)
	}
	return fmt.Sprintf("%s %s\n         %s", str.Alert(), str.Ci(i.Issue), str.Cf(fmt.Sprint(i.Err)))
}

// Exit prints the IssueErr and causes the program to exit with the given status code.
func (i IssueErr) Exit(code int) {
	fmt.Println(i.String())
	os.Exit(code)
}
