package logs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// InvalidCommand prints a problem highlighting the unsupported command.
func InvalidCommand(usage string, args ...string) {
	err := fmt.Errorf("%w: %s", ErrCmdExist, args[0])
	args = append(args, usage)
	Execute(err, args...)
}

func InvalidChoice(name, value string, choices ...string) {
	c := ProblemFlag(name, value, ErrFlagChoice)
	fmt.Println(c)
	fmt.Printf("          choices: %s\n", strings.Join(choices, ", "))
	os.Exit(1)
}

// Execute is the error handler for command flags and arguments.
func Execute(err error, args ...string) {
	const (
		minWords       = 3
		rt             = "retrotxt"
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
		SaveFatal(fmt.Errorf("cmd error args: %w", ErrShort))
	}
	if argsCnt == 0 {
		SaveFatal(fmt.Errorf("cmd error err: %w", ErrEmpty))
	}

	mark, name := words[wordCnt-1], args[0]
	if mark == name {
		name = rt
	}
	problem := strings.Join(words[0:2], " ")
	var c string
	switch problem {
	case flagSyntax:
		c = ProblemFlag(name, mark, err)
	case invalidFlag:
		// retroxt config shell -i
		c = ProblemFlag(name, mark, ErrNotNil)
	case invalidType:
		// retroxt --help=foo
		mark = strings.Join(words[4:6], " ")
		c = parseType(name, mark, err)
	case invalidSlice:
		c = "invalidSlice placeholder"
	case invalidCommand:
		// retrotxt config foo
		c = Hint(fmt.Sprintf("%s --help", mark), ErrCmdExist)
	case flagRequired:
		// retrotxt config shell
		c = ProblemCmd(mark, ErrFlagNil)
	case unknownCmd:
		// retrotxt foo
		mark = words[2]
		c = Hint("--help", fmt.Errorf("%w: %s", ErrCmdExist, mark))
	case unknownFlag:
		// retrotxt --foo
		mark = words[2]
		if mark == name {
			name = rt
		}
		c = ProblemFlag(name, mark, ErrFlag)
	case unknownShort:
		// retrotxt -foo
		mark = words[5]
		if mark == name {
			name = rt
		}
		c = ProblemFlag(name, mark, ErrFlag)
	case flagChoice:
		c = "flagChoice placeholder"
	default:
		if errors.As(err, &ErrCmdExist) {
			mark = strings.Join(args[1:], " ")
			c = Hint(fmt.Sprintf("%s --help", mark), ErrCmdExist)
			break
		}
		c = Errorf(err)
	}
	if c != "" {
		fmt.Println(c)
		os.Exit(1)
	}
	log.Panic(err)
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
		return ProblemFlag(name, flag, ErrNotBool)
	case strings.Contains(s, invalidInt):
		return ProblemFlag(name, flag, ErrNotInt)
	case strings.Contains(s, invalidStr):
		return ProblemFlag(name, flag, ErrNotInts)
	default:
		return ProblemFlag(name, flag, err)
	}
}
