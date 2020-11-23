// Package logs for the display of errors.
package logs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gookit/color"
	gap "github.com/muesli/go-app-paths"
)

const (
	// default error log.
	filename = "errors.log"
	// posix permissions for the configuration file and directory.
	filemode os.FileMode = 0600
	dirmode  os.FileMode = 0700
	// Panic uses log.Panic to print all saved errors.
	Panic = false
)

// ErrNil cannot be nil.
var ErrNil = errors.New("error value cannot be nil")

// Log saves the error and continues the program.
func Log(err error) {
	if err != nil {
		// save error to log file
		if err = save(err, ""); err != nil {
			log.Fatalf("%s %s", color.Danger.Sprint("!"), err)
		}
	}
}

// LogFatal saves the error and exits.
func LogFatal(err error) {
	if err != nil {
		// save error to log file
		if err = save(err, ""); err != nil {
			log.Fatalf("%s %s", color.Danger.Sprint("!"), err)
		}
		// print error
		switch Panic {
		case true:
			log.Println(fmt.Sprintf("error type: %T\tmsg: %v", err, err))
			log.Panic(err)
		default:
			Fatal("log error", "", err)
		}
	}
}

// Path is the absolute path and filename of the error log file.
func Path() string {
	fp, err := gap.NewScope(gap.User, "df2").LogPath(filename)
	if err != nil {
		h, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(fmt.Errorf("log path userhomedir: %w", err))
		}
		return path.Join(h, filename)
	}
	return fp
}

// save an error to the log directory, an optional named file is available for unit tests.
func save(err error, name string) error {
	if err == nil || fmt.Sprintf("%v", err) == "" {
		return fmt.Errorf("logs save: %w", ErrNil)
	}
	// use UTC date and times in the log file
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	if name == "" {
		name = Path()
	}
	p := filepath.Dir(name)
	if _, e := os.Stat(p); os.IsNotExist(e) {
		if e := os.MkdirAll(p, dirmode); e != nil {
			return e
		}
	}
	file, e := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filemode)
	if e != nil {
		return e
	}
	defer file.Close()
	log.SetOutput(file)
	log.Print(err)
	log.SetOutput(os.Stderr)
	return file.Close()
}
