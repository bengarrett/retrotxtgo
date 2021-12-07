// Package filesystem handles the opening and reading of text files.
package filesystem

import (
	"archive/tar"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/filesystem/internal/util"
	"github.com/bengarrett/retrotxtgo/lib/internal/save"
	"github.com/bengarrett/retrotxtgo/lib/logs"
)

// ErrStdErr could not print to stderr.
var (
	ErrNotFound = errors.New("cannot find the file or sample file")
	ErrStdErr   = errors.New("failed to print to stderr")
)

// Clean removes the named file or directory.
func Clean(name string) {
	if err := os.RemoveAll(name); err != nil {
		if _, err = fmt.Fprintf(os.Stderr, "failed to clean %q: %s", name, err); err != nil {
			logs.FatalSave(fmt.Errorf("clean %s: %w", name, ErrStdErr))
		}
	}
}

// DirExpansion returns the absolute directory path from a named path using shell-like expansions.
// It currently supports limited Bash tilde, shell dot and double dot syntax.
func DirExpansion(name string) string {
	if name == "" {
		return ""
	}
	root := func() bool {
		return name[:1] == string(os.PathSeparator)
	}
	const homeDir, currentDir, parentDir = "~", ".", ".."
	var err error
	// Bash tilde expension http://www.gnu.org/software/bash/manual/html_node/Tilde-Expansion.html
	dir, paths := "", strings.Split(name, string(os.PathSeparator))
	for i, s := range paths {
		var p string
		switch s {
		case homeDir:
			if p, err = os.UserHomeDir(); err != nil {
				logs.FatalSave(err)
			}
		case currentDir:
			if i != 0 {
				continue
			}
			if p, err = os.Getwd(); err != nil {
				logs.FatalSave(err)
			}
		case parentDir:
			if i != 0 {
				dir = filepath.Dir(dir)
				continue
			}
			wd, err := os.Getwd()
			if err != nil {
				logs.FatalSave(err)
			}
			p = filepath.Dir(wd)
		default:
			p = s
		}
		var cont bool
		if dir, cont = util.Windows(i, p, runtime.GOOS, dir); cont {
			continue
		}
		dir = filepath.Join(dir, p)
	}
	if root() {
		dir = string(os.PathSeparator) + dir
	}
	return dir
}

// SaveTemp saves bytes to a named temporary file.
func SaveTemp(name string, b ...byte) (string, error) {
	_, path, err := save.Save(util.Temp(name), b...)
	if err != nil {
		return path, fmt.Errorf("could not save the temporary file: %w", err)
	}
	return path, nil
}

// Tar add files to a named tar file archive.
func Tar(name string, files ...string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	w := tar.NewWriter(f)
	defer w.Close()
	for _, file := range files {
		if err := util.AddTar(file, w); err != nil {
			return err
		}
	}
	return nil
}

// Touch creates an empty file at the named location.
func Touch(name string) (string, error) {
	_, path, err := save.Save(name, nil...)
	if err != nil {
		return path, fmt.Errorf("could not touch a new file: %w", err)
	}
	return path, nil
}

// Write b to the named file.
func Write(name string, b ...byte) (written int, path string, err error) {
	return save.Save(name, b...)
}
