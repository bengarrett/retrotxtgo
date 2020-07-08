package logs

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
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

// Hint is a generic error type used to apply color to errors and offer a hint
type Hint struct {
	Issue string // Issue is a summary of the problem
	Arg   string // Arg is the argument, flag or item that triggered the error
	Msg   error  // Msg is the actual error generated
	Hint  string // Hint is an optional solution such as a retrotxt command
}

func (h Hint) String() string {
	e := Generic{Issue: h.Issue, Arg: h.Arg, Err: h.Msg}
	return e.String() + fmt.Sprintf("\n         run %s", str.Example("retrotxt "+h.Hint))
}

// CmdErr is a root command error type to handle command arguments and flags
type CmdErr struct {
	Args []string // Command line arguments
	Err  error    // rootCmd.Execute output
}

func (e CmdErr) Error() Generic {
	quote := func(s string) string {
		return fmt.Sprintf("%q", s)
	}
	s := e.Strings()
	l := len(s)
	if l < 3 {
		Log(errors.New("cmderr.args: word count is < 3"))
	} else if len(e.Args) < 1 {
		Log(errors.New("cmderr.err: value is empty"))
	}
	a := fmt.Sprintf("%q", e.Args[0])
	switch strings.Join(s[0:2], " ") {
	case "bad flag":
		return Generic{Issue: "flag syntax",
			Arg: quote(s[l-1]),
			Err: errors.New("flags can only be in -s (short) or --long (long) form")}
	case "flag needs":
		return Generic{Issue: "invalid flag",
			Arg: quote(s[l-1]),
			Err: errors.New("cannot be empty and requires a value")}
	case "invalid argument":
		m := strings.Split(fmt.Sprint(e.Err), ":")
		return Generic{Issue: "flag value",
			Arg: a,
			Err: errors.New(m[0])}
	case "invalid slice":
		return Generic{Issue: "flag value",
			Arg: quote(s[l-1]),
			Err: errors.New("is not a valid choice for --" + s[l-2])}
	case "invalid command":
		return Generic{Issue: "invalid command",
			Arg: quote(s[l-1]),
			Err: errors.New("choose another command from the available commands")}
	case "required flag(s)":
		return Generic{Issue: "a required flag is missing",
			Arg: s[2],
			Err: errors.New("you must include this flag in your command")}
	case "subcommand is":
		fmt.Printf("SUBCMD DEBUG: %+v", e.Err)
		return Generic{} // ignore error
	case "unknown command":
		return Generic{Issue: "invalid command",
			Arg: a,
			Err: errors.New("choose a command from the list available")}
	case "unknown flag:":
		return Generic{Issue: "unknown flag",
			Arg: s[2],
			Err: errors.New("use a flag from the list of flags")}
	case "unknown shorthand":
		return Generic{Issue: "unknown shorthand flag",
			Arg: s[5],
			Err: errors.New("use a flag from the list of flags")}
	}
	//fmt.Printf("DEBUG: %+v\n", e.Err)
	return Generic{Issue: "command", Arg: "execute", Err: e.Err}
}

func (e CmdErr) String() {
	fmt.Println(e.Error())
}

// Strings splits the CmdErr.Err value.
func (e CmdErr) Strings() []string {
	return strings.Split(fmt.Sprintf("%s", e.Err), " ")
}

// Exit prints the CmdErr and causes the program to exit with the given status code.
func (e CmdErr) Exit(code int) {
	fmt.Println(e.Error())
	os.Exit(code)
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

// CheckCmd returns instructions for invalid command arguments.
func CheckCmd(args []string) {
	err := CmdErr{Args: args,
		Err: errors.New("invalid command " + args[0])}
	err.Exit(1)
}

// CheckFlag returns flag options when an invalid choice is used.
func CheckFlag(name, value string, args []string) {
	err := CmdErr{Args: args,
		Err: fmt.Errorf("invalid slice %s %s", name, value)}
	err.String()
	fmt.Printf("         choices: %s\n", strings.Join(args, ", "))
	os.Exit(1)
}

func checkCmd(arg string, args []string) (msg string, code int) {
	code = 10
	msg += fmt.Sprintf("%s invalid argument!%s", str.Alert(),
		color.Bold.Sprintf(" %q", arg))
	if len(args) > 1 {
		msg += fmt.Sprintf(", choices: %s\n%s",
			color.Info.Sprintf("%s", strings.Join(args, ", ")),
			color.Warn.Sprint("please use one of the ")+
				color.Info.Sprint("argument choices")+
				color.Warn.Sprint(" shown above"))
	}
	return msg, code
}
