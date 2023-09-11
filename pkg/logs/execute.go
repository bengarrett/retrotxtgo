package logs

import (
	"errors"
	"fmt"
	"log"
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
//
//nolint:funlen,cyclop
func Execute(err error, test bool, args ...string) string {
	if err == nil {
		return ""
	}
	const (
		minWords       = 3
		rt             = meta.Bin
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

		exec = "logs executer problem"
	)
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
		// do nothing
	}
	// 	e := fmt.Errorf("%s: %w", exec, ErrNoArgs)
	// 	if test {
	// 		return e.Error()
	// 	}
	// 	Fatal(e)
	// }
	mark, name := words[wordCnt-1], args[0]
	if mark == name {
		name = rt
	}
	problem := strings.Join(words[0:2], " ")
	var c string
	switch problem {
	case flagSyntax:
		c = SprintFlag(err, name, mark)
	case invalidFlag:
		// retroxt config shell -i
		c = SprintFlag(ErrNotNil, name, mark)
	case invalidType:
		// retroxt --help=foo
		const min = 6
		if len(words) >= min {
			mark = strings.Join(words[4:6], " ")
		}
		c = parseType(name, mark, err)
	case invalidSlice:
		c = "invalidSlice placeholder"
	case invalidCommand:
		// retrotxt config foo
		c = Hint(ErrCmd, fmt.Sprintf("%s --help", mark))
	case flagRequired:
		// retrotxt config shell
		c = SprintCmd(ErrFlagNil, mark)
	case unknownCmd:
		// retrotxt foo
		mark = words[2]
		c = Hint(fmt.Errorf("%w: %s", ErrCmd, mark), "--help")
	case unknownFlag:
		// retrotxt --foo
		mark = words[2]
		if mark == name {
			name = rt
		}
		c = SprintFlag(ErrFlag, name, mark)
	case unknownShort:
		// retrotxt -foo
		mark = words[5]
		if mark == name {
			name = rt
		}
		c = SprintFlag(ErrFlag, name, mark)
	case flagChoice:
		c = "flagChoice placeholder"
	default:
		if errors.Is(err, ErrCmd) {
			mark = strings.Join(args[1:], " ")
			c = Hint(ErrCmd, fmt.Sprintf("%s --help", mark))
			break
		}
		c = Sprint(err)
	}
	if c != "" {
		return c
	}
	log.Panic(err)
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
