package logs

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/viper"
)

// the lowest, largest and recommended network ports to serve HTTP.
const (
	PortMin uint = 0
	PortMax uint = 65535
	PortRec uint = 8080
)

// Filename is the default error log filename
const Filename = "errors.log"

// posix permissions for the configuration file and directory.
const perm = 0600
const permDir = 0700

var (
	scope = gap.NewScope(gap.User, "df2")
	// Panic uses the panic function to handle all error logs.
	Panic = false
)

// FileMissingErr exits with a missing FILE error.
// TODO: probably needs replacing
func FileMissingErr() {
	i := Ci("missing the --name flag")
	m := Cf("you need to provide a path to a text file")
	fmt.Printf("\n%s %s %s\n", Alert(), i, m)
	os.Exit(1)
}

// ReCheck is a temp function used by cmd functions that will be replaced.
func ReCheck(err error) {
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}

// Check prints an error issue and message then exits the program.
func Check(issue string, err error) (ok bool) {
	if err != nil {
		Exit(check(issue, err))
	}
	return true
}

func check(issue string, err error) (msg string, code int) {
	code = 1
	if err == nil {
		err = errors.New("")
	}
	if issue == "" {
		msg = fmt.Sprintf("%s\n", err)
	} else {
		msg = fmt.Sprintf("%s %s\n", issue, err)
	}
	return msg, code
}

// CheckArg returns instructions for invalid command arguments.
func CheckArg(arg string, args []string) {
	if len(args) == 0 {
		return
	}
	Exit(checkArgument(arg, args))
}

func checkArgument(arg string, args []string) (msg string, code int) {
	code = 10
	msg += fmt.Sprintf("%s invalid argument%s", Alert(),
		color.Bold.Sprintf(" %q", arg))
	if len(args) > 1 {
		msg += fmt.Sprintf(", choices: %s\n%s",
			color.Info.Sprintf("%s", strings.Join(args, ", ")),
			color.Warn.Sprint("please use one of the ")+
				color.Info.Sprint("argument choices")+
				color.Warn.Sprint(" shown above"))
	}
	return msg, code
}

// CheckNilArg returns instructions for empty command arguments.
func CheckNilArg(arg string, args []string) {
	if len(args) != 0 {
		return
	}
	Exit(checkNilArgument(arg, args))
}

func checkNilArgument(arg string, args []string) (msg string, code int) {
	msg += fmt.Sprintf("%s required argument%s cannot be empty and requires a value\n", Alert(), color.Bold.Sprintf(" %q", arg))
	return msg, code
}

// ChkErr prints an error message and exits the program.
func ChkErr(e Err) {
	if e.Msg != nil {
		Exit(e.check())
	}
}

func (e Err) check() (msg string, code int) {
	code = 1
	msg = e.String()
	return msg, code
}

// ColorHTML prints colored syntax highlighting to HTML elements.
func ColorHTML(elm string) {
	fmt.Println(colorhtml(&elm))
}

func colorhtml(elm *string) string {
	var buf bytes.Buffer
	style := viper.GetString("style.html")
	if err := HighlightWriter(&buf, *elm, "html", style); err != nil {
		return fmt.Sprintf("\n%s\n", *elm)
	}
	if len(buf.String()) == 0 {
		return ""
	}
	return fmt.Sprintf("\n%v\n", buf.String())
}

// Exit prints the message and causes the program to exit.
// TODO: make into an error method
func Exit(msg string, code int) {
	i, err := fmt.Println(msg)
	if err != nil {
		log.Fatalf("logs.exit println at %dB: %s", i, err)
	}
	os.Exit(code)
}

// Log the error and exit to the operating system with the error code 1.
func Log(err error) {
	if err != nil {
		// save error to log file
		if err := save(err, ""); err != nil {
			log.Fatalf("%s %s", color.Danger.Sprint("!"), err)
		}
		// print error
		switch Panic {
		case true:
			log.Println(fmt.Sprintf("error type: %T\tmsg: %v", err, err))
			log.Panic(err)
		default:
			log.Fatal(color.Danger.Sprint("ERROR: "), err)
		}
	}
}

// save an error to the log directory, an optional named file is available for unit tests.
func save(err error, name string) error {
	if err == nil || fmt.Sprintf("%v", err) == "" {
		return errors.New("logs save: err value is nil")
	}
	// use UTC date and times in the log file
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	if name == "" {
		name = Filepath()
	}
	p := filepath.Dir(name)
	if _, e := os.Stat(p); os.IsNotExist(e) {
		if e := os.MkdirAll(p, permDir); e != nil {
			return e
		}
	}
	file, e := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if e != nil {
		return e
	}
	defer file.Close()
	log.SetOutput(file)
	log.Print(err)
	log.SetOutput(os.Stderr)
	return nil
}

// Filepath is the absolute path and filename of the error log file.
func Filepath() string {
	fp, err := scope.LogPath(Filename)
	if err != nil {
		h, _ := os.UserHomeDir()
		return path.Join(h, Filename)
	}
	return fp
}
