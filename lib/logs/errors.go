package logs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"retrotxt.com/retrotxt/lib/str"
)

var (
	// generic errors
	ErrConfigFile = errors.New("could not open the configuration file")
	ErrEncode     = errors.New("text encoding not known or supported")
	ErrHighlight  = errors.New("could not format or colorize the element")
	ErrOpenFile   = errors.New("could not open the file")
	ErrPipe       = errors.New("could not read text stream from piped stdin (standard input)")
	ErrPipeParse  = errors.New("could not parse the text stream from piped stdin (standard input)")
	ErrSaveDir    = errors.New("could not save file to the directory")
	ErrSaveFile   = errors.New("could not save the file")
	ErrSampFile   = errors.New("unknown sample filename")
	ErrTmpClean   = errors.New("could not cleanup the temporary directory")
	ErrTmpDir     = errors.New("could not save file to the temporary directory")
	ErrTmpSave    = errors.New("could not save the temporary file")
	ErrTabFlush   = errors.New("tab writer could not write or flush")
	ErrZipFile    = errors.New("could not create the zip archive")
	// command (cmd library) and argument errors
	ErrHelp        = errors.New("command help could not display")
	ErrMarkRequire = errors.New("command flag could not be marked as required")
	ErrUsage       = errors.New("command usage could not display")
	// type errors
	ErrCmd     = errors.New("choose a command from the list available")
	ErrNewCmd  = errors.New("choose another command from the available commands")
	ErrNoCmd   = errors.New("invalid command")
	ErrEmpty   = errors.New("value is empty")
	ErrFlag    = errors.New("use a flag from the list of flags")
	ErrSyntax  = errors.New("flags can only be in -s (short) or --long (long) form")
	ErrNoFlagx = errors.New("cannot be empty and requires a value")
	ErrReqFlag = errors.New("you must include this flag in your command")
	ErrSlice   = errors.New("invalid option choice")
	ErrShort   = errors.New("word count is too short, less than 3")
	ErrVal     = errors.New("value is not a valid choice")
)

// Argument is used to highlight user supplied argument values in error replies.
type Argument struct {
	Issue string // Issue is a human readable summary of the problem
	Value string // Value of the the argument, flag or item that triggered the error
	Err   error  // Err is required error
}

func CmdProblem(name string, err error) string {
	alert := str.Alert()
	return fmt.Sprintf("%s the command %s does not exist, %s", alert, name, err)
}

func CmdProblemFatal(name, flag string, err error) {
	fmt.Println(FlagProblem(name, flag, err))
	os.Exit(1)
}

// rename to Fatal
func ErrorFatal(err error) {
	Errorf(err)
	os.Exit(1)
}

func Errorf(err error) {
	fmt.Printf("%s%s\n", str.Alert(), err) // TODO: change to string return?
}

func FlagProblem(name, flag string, err error) string {
	alert, toggle := str.Alert(), "--"
	fmt.Println("FLAG:", flag)
	if strings.Contains(flag, "-") {
		toggle = ""
	} else if len(flag) == 1 {
		toggle = "-"
	}
	return fmt.Sprintf("%s with the %s %s%s flag, %s", alert, name, toggle, flag, err)
}

func MarkProblem(value string, new, err error) string {
	v := value
	a := str.Alert()
	n := str.Cf(fmt.Sprintf("%v", new))
	e := str.Cf(fmt.Sprintf("%v", err))

	return fmt.Sprintf("%s %s %q: %s", a, n, v, e)
}

func MarkProblemFatal(value string, new, err error) {
	fmt.Println(MarkProblem(value, new, err))
	os.Exit(1)
}

// ok
func Problemln(new, err error) {
	e := fmt.Errorf("%s: %w", new, err)
	fmt.Printf("%s%s\n", str.Alert(), e)
}

// ok
func ProblemFatal(new, err error) {
	Problemln(new, err)
	os.Exit(1)
}

func SubCmdProblem(name string, err error) string {
	alert := str.Alert()
	return fmt.Sprintf("%s the subcommand %s does not exist, %s", alert, name, err)
}

// Fatal prints a generic error and exits.
func Fatal(issue, arg string, msg error) {
	Println(issue, arg, msg)
	os.Exit(1)
}

// Println prints a generic error.
// TODO, make into (a Argument) Println ?
// (a Argument) Println() {}
func Println(i, v string, err error) {
	var g = Argument{
		Issue: i,
		Value: v,
		Err:   err,
	}
	fmt.Println(g.String())
}

func (g Argument) String() string {
	if g.Err == nil {
		return ""
	}
	a, c := str.Alert(), str.Cf(fmt.Sprintf("%v", g.Err))

	if s := g.unWrap(); s != "" {
		return fmt.Sprintf("%s %s", a, s)
	}

	if g.Issue == "" && g.Value == "" {
		return fmt.Sprintf("%s %s", a, c) // alert and err
	}
	var b string
	if g.Value == "" {
		b = str.Ci(fmt.Sprintf("%s,", g.Issue)) // alert, issue and err
	} else {
		b = str.Ci(fmt.Sprintf("%s %s,", g.Issue, g.Value)) // alert, issue, arg, err
	}
	return fmt.Sprintf("%s %s %s", a, b, c)
}

func (g Argument) unWrap() string {
	var fp *fs.PathError
	uw := errors.Unwrap(g.Err)
	if uw == nil {
		return ""
	}
	if errors.As(uw, &fp) {
		return fmt.Sprintf("cannot open file %q, %s", g.Value, str.Cf("is there a typo?"))
	}

	fmt.Printf("\n%T %+v\n", uw, uw)
	return ""
}

// Fatal prints a generic error and exits.
func (g Argument) XFatal() {
	fmt.Println(g.String())
	os.Exit(1)
}

// Hint is a standard error type that also offers the user a command hint.
type Hint struct {
	Error Argument
	Hint  string // Hint is an optional solution such as a retrotxt command
}

func (h Hint) String() string {
	err := h.Error
	return err.String() + fmt.Sprintf("\n         run %s", str.Example("retrotxt "+h.Hint))
}
