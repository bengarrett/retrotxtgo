package logs

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Cmd is a command error type to handle command arguments and flags
type Cmd struct {
	Args []string // Command line arguments
	Err  error    // rootCmd.Execute output
}

// Execute is the error handler for command flags and arguments.
func Execute(err error, args ...string) {
	cmd := Cmd{Args: args, Err: err}
	fmt.Println(cmd.error().String())
}

// ArgFatal returns instructions for invalid command arguments.
func ArgFatal(args ...string) {
	cmd := Cmd{Args: args, Err: errors.New("invalid command " + args[0])}
	fmt.Println(cmd.error().String())
	os.Exit(1)
}

// FlagFatal returns flag options when an invalid choice is used.
func FlagFatal(name, value string, args ...string) {
	cmd := Cmd{Args: args, Err: fmt.Errorf("invalid slice %s %s", name, value)}
	fmt.Println(cmd.error().String())
	fmt.Printf("         choices: %s\n", strings.Join(args, ", "))
	os.Exit(1)
}

func (c Cmd) error() Generic {
	var (
		s     = strings.Split(fmt.Sprintf("%s", c.Err), " ")
		l     = len(s)
		quote = func(s string) string {
			return fmt.Sprintf("%q", s)
		}
	)
	const min = 3
	if l < min {
		LogFatal(errors.New("cmd error args: word count is too short, less than 3"))
	} else if len(c.Args) < 1 {
		LogFatal(errors.New("cmd error err: value is empty"))
	}
	a := fmt.Sprintf("%q", c.Args[0])
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
		m := strings.Split(fmt.Sprint(c.Err), ":")
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
		fmt.Printf("SUBCMD DEBUG: %+v", c.Err)
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
	//fmt.Printf("DEBUG: %+v\n", c.Err)
	return Generic{Issue: "command", Arg: "execute", Err: c.Err}
}
