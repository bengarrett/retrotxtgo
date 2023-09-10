// Package fsys handles the opening, reading and writing of files.
package fsys

import (
	"archive/tar"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bengarrett/retrotxtgo/pkg/fsys/internal/util"
	"github.com/bengarrett/retrotxtgo/pkg/internal/save"
)

// ErrStd could not print to stderr.
var (
	ErrNotFound = errors.New("cannot find the file or sample file")
)

// Clean removes the named file or directory.
func Clean(name string) {
	if err := os.RemoveAll(name); err != nil {
		fmt.Fprintf(os.Stderr, "failed to clean %q: %s", name, err)
	}
}

// DirExpansion returns the absolute directory path from a named path using shell-like expansions.
// It currently supports limited Bash tilde, shell dot and double dot syntax.
func DirExpansion(name string) (string, error) {
	if name == "" {
		return "", nil
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
				return "", err
			}
		case currentDir:
			if i != 0 {
				continue
			}
			if p, err = os.Getwd(); err != nil {
				return "", err
			}
		case parentDir:
			if i != 0 {
				dir = filepath.Dir(dir)
				continue
			}
			wd, err := os.Getwd()
			if err != nil {
				return "", err
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
	return dir, nil
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
		if err := util.InsertTar(w, file); err != nil {
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
// The number of bytes written and the path to the file are returned.
func Write(name string, b ...byte) (int, string, error) {
	return save.Save(name, b...)
}
