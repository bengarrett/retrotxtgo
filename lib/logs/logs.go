package logs

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/gookit/color"
	gap "github.com/muesli/go-app-paths"
)

// the lowest, largest and recommended network ports to serve HTTP.
const (
	PortMin uint = 0
	PortMax uint = 65535
	PortRec uint = 8080
)

// Err is an interface for error messages
type Err struct {
	Issue string
	Arg   string
	Msg   error
}

func (e Err) String() string {
	ia := Ci(fmt.Sprintf("%s %s,", e.Issue, e.Arg))
	m := Cf(fmt.Sprintf("%v", e.Msg))
	return fmt.Sprintf("%s %s %s", Alert(), ia, m)
}

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

// color aliases
var (
	Alert = func() string {
		return color.Error.Sprint("problem:")
	}
	Cc = func(t string) string {
		return color.Comment.Sprint(t)
	}
	Ce = func(t string) string {
		return color.Warn.Sprint(t)
	}
	Cf = func(t string) string {
		return color.OpFuzzy.Sprint(t)
	}
	Ci = func(t string) string {
		return color.OpItalic.Sprint(t)
	}
	Cinf = func(t string) string {
		return color.Info.Sprint(t)
	}
	Cp = func(t string) string {
		return color.Primary.Sprint(t)
	}
	Cs = func(t string) string {
		return color.Success.Sprint(t)
	}
)

// ChkArg returns instructions for invalid command arguments.
func ChkArg(arg string, args []string) {
	if len(args) == 0 {
		return
	}
	fmt.Printf("%s invalid argument%s",
		Alert(),
		color.Bold.Sprintf(" %q", arg))
	if len(args) > 1 {
		fmt.Printf(" choices: %s\n%s",
			color.Info.Sprintf("%s", strings.Join(args, ", ")),
			color.Warn.Sprint("please use one of the argument choices shown above"))
	}
	fmt.Println()
	os.Exit(10)
}

// ChkErr prints an error issue and message then exits the program.
func ChkErr(issue string, err error) {
	if err != nil {
		if issue == "" {
			fmt.Printf("%s\n", err)
		} else {
			fmt.Printf("%s %s\n", issue, err)
		}
		os.Exit(1)
	}
}

// Check prints an error message and exits the program.
func Check(e Err) {
	if e.Msg != nil {
		println(e.String())
		os.Exit(1)
	}
}

// ColorHTML prints colored syntax highlighting to HTML elements.
func ColorHTML(elm string) {
	var buf bytes.Buffer
	if err := quick.Highlight(&buf, elm, "html", "terminal256", "lovelace"); err != nil {
		fmt.Printf("\n%s\n", elm)
	} else {
		fmt.Printf("\n%v\n", buf.String())
	}
}

// Save logs any errors and exits to the operating system with error code 1.
func Save(err error) {
	if err != nil {
		save(err, "")
		switch Panic {
		case true:
			log.Println(fmt.Sprintf("error type: %T\tmsg: %v", err, err))
			log.Panic(err)
		default:
			log.Fatal(color.Danger.Sprint("ERROR: "), err)
		}
	}
}

// save an error to the logs, optional path is available for unit tests.
func save(err error, path string) (ok bool) {
	if err == nil || fmt.Sprintf("%v", err) == "" {
		return false
	}
	// use UTC date and times in the log file
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	if path == "" {
		path = Filepath()
	}
	p := filepath.Dir(path)
	if _, e := os.Stat(p); os.IsNotExist(e) {
		e2 := os.MkdirAll(p, permDir)
		check(e2)
	}
	file, e := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	check(e)
	defer file.Close()
	log.SetOutput(file)
	log.Print(err)
	log.SetOutput(os.Stderr)
	return true
}

func check(e error) {
	if e != nil {
		log.Printf("%s %s", color.Danger.Sprint("!"), e)
		os.Exit(19)
	}
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
