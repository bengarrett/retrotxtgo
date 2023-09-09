package util

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrNoWriter = errors.New("the tar.writer cannot be nil")
)

// InsertTar inserts the named file to the TAR writer.
func InsertTar(w *tar.Writer, name string) error {
	if w == nil {
		return ErrNoWriter
	}
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
		return nil //nolint:nilerr
	}
	return f.Close()
}

// Temp returns a path to the named file
// if it was stored in the system's temporary directory.
func Temp(name string) string {
	path := name
	if filepath.Base(name) == name {
		path = filepath.Join(os.TempDir(), name)
	}
	return path
}

// Windows appends Windows style syntax to the directory.
func Windows(i int, p, platform, dir string) (string, bool) {
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
