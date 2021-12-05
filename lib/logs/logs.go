// Package logs for the display of errors.
package logs

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/bengarrett/retrotxtgo/lib/internal/save"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
	gap "github.com/muesli/go-app-paths"
)

const (
	filename = "errors.log"

	Panic = false
)

// FatalSave saves the error to the logfile and exits.
func FatalSave(err error) {
	if err == nil {
		return
	}
	// save error to log file
	if err = SaveErr(err, ""); err != nil {
		log.Fatalf("%s %s", color.Danger.Sprint("!"), err)
	}
	// print error
	switch Panic {
	case true:
		log.Println(fmt.Sprintf("error type: %T\tmsg: %v", err, err))
		log.Panic(err)
	default:
		FatalWrap(ErrSave, err)
	}
}

// Save by appending the error to the logfile.
func Save(err error) {
	if err == nil {
		return
	}
	// save error to log file
	if err = SaveErr(err, ""); err != nil {
		log.Fatalf("%s %s", color.Danger.Sprint("!"), err)
	}
}

// LastEntry returns the most recent saved entry in the error log file.
func LastEntry() (s string, err error) {
	name := Name()
	file, err := os.Open(name)
	if err != nil {
		return "", fmt.Errorf("read tail could not open file: %q: %w", name, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		s = scanner.Text()
	}
	if err = scanner.Err(); err != nil {
		return "", fmt.Errorf("read tail could scan file bytes: %q: %w", name, err)
	}
	return s, file.Close()
}

// Name is the absolute path and filename of the error log file.
func Name() string {
	fp, err := gap.NewScope(gap.User, meta.Dir).LogPath(filename)
	if err == nil {
		return fp
	}
	h, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(fmt.Errorf("log path userhomedir: %w", err))
	}
	return path.Join(h, filename)
}

// Save an error to the log directory.
// An optional named file is available for unit tests.
func SaveErr(err error, name string) error {
	if err == nil || fmt.Sprintf("%v", err) == "" {
		return fmt.Errorf("logs save: %w", ErrErrorNil)
	}
	// use UTC date and times in the log file
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	if name == "" {
		name = Name()
	}
	p := filepath.Dir(name)
	if _, e := os.Stat(p); os.IsNotExist(e) {
		if e := os.MkdirAll(p, save.DirMode); e != nil {
			return e
		}
	}
	const appendFile = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	file, e := os.OpenFile(name, appendFile, save.LogFileMode)
	if e != nil {
		return e
	}
	defer file.Close()
	log.SetOutput(file)
	log.Print(err)
	log.SetOutput(os.Stderr)
	return file.Close()
}
