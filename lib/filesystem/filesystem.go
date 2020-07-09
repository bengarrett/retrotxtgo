//Package filesystem to handle the opening and reading of text files
package filesystem

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"retrotxt.com/retrotxt/lib/logs"
)

// T are short strings shared between various app tests.
var T = map[string]string{
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

const (
	cr = "\x0d"
	lf = "\x0a"
	// posix permission bits for files
	filemode os.FileMode = 0660
	// posix permission bits for directories
	dirmode os.FileMode = 0700
)

// Clean removes the named file or directory.
func Clean(name string) {
	if err := os.RemoveAll(name); err != nil {
		fmt.Fprintln(os.Stderr, "removing path:", err)
	}
}

// DirExpansion traverses the named directory to apply shell-like expansions.
// It currently supports limited Bash tilde, shell dot and double dot syntax.
func DirExpansion(name string) (dir string) {
	if name == "" {
		return dir
	}
	var root = func() bool {
		if name[:1] == string(os.PathSeparator) {
			return true
		}
		return false
	}
	var err error
	// Bash tilde expension http://www.gnu.org/software/bash/manual/html_node/Tilde-Expansion.html
	paths := strings.Split(name, string(os.PathSeparator))
	for i, s := range paths {
		p := ""
		switch s {
		case "~":
			p, err = os.UserHomeDir()
			if err != nil {
				logs.LogFatal(err)
			}
		case ".":
			if i != 0 {
				continue
			}
			p, err = os.Getwd()
			if err != nil {
				logs.LogFatal(err)
			}
		case "..":
			if i == 0 {
				wd, err := os.Getwd()
				if err != nil {
					logs.LogFatal(err)
				}
				p = filepath.Dir(wd)
			} else {
				dir = filepath.Dir(dir)
				continue
			}
		default:
			p = s
		}
		dir = filepath.Join(dir, p)
	}
	if root() {
		dir = string(os.PathSeparator) + dir
	}
	return dir
}

// Save bytes to a named file location.
func Save(filename string, b []byte) (path string, err error) {
	path, err = dir(filename)
	if err != nil {
		return path, err
	}
	path = filename
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filemode)
	if err != nil {
		return path, err
	}
	// bufio is the most performant
	writer := bufio.NewWriter(file)
	for _, c := range b {
		if err = writer.WriteByte(c); err != nil {
			return path, err
		}
	}
	if err = writer.Flush(); err != nil {
		return path, err
	}
	if err = file.Close(); err != nil {
		return path, err
	}
	//ioutil.WriteFile(filename,data,perm)
	return filepath.Abs(file.Name())
}

// SaveTemp saves bytes to a named temporary file.
func SaveTemp(filename string, b []byte) (path string, err error) {
	return Save(tempFile(filename), b)
}

// Tar blah
func Tar(name string, files []string) error {
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
	if err = w.WriteHeader(h); err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	if err != nil {
		return nil
	}

	return nil
}

// Touch creates an empty file at the named location.
func Touch(name string) (path string, err error) {
	return Save(name, nil)
}

func dir(name string) (path string, err error) {
	path = filepath.Dir(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, dirmode)
	}
	return path, err
}

// nl returns a platform's newline character.
func nl(platform string) string {
	switch platform {
	case "dos", "windows":
		return cr + lf
	case "c64", "darwin", "mac":
		return cr
	case "amiga", "linux", "unix":
		return lf
	default: // use operating system default
		return "\n"
	}
}

func tempFile(name string) (path string) {
	path = name
	if filepath.Base(name) == name {
		path = filepath.Join(os.TempDir(), name)
	}
	return path
}
