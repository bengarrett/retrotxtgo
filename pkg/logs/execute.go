package logs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/meta"
)

// FatalCmd prints a problem highlighting the unsupported command.
func FatalCmd(usage string, args ...string) {
	args = append(args, usage)
	var err error
	if len(args) > 0 {
		err = fmt.Errorf("%w: %s", ErrCmd, args[0])
	}
	if s := Execute(err, false, args...); s != "" {
		fmt.Fprintln(os.Stderr, s)
		os.Exit(OSErrCode)
	}
}

// FatalExecute is the error handler for the root command flags and arguments.
func FatalExecute(err error, args ...string) {
	if s := Execute(err, false, args...); s != "" {
		fmt.Fprintln(os.Stderr, s)
		os.Exit(OSErrCode)
	}
}

// Execute is the error handler for command flags and arguments.
func Execute(err error, test bool, args ...string) string { //nolint:funlen
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
	)
	words := strings.Split(fmt.Sprintf("%s", err), " ")
	argsCnt, wordCnt := len(args), len(words)
	if wordCnt < minWords {
		e := fmt.Errorf("cmd error args: %w", ErrShort)
		if test {
			return e.Error()
		}
		FatalSave(e)
	}
	if argsCnt == 0 {
		e := fmt.Errorf("cmd error err: %w", ErrEmpty)
		if test {
			return e.Error()
		}
		FatalSave(e)
	}
	mark, name := words[wordCnt-1], args[0]
	if mark == name {
		name = rt
	}
	problem := strings.Join(words[0:2], " ")
	var c string
	switch problem {
	case flagSyntax:
		c = SprintFlag(name, mark, err)
	case invalidFlag:
		// retroxt config shell -i
		c = SprintFlag(name, mark, ErrNotNil)
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
		c = Hint(fmt.Sprintf("%s --help", mark), ErrCmd)
	case flagRequired:
		// retrotxt config shell
		c = SprintCmd(mark, ErrFlagNil)
	case unknownCmd:
		// retrotxt foo
		mark = words[2]
		c = Hint("--help", fmt.Errorf("%w: %s", ErrCmd, mark))
	case unknownFlag:
		// retrotxt --foo
		mark = words[2]
		if mark == name {
			name = rt
		}
		c = SprintFlag(name, mark, ErrFlag)
	case unknownShort:
		// retrotxt -foo
		mark = words[5]
		if mark == name {
			name = rt
		}
		c = SprintFlag(name, mark, ErrFlag)
	case flagChoice:
		c = "flagChoice placeholder"
	default:
		if errors.Is(err, ErrCmd) {
			mark = strings.Join(args[1:], " ")
			c = Hint(fmt.Sprintf("%s --help", mark), ErrCmd)
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
		return SprintFlag(name, flag, ErrNotBool)
	case strings.Contains(s, invalidInt):
		return SprintFlag(name, flag, ErrNotInt)
	case strings.Contains(s, invalidStr):
		return SprintFlag(name, flag, ErrNotInts)
	default:
		return SprintFlag(name, flag, err)
	}
}
