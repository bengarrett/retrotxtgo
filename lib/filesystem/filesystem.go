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
	// posix permission bits for files.
	filemode os.FileMode = 0660
	// posix permission bits for directories.
	dirmode os.FileMode = 0700
)

// ErrStdErr could not print to stderr.
var ErrStdErr = errors.New("failed to print to stderr")

// Clean removes the named file or directory.
func Clean(name string) {
	if err := os.RemoveAll(name); err != nil {
		if _, err = fmt.Fprintf(os.Stderr, "failed to clean %q: %s", name, err); err != nil {
			logs.SaveFatal(fmt.Errorf("clean %s: %w", name, ErrStdErr))
		}
	}
}

// DirExpansion traverses the named directory to apply shell-like expansions.
// It currently supports limited Bash tilde, shell dot and double dot syntax.
func DirExpansion(name string) (dir string) {
	if name == "" {
		return dir
	}
	const homeDir, currentDir, parentDir = "~", ".", ".."
	var root = func() bool {
		return name[:1] == string(os.PathSeparator)
	}

	var err error
	// Bash tilde expension http://www.gnu.org/software/bash/manual/html_node/Tilde-Expansion.html
	paths := strings.Split(name, string(os.PathSeparator))
	for i, s := range paths {
		p := ""
		switch s {
		case homeDir:
			if p, err = os.UserHomeDir(); err != nil {
				logs.SaveFatal(err)
			}
		case currentDir:
			if i != 0 {
				continue
			}
			if p, err = os.Getwd(); err != nil {
				logs.SaveFatal(err)
			}
		case parentDir:
			if i != 0 {
				dir = filepath.Dir(dir)
				continue
			}
			wd, err := os.Getwd()
			if err != nil {
				logs.SaveFatal(err)
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

// WinDir appends Windows style syntax to the directory.
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

// Save bytes to a named file location.
func Save(name string, b ...byte) (nn int, path string, err error) {
	path, err = dir(name)
	if err != nil {
		return nn, path, fmt.Errorf("save could not open directory %q: %w", name, err)
	}
	path = name
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filemode)
	if err != nil {
		return nn, path, fmt.Errorf("save could not open file %q: %w", path, err)
	}
	defer file.Close()
	// bufio is the most performant
	writer := bufio.NewWriter(file)
	for i, c := range b {
		nn = i
		if err = writer.WriteByte(c); err != nil {
			return nn, path, fmt.Errorf("save could not write bytes: %w", err)
		}
	}
	if err = writer.Flush(); err != nil {
		return nn, path, fmt.Errorf("save could not flush the writer: %w", err)
	}
	path, err = filepath.Abs(file.Name())
	if err != nil {
		return nn, path, fmt.Errorf("save could not find the absolute filename: %w", err)
	}
	return nn, path, file.Close()
}

// SaveTemp saves bytes to a named temporary file.
func SaveTemp(filename string, b ...byte) (path string, err error) {
	_, path, err = Save(tempFile(filename), b...)
	if err != nil {
		return path, fmt.Errorf("could not save the temporary file: %w", err)
	}
	return path, nil
}

// T are short strings shared between various app tests.
func T() map[string]string {
	return map[string]string{
		// Newline sample using yjr operating system defaults
		"Newline": "a\nb\nc...\n",
		// Symbols for Unicode Wingdings
		"Symbols": `[☠|☮|♺]`,
		// Tabs and Unicode glyphs
		"Tabs": "☠\tSkull and crossbones\n\n☮\tPeace symbol\n\n♺\tRecycling",
		// Escapes and control codes.
		"Escapes": "bell:\a,back:\b,tab:\t,form:\f,vertical:\v,quote:\"",
		// Digits in various formats
		"Digits": "\xb0\260\u0170\U00000170",
	}
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
func Touch(name string) (path string, err error) {
	_, path, err = Save(name, nil...)
	if err != nil {
		return path, fmt.Errorf("could not touch a new file: %w", err)
	}
	return path, nil
}

func dir(name string) (path string, err error) {
	path = filepath.Dir(name)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, dirmode); err != nil {
			return "", fmt.Errorf("dir could not make the directory: %s %s: %w", dirmode, path, err)
		}
	}
	return path, err
}

// LineBreak returns line break character for the system platform.
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

// TempFile returns a path to the named file
// if it was stored in the system's temporary directory.
func tempFile(name string) (path string) {
	path = name
	if filepath.Base(name) == name {
		path = filepath.Join(os.TempDir(), name)
	}
	return path
}
