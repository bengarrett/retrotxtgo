package logs

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Cmd is a command error type to handle command arguments and flags.
type Cmd struct {
	Args []string // Arguments
	Err  error    // rootCmd.Execute output
}

// Execute is the error handler for command flags and arguments.
func Execute(err error, args ...string) {
	//Parse(err, args...)
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
	cmd := Cmd{Args: args, Err: fmt.Errorf("--%s=%q: %w", name, value, ErrSlice)}
	fmt.Println(cmd.error().String())
	fmt.Printf("         choices: %s\n", strings.Join(args, ", "))
	os.Exit(1)
}

// Parse placeholder
func Parse(err error, args ...string) {
	uw := errors.Unwrap(err)
	fmt.Println("uw", uw)
	fmt.Println("err", err)
	s := fmt.Sprint(err)
	switch {
	case strings.Contains(s, "strconv.ParseInt"):
		fmt.Println("hehhh")
	}
}

//
func (c Cmd) error() Argument {
	s := strings.Split(fmt.Sprintf("%s", c.Err), " ")
	l := len(s)
	quote := func(s string) string {
		return fmt.Sprintf("%q", s)
	}
	const min = 3
	if l < min {
		LogFatal(fmt.Errorf("cmd error args: %w", ErrShort))
	} else if len(c.Args) < 1 {
		LogFatal(fmt.Errorf("cmd error err: %w", ErrEmpty))
	}
	a := fmt.Sprintf("%q", c.Args[0])
	switch strings.Join(s[0:2], " ") {
	case "bad flag":
		return Argument{Issue: "flag syntax",
			Value: quote(s[l-1]),
			Err:   ErrSyntax}
	case "flag needs":
		return Argument{Issue: "invalid flag",
			Value: quote(s[l-1]),
			Err:   ErrNoFlag}
	case "invalid argument":
		m := strings.Split(fmt.Sprint(c.Err), ":")
		return Argument{Issue: "flag value",
			Value: a,
			Err:   fmt.Errorf("%s: %w", m[0], ErrNoCmd)}
	case "invalid slice":
		return Argument{Issue: "flag value",
			Value: quote(s[l-1]),
			Err:   fmt.Errorf("--%s: %w", s[l-2], ErrVal)}
	case "invalid command":
		return Argument{Issue: "invalid command",
			Value: quote(s[l-1]),
			Err:   ErrNewCmd}
	case "required flag(s)":
		return Argument{Issue: "a required flag is missing",
			Value: s[2],
			Err:   ErrReqFlag}
	case "subcommand is":
		fmt.Printf("SUBCMD DEBUG: %+v", c.Err)
		return Argument{} // ignore error
	case "unknown command":
		return Argument{Issue: "invalid command",
			Value: a,
			Err:   ErrCmd}
	case "unknown flag:":
		return Argument{Issue: "unknown flag",
			Value: s[2],
			Err:   ErrFlag}
	case "unknown shorthand":
		return Argument{Issue: "unknown shorthand flag",
			Value: s[5],
			Err:   ErrFlag}
	}
	return Argument{Issue: "command", Value: "execute", Err: c.Err}
}
