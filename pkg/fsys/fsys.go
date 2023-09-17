// Package fsys handles the opening, reading and writing of files.
package fsys

import (
	"archive/tar"
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/pkg/fsys/internal/util"
	"github.com/bengarrett/retrotxtgo/pkg/internal/save"
)

var (
	ErrLB        = errors.New("linebreak runes cannot be empty")
	ErrMax       = errors.New("maximum attempts reached")
	ErrName      = errors.New("name file cannot be a directory")
	ErrNotFound  = errors.New("cannot find the file or sample file")
	ErrPipeEmpty = errors.New("empty text stream from piped stdin (standard input)")
	ErrReader    = errors.New("the r reader cannot be nil")
)

// Clean removes the named file or directory.
// It is intended to be used as a helper in unit tests.
func Clean(name string) {
	if err := os.RemoveAll(name); err != nil {
		fmt.Fprintf(os.Stderr, "failed to clean %q: %s", name, err)
	}
}

// SaveTemp saves bytes to a named temporary file.
// The path to the file is returned.
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
