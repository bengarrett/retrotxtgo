//Package filesystem to handle the opening and reading of text files
package filesystem

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
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
		if _, err := fmt.Fprintf(os.Stderr, "failed to clean: %q: %s", name, err); err != nil {
			logs.LogFatal(fmt.Errorf("failed to print to stderr and to clean: %s", name))
		}
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
		if runtime.GOOS == "windows" && dir == "" {
			dir = p + "\\"
		} else {
			dir = filepath.Join(dir, p)
		}
	}
	if root() {
		dir = string(os.PathSeparator) + dir
	}
	return dir
}

// Save bytes to a named file location.
func Save(name string, b ...byte) (path string, err error) {
	path, err = dir(name)
	if err != nil {
		return path, fmt.Errorf("save could not open directory: %q: %s", name, err)
	}
	path = name
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filemode)
	if err != nil {
		return path, fmt.Errorf("save could not open file: %q: %s", path, err)
	}
	defer file.Close()
	// bufio is the most performant
	writer := bufio.NewWriter(file)
	for _, c := range b {
		if err = writer.WriteByte(c); err != nil {
			return path, fmt.Errorf("save could not write bytes: %s", err)
		}
	}
	if err = writer.Flush(); err != nil {
		return path, fmt.Errorf("save could not flush the writer: %s", err)
	}
	//ioutil.WriteFile(filename,data,perm)
	path, err = filepath.Abs(file.Name())
	if err != nil {
		return path, fmt.Errorf("save could not find the absolute filename: %s", err)
	}
	return path, file.Close()
}

// SaveTemp saves bytes to a named temporary file.
func SaveTemp(filename string, b ...byte) (path string, err error) {
	path, err = Save(tempFile(filename), b...)
	if err != nil {
		return path, fmt.Errorf("could not save the temporary file: %s", err)
	}
	return path, nil
}

// Tar blah
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
	return f.Close()
}

// Touch creates an empty file at the named location.
func Touch(name string) (path string, err error) {
	path, err = Save(name, nil...)
	if err != nil {
		return path, fmt.Errorf("could not touch a new file: %s", err)
	}
	return path, nil
}

func dir(name string) (path string, err error) {
	path = filepath.Dir(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, dirmode); err != nil {
			return "", fmt.Errorf("dir could not make the directory: %s %s: %s", dirmode, path, err)
		}
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
