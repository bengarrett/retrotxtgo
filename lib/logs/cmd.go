package logs

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrCmd     = errors.New("choose a command from the list available")
	ErrNewCmd  = errors.New("choose another command from the available commands")
	ErrNoCmd   = errors.New("invalid command")
	ErrEmpty   = errors.New("value is empty")
	ErrFlag    = errors.New("use a flag from the list of flags")
	ErrSyntax  = errors.New("flags can only be in -s (short) or --long (long) form")
	ErrNoFlag  = errors.New("cannot be empty and requires a value")
	ErrReqFlag = errors.New("you must include this flag in your command")
	ErrSlice   = errors.New("invalid slice")
	ErrShort   = errors.New("word count is too short, less than 3")
	ErrVal     = errors.New("value is not a valid choice")
)

// Cmd is a command error type to handle command arguments and flags.
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
	cmd := Cmd{Args: args, Err: fmt.Errorf("%q: %w", args[0], ErrNoCmd)}
	fmt.Println(cmd.error().String())
	os.Exit(1)
}

// FlagFatal returns flag options when an invalid choice is used.
func FlagFatal(name, value string, args ...string) {
	cmd := Cmd{Args: args, Err: fmt.Errorf("%q: %w", name, ErrSlice)}
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
		LogFatal(fmt.Errorf("cmd error args: %w", ErrShort))
	} else if len(c.Args) < 1 {
		LogFatal(fmt.Errorf("cmd error err: %w", ErrEmpty))
	}
	a := fmt.Sprintf("%q", c.Args[0])
	switch strings.Join(s[0:2], " ") {
	case "bad flag":
		return Generic{Issue: "flag syntax",
			Arg: quote(s[l-1]),
			Err: ErrSyntax}
	case "flag needs":
		return Generic{Issue: "invalid flag",
			Arg: quote(s[l-1]),
			Err: ErrNoFlag}
	case "invalid argument":
		m := strings.Split(fmt.Sprint(c.Err), ":")
		return Generic{Issue: "flag value",
			Arg: a,
			Err: fmt.Errorf("%s: %w", m[0], ErrNoCmd)}
	case "invalid slice":
		return Generic{Issue: "flag value",
			Arg: quote(s[l-1]),
			Err: fmt.Errorf("--%s: %w", s[l-2], ErrVal)}
	case "invalid command":
		return Generic{Issue: "invalid command",
			Arg: quote(s[l-1]),
			Err: ErrNewCmd}
	case "required flag(s)":
		return Generic{Issue: "a required flag is missing",
			Arg: s[2],
			Err: ErrReqFlag}
	case "subcommand is":
		fmt.Printf("SUBCMD DEBUG: %+v", c.Err)
		return Generic{} // ignore error
	case "unknown command":
		return Generic{Issue: "invalid command",
			Arg: a,
			Err: ErrCmd}
	case "unknown flag:":
		return Generic{Issue: "unknown flag",
			Arg: s[2],
			Err: ErrFlag}
	case "unknown shorthand":
		return Generic{Issue: "unknown shorthand flag",
			Arg: s[5],
			Err: ErrFlag}
	}
	//fmt.Printf("DEBUG: %+v\n", c.Err)
	return Generic{Issue: "command", Arg: "execute", Err: c.Err}
}
