// Package filesystem handles the opening and reading of text files.
package filesystem

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
		p := ""
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
		if dir, cont = winDir(i, p, runtime.GOOS, dir); cont {
			continue
		}
		dir = filepath.Join(dir, p)
	}
	if root() {
		dir = string(os.PathSeparator) + dir
	}
	return dir
}

// winDir appends Windows style syntax to the directory.
func winDir(i int, p, platform, dir string) (s string, cont bool) {
	if platform == "windows" {
		if len(p) == 2 && p[1:] == ":" {
			dir = strings.ToUpper(p) + "\\"
			return dir, true
		}
		if dir == "" && i > 0 {
			dir = p + "\\"
			return dir, true
		}
	}
	return dir, false
}

// SaveTemp saves bytes to a named temporary file.
func SaveTemp(name string, b ...byte) (string, error) {
	_, path, err := save.Save(tempFile(name), b...)
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
		if err := addTar(file, w); err != nil {
			return err
		}
	}
	return nil
}

// AddTar inserts the named file to the TAR writer.
func addTar(name string, w *tar.Writer) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return err
	}
	h := &tar.Header{
		Name:    f.Name(),
		Size:    s.Size(),
		Mode:    int64(s.Mode()),
		ModTime: s.ModTime(),
	}
	if err1 := w.WriteHeader(h); err1 != nil {
		return err1
	}
	_, err = io.Copy(w, f)
	if err != nil {
		return nil
	}
	return f.Close()
}

// Touch creates an empty file at the named location.
func Touch(name string) (string, error) {
	_, path, err := save.Save(name, nil...)
	if err != nil {
		return path, fmt.Errorf("could not touch a new file: %w", err)
	}
	return path, nil
}

// tempFile returns a path to the named file
// if it was stored in the system's temporary directory.
func tempFile(name string) string {
	path := name
	if filepath.Base(name) == name {
		path = filepath.Join(os.TempDir(), name)
	}
	return path
}

// Write b to the named file.
func Write(name string, b ...byte) (written int, path string, err error) {
	return save.Save(name, b...)
}
