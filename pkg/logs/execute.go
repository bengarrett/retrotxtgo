package logs

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/meta"
)

var (
	ErrNoArgs  = errors.New("no arguments were given for the logs executer")
	ErrShort   = errors.New("word count is too short, it requires at least 3 words")
	ErrCmd     = errors.New("the command is invalid")
	ErrFlag    = errors.New("the flag does not work with this command")
	ErrFlagNil = errors.New("the flag with a value must be included with this command")
	ErrNotBool = errors.New("the value must be either true or false")
	ErrNotInt  = errors.New("the value must be a number")
	ErrNotInts = errors.New("the value must be a single or a list of numbers")
	ErrNotNil  = errors.New("the value cannot be empty")
)

const (
	minWords       = 3
	flagChoice     = "invalid option choice"
	flagRequired   = "required flag(s)"
	flagSyntax     = "bad flag"
	invalidCommand = "invalid command"
	invalidFlag    = "flag needs"
	invalidSlice   = "invalid slice"
	invalidType    = "invalid argument"
	exec           = "logs executer problem"
)

// FatalSubCmd prints a problem highlighting the unsupported sub-command.
func FatalSubCmd(usage string, args ...string) {
	args = append(args, usage)
	var err error
	if len(args) > 0 {
		err = fmt.Errorf("%w: %s", ErrCmd, args[0])
	}
	if s := Execute(err, false, args...); s != "" {
		fmt.Fprintln(os.Stderr, s)
		os.Exit(OSErr)
	}
}

// Execute is the error handler for command flags and arguments.
func Execute(err error, test bool, args ...string) string {
	if err == nil {
		return ""
	}
	words := strings.Split(fmt.Sprintf("%s", err), " ")
	argsCnt, wordCnt := len(args), len(words)
	if wordCnt < minWords {
		e := fmt.Errorf("%s: %w", exec, ErrShort)
		if test {
			return e.Error()
		}
		Fatal(e)
	}
	if argsCnt == 0 {
		Fatal(err)
	}
	mark, name := words[wordCnt-1], args[0]
	if mark == name {
		name = meta.Bin
	}
	if x := invalid(err, mark, name, words...); x != "" {
		return x
	}
	if x := unknown(name, words...); x != "" {
		return x
	}
	if errors.Is(err, ErrCmd) {
		mark = strings.Join(args[1:], " ")
		return Hint(ErrCmd, fmt.Sprintf("%s --help", mark))
	}
	return Sprint(err)
}

func invalid(err error, mark, name string, words ...string) string {
	problem := strings.Join(words[0:2], " ")
	switch problem {
	case flagSyntax:
		return SprintFlag(err, name, mark)
	case flagRequired: // retrotxt config shell
		return SprintCmd(ErrFlagNil, mark)
	case flagChoice:
		return "flagChoice placeholder"
	case invalidFlag: // retroxt config shell -i
		return SprintFlag(ErrNotNil, name, mark)
	case invalidType: // retroxt --help=foo
		const min = 6
		if len(words) >= min {
			mark = strings.Join(words[4:6], " ")
		}
		return parseType(name, mark, err)
	case invalidSlice:
		return "invalidSlice placeholder"
	case invalidCommand: // retrotxt config foo
		return Hint(ErrCmd, fmt.Sprintf("%s --help", mark))
	}
	return ""
}

func unknown(name string, words ...string) string {
	const (
		unknownCmd   = "unknown command"
		unknownFlag  = "unknown flag:"
		unknownShort = "unknown shorthand"
		rt           = meta.Bin
	)
	problem := strings.Join(words[0:2], " ")
	switch problem {
	case unknownCmd: // retrotxt foo
		return Hint(fmt.Errorf("%w: %s", ErrCmd, words[2]), "--help")
	case unknownFlag: // retrotxt --foo
		mark := words[2]
		if mark == name {
			name = rt
		}
		return SprintFlag(ErrFlag, name, mark)
	case unknownShort: // retrotxt -foo
		mark := words[5]
		if mark == name {
			name = rt
		}
		return SprintFlag(ErrFlag, name, mark)
	}
	return ""
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
		return SprintFlag(ErrNotBool, name, flag)
	case strings.Contains(s, invalidInt):
		return SprintFlag(ErrNotInt, name, flag)
	case strings.Contains(s, invalidStr):
		return SprintFlag(ErrNotInts, name, flag)
	default:
		return SprintFlag(err, name, flag)
	}
}
