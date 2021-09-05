// Package filesystem handles the opening and reading of text files.
package filesystem

import (
	"archive/tar"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/logs"
)

const (
	fileMode os.FileMode = 0660
	dirMode  os.FileMode = 0700
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

// Save bytes to the named file location.
func Save(name string, b ...byte) (written int, path string, err error) {
	path, err = dir(name)
	if err != nil {
		return 0, path, fmt.Errorf("save could not open directory %q: %w", name, err)
	}
	path = name
	const overwrite = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(path, overwrite, fileMode)
	if err != nil {
		return 0, path, fmt.Errorf("save could not open file %q: %w", path, err)
	}
	defer file.Close()
	// bufio is the most performant
	writer := bufio.NewWriter(file)
	for i, c := range b {
		written = i
		if err = writer.WriteByte(c); err != nil {
			return 0, path, fmt.Errorf("save could not write bytes: %w", err)
		}
	}
	if err = writer.Flush(); err != nil {
		return 0, path, fmt.Errorf("save could not flush the writer: %w", err)
	}
	path, err = filepath.Abs(file.Name())
	if err != nil {
		return 0, path, fmt.Errorf("save could not find the absolute filename: %w", err)
	}
	return written, path, file.Close()
}

// SaveTemp saves bytes to a named temporary file.
func SaveTemp(name string, b ...byte) (string, error) {
	_, path, err := Save(tempFile(name), b...)
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
	_, path, err := Save(name, nil...)
	if err != nil {
		return path, fmt.Errorf("could not touch a new file: %w", err)
	}
	return path, nil
}

// dir creates the named path directory if it doesn't exist.
func dir(name string) (string, error) {
	path := filepath.Dir(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, dirMode); err != nil {
			return "", fmt.Errorf("dir could not make the directory: %s %s: %w", dirMode, path, err)
		}
	}
	return path, nil
}

// lineBreak returns line break character for the system platform.
func lineBreak(platform lineBreaks) string {
	switch platform {
	case dos, win:
		return fmt.Sprintf("%x%x", CarriageReturn, Linefeed)
	case c64, darwin, mac:
		return fmt.Sprintf("%x", CarriageReturn)
	case amiga, linux, unix:
		return fmt.Sprintf("%x", Linefeed)
	case nl: // use operating system default
		return "\n"
	}
	return "\n"
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
