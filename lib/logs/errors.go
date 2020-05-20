package logs

import (
	"errors"
	"fmt"
	"strings"
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
	return fmt.Sprintf("%s %s %s", Alert(),
		Ci(fmt.Sprintf("%s %s,", e.Issue, e.Arg)), // issue & argument
		Cf(fmt.Sprintf("%v", e.Msg)))              // error message
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
	s := strings.Split(e.String(), " ")
	l := len(s)
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
	case "required flag(s)":
		return Err{Issue: "a required flag missing",
			Arg: a,
			Msg: e.Err} // TODO: test
	case "subcommand is":
		fmt.Printf("SUBCMD DEBUG: %+v", e.Err)
		return Err{} // ignore error
	case "unknown command":
		return Err{Issue: "invalid command",
			Arg: a,
			Msg: errors.New("choose a command from the list available")}
	case "unknown flag:", "unknown shorthand":
		return Err{Issue: "unknown flag",
			Arg: a,
			Msg: errors.New("use a flag from the list of flags")}
	}

	fmt.Printf("DEBUG: %+v\n", e.Err)
	return Err{Issue: "command", Arg: "execute", Msg: e.Err}
}

func (e CmdErr) String() string {
	return fmt.Sprintf("%s", e.Err)
}

// ConfigErr ...
type ConfigErr struct {
	FileUsed string
	Err      error
}

func (e ConfigErr) String() string {
	return (Err{
		Issue: "config file",
		Arg:   e.FileUsed,
		Msg:   e.Err}).String()
}
