package logs

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrCmdExist   = errors.New("the command does not exist")
	ErrCmdChoose  = errors.New("choose an available command from the --help list")
	ErrFlagExist  = errors.New("the flag does not work with this command")
	ErrFlagChoice = errors.New("choose a value from the following")
	ErrNoFlag     = errors.New("the flag with a value must be included with this command")
	ErrNoFlagVal  = errors.New("the value cannot be left empty")
	ErrNotBool    = errors.New("the value must be either true or false")
	ErrNotInt     = errors.New("the value must be a number")
	ErrNotInts    = errors.New("the value must be a single or a list of numbers")
)

// InvalidCommand prints a problem highlighting the unsupported command.
func InvalidCommand(args ...string) {
	err := fmt.Errorf("invalid command %s", args[0])
	Execute(err, args...)
}

func InvalidChoice(name, value string, choices ...string) {
	c := FlagProblem(name, value, ErrFlagChoice)
	fmt.Println(c)
	fmt.Printf("          choices: %s\n", strings.Join(choices, ", "))
	os.Exit(1)
}

// Execute is the error handler for command flags and arguments.
func Execute(err error, args ...string) {
	const (
		minWords       = 3
		flagChoice     = "invalid option choice"
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
	case flagChoice:
		c = "honk"
	}
	if c != "" {
		fmt.Println(c)
		os.Exit(1)
	}
	// TODO: handle unknown/empty
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
