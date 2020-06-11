package logs

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/gookit/color"
)

// Err is a generic error type used to apply color to errors
type Err struct {
	Issue string // Issue is a summary of the problem
	Arg   string // Arg is the argument, flag or item that triggered the error
	Msg   error  // Msg is the actual error generated
}

func (e Err) String() string {
	if e.Msg == nil {
		return ""
	}
	return fmt.Sprintf("%s %s %s", str.Alert(),
		str.Ci(fmt.Sprintf("%s %s,", e.Issue, e.Arg)), // issue & argument
		str.Cf(fmt.Sprintf("%v", e.Msg)))              // error message
}

// Hint is a generic error type used to apply color to errors and offer a hint
type Hint struct {
	Issue string // Issue is a summary of the problem
	Arg   string // Arg is the argument, flag or item that triggered the error
	Msg   error  // Msg is the actual error generated
	Hint  string // Hint is an optional solution such as a retrotxt command
}

func (h Hint) String() string {
	e := Err{Issue: h.Issue, Arg: h.Arg, Msg: h.Msg}
	return e.String() + fmt.Sprintf("\n         run %s", str.Example("retrotxt "+h.Hint))
}

// CmdErr is a root command error type to handle command arguments and flags
type CmdErr struct {
	Args []string // Command line arguments
	Err  error    // rootCmd.Execute output
}

func (e CmdErr) Error() Err {
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
		return Err{Issue: "flag syntax",
			Arg: quote(s[l-1]),
			Msg: errors.New("flags can only be in -s (short) or --long (long) form")}
	case "flag needs":
		return Err{Issue: "invalid flag",
			Arg: quote(s[l-1]),
			Msg: errors.New("cannot be empty and requires a value")}
	case "invalid argument":
		m := strings.Split(fmt.Sprint(e.Err), ":")
		return Err{Issue: "flag value",
			Arg: a,
			Msg: errors.New(m[0])}
	case "invalid slice":
		return Err{Issue: "flag value",
			Arg: quote(s[l-1]),
			Msg: errors.New("is not a valid choice for --" + s[l-2])}
	case "invalid command":
		return Err{Issue: "invalid command",
			Arg: quote(s[l-1]),
			Msg: errors.New("choose another command from the available commands")}
	case "required flag(s)":
		return Err{Issue: "a required flag is missing",
			Arg: s[2],
			Msg: errors.New("you must include this flag in your command")}
	case "subcommand is":
		fmt.Printf("SUBCMD DEBUG: %+v", e.Err)
		return Err{} // ignore error
	case "unknown command":
		return Err{Issue: "invalid command",
			Arg: a,
			Msg: errors.New("choose a command from the list available")}
	case "unknown flag:":
		return Err{Issue: "unknown flag",
			Arg: s[2],
			Msg: errors.New("use a flag from the list of flags")}
	case "unknown shorthand":
		return Err{Issue: "unknown shorthand flag",
			Arg: s[5],
			Msg: errors.New("use a flag from the list of flags")}
	}
	//fmt.Printf("DEBUG: %+v\n", e.Err)
	return Err{Issue: "command", Arg: "execute", Msg: e.Err}
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
