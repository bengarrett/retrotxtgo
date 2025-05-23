// Package fsys handles the opening, reading and writing of files.
package fsys

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bengarrett/retrotxtgo/internal/save"
)

var (
	ErrLB        = errors.New("linebreak runes cannot be empty")
	ErrMax       = errors.New("maximum attempts reached")
	ErrName      = errors.New("name file cannot be a directory")
	ErrNotFound  = errors.New("cannot find the file or sample file")
	ErrPipeEmpty = errors.New("empty text stream from piped stdin (standard input)")
	ErrReader    = errors.New("the r reader cannot be nil")
	ErrWriter    = errors.New("the w writer cannot be nil")
)

// SaveTemp saves bytes to a named temporary file.
// The path to the file is returned.
func SaveTemp(name string, b ...byte) (string, error) {
	_, path, err := save.Save(temp(name), b...)
	if err != nil {
		return path, fmt.Errorf("could not save the temporary file: %w", err)
	}
	return path, nil
}

// temp returns a path to the named file
// if it was stored in the system's temporary directory.
func temp(name string) string {
	path := name
	if filepath.Base(name) == name {
		path = filepath.Join(os.TempDir(), name)
	}
	return path
}

// Tar add files to a named tar file archive.
func Tar(name string, files ...string) error {
	w, err := os.Create(name)
	if err != nil {
		return err
	}
	defer w.Close()
	dst := tar.NewWriter(w)
	defer dst.Close()
	for _, file := range files {
		if err := InsertTar(dst, file); err != nil {
			return err
		}
	}
	return nil
}

// InsertTar inserts the named file to the TAR writer.
func InsertTar(dst *tar.Writer, name string) error {
	if dst == nil {
		return ErrWriter
	}
	src, err := os.Open(name)
	if err != nil {
		return err
	}
	defer src.Close()
	s, err := src.Stat()
	if err != nil {
		return err
	}
	h := &tar.Header{
		Name:    src.Name(),
		Size:    s.Size(),
		Mode:    int64(s.Mode()),
		ModTime: s.ModTime(),
	}
	if err := dst.WriteHeader(h); err != nil {
		return err
	}
	const size = 4 * 1024
	buf := make([]byte, size)
	if _, err = io.CopyBuffer(dst, src, buf); err != nil {
		return err
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
