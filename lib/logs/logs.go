package logs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gookit/color"
	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/str"
)

// the lowest, largest and recommended network ports to serve HTTP.
const (
	PortMin uint = 0
	PortMax uint = 65535
	PortRec uint = 8080
)

const (
	// Filename is the default error log filename
	Filename = "errors.log"
	// posix permissions for the configuration file and directory.
	permf os.FileMode = 0600
	permd os.FileMode = 0700
)

var (
	scope = gap.NewScope(gap.User, "df2")
	// Panic uses the panic function to handle all error logs.
	Panic = false
)

// ColorCSS prints colored CSS syntax highlighting.
func ColorCSS(elm string) string {
	style := viper.GetString("style.html")
	return colorElm(elm, "css", style)
}

// ColorHTML prints colored syntax highlighting to HTML elements.
func ColorHTML(elm string) string {
	style := viper.GetString("style.html")
	return colorElm(elm, "html", style)
}

func colorElm(elm, lexer, style string) string {
	if elm == "" {
		return ""
	}
	var b bytes.Buffer
	_ = io.Writer(&b)
	if err := str.HighlightWriter(&b, elm, lexer, style); err != nil {
		Fatal("logs", "colorhtml", err)
	}
	return fmt.Sprintf("\n%s\n", b.String())
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

// LogCont logs the error and continues the program.
func LogCont(err error) {
	if err != nil {
		// save error to log file
		if err := save(err, ""); err != nil {
			log.Fatalf("%s %s", color.Danger.Sprint("!"), err)
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
		name = Path()
	}
	p := filepath.Dir(name)
	if _, e := os.Stat(p); os.IsNotExist(e) {
		if e := os.MkdirAll(p, permd); e != nil {
			return e
		}
	}
	file, e := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, permf)
	if e != nil {
		return e
	}
	defer file.Close()
	log.SetOutput(file)
	log.Print(err)
	log.SetOutput(os.Stderr)
	return nil
}

// Path is the absolute path and filename of the error log file.
func Path() string {
	fp, err := scope.LogPath(Filename)
	if err != nil {
		h, _ := os.UserHomeDir()
		return path.Join(h, Filename)
	}
	return fp
}
