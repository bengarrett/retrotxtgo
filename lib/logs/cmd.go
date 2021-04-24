package logs

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrCmdExist  = errors.New("the command does not exist")
	ErrCmdChoose = errors.New("choose an available command from the --help list")
	ErrFlagExist = errors.New("the flag does not work with this command")
	ErrNoFlag    = errors.New("the flag with a value must be included with this command")
	ErrNoFlagVal = errors.New("the value cannot be left empty")
	ErrNotBool   = errors.New("the value must be either true or false")
	ErrNotInt    = errors.New("the value must be a number")
	ErrNotInts   = errors.New("the value must be a single or a list of numbers")
)

// Cmd is a command error type to handle command arguments and flags.
type Cmd struct {
	Args []string // Arguments
	Err  error    // rootCmd.Execute output
}

// func Execute(err error, args ...string) {
// 	Parse(err, args...)
// 	cmd := Cmd{Args: args, Err: err}
// 	fmt.Println(cmd.error().String())
// }

func InvalidCommand(args ...string) {
	err := fmt.Errorf("invalid command %s", args[0])
	Execute(err, args...)
}

// FlagFatal returns flag options when an invalid choice is used.
func FlagFatal(name, value string, args ...string) {
	cmd := Cmd{Args: args, Err: fmt.Errorf("--%s=%q: %w", name, value, ErrSlice)}
	fmt.Println(cmd.error().String())
	fmt.Printf("         choices: %s\n", strings.Join(args, ", "))
	os.Exit(1)
}

// Execute is the error handler for command flags and arguments.
func Execute(err error, args ...string) {
	const (
		minWords       = 3
		flagRequired   = "required flag(s)"
		flagSyntax     = "bad flag"
		invalidCommand = "invalid command"
		invalidFlag    = "flag needs"
		invalidSlice   = "invalid slice"
		invalidType    = "invalid argument"
		unknownCmd     = "unknown command"
		unknownFlag    = "unknown flag:"
		unknownShort   = "unknown shorthand"
	)

	words := strings.Split(fmt.Sprintf("%s", err), " ")
	argsCnt, wordCnt := len(args), len(words)
	if wordCnt < minWords {
		LogFatal(fmt.Errorf("cmd error args: %w", ErrShort))
	}
	if argsCnt == 0 {
		LogFatal(fmt.Errorf("cmd error err: %w", ErrEmpty))
	}

	mark, name := words[wordCnt-1], args[0]
	problem := strings.Join(words[0:2], " ")
	var c string
	switch problem {
	case flagSyntax:
		c = FlagProblem(name, mark, err)
	case invalidFlag:
		c = FlagProblem(name, mark, ErrNoFlagVal)
	case invalidType:
		mark = strings.Join(words[4:6], " ")
		c = parseType(name, mark, err)
	case invalidSlice:
		return // TODO:
	case invalidCommand:
		c = SubCmdProblem(mark, ErrCmdChoose)
	case flagRequired:
		c = CmdProblem(mark, ErrNoFlag)
	case unknownCmd:
		mark = words[2]
		c = CmdProblem(mark, ErrCmdChoose)
	case unknownFlag, unknownShort:
		c = FlagProblem(name, mark, ErrFlagExist)
	}
	if c != "" {
		fmt.Println(c)
		os.Exit(1)
	}

}

func parseType(name, flag string, err error) string {
	const (
		invalidBool = "strconv.ParseBool"
		invalidInt  = "strconv.ParseInt"
		invalidStr  = "strconv.Atoi"
	)
	s := err.Error()
	switch {
	case strings.Contains(s, invalidBool):
		return FlagProblem(name, flag, ErrNotBool)
	case strings.Contains(s, invalidInt):
		return FlagProblem(name, flag, ErrNotInt)
	case strings.Contains(s, invalidStr):
		return FlagProblem(name, flag, ErrNotInts)
	default:
		return FlagProblem(name, flag, err)
	}
}

// func CmdProblem(name, flag string, err error) string {
// 	alert := str.Alert()
// 	return fmt.Sprintf("%s %s flag --%s %s", alert, name, flag, err)
// }

// func CmdProblemFatal(name, flag string, err error) {

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
	fmt.Println("Cmd.error(): ", s)
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
