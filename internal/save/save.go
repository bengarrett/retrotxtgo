package save

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DirMode     os.FileMode = 0o700
	FileMode    os.FileMode = 0o660
	LogFileMode os.FileMode = 0o600
)

// Save bytes to the named file location.
// The return values are the number of bytes written, the absolute path and filename, and any error.
func Save(name string, b ...byte) (int, string, error) {
	if _, err := dir(name); err != nil {
		return 0, "", fmt.Errorf("save could not open directory %q: %w", name, err)
	}
	path := name
	const overwrite = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(path, overwrite, FileMode)
	if err != nil {
		return 0, path, fmt.Errorf("save could not open file %q: %w", path, err)
	}
	defer file.Close()
	// bufio is the most performant
	w := bufio.NewWriter(file)
	written := 0
	for i, c := range b {
		written = i
		if err := w.WriteByte(c); err != nil {
			return 0, path, fmt.Errorf("save could not write bytes: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return 0, path, fmt.Errorf("save could not flush the writer: %w", err)
	}
	path, err = filepath.Abs(file.Name())
	if err != nil {
		return 0, path, fmt.Errorf("save could not find the absolute filename: %w", err)
	}
	return written, path, nil
}

// dir creates the named path directory if it doesn't exist.
func dir(name string) (string, error) {
	path := filepath.Dir(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, DirMode); err != nil {
			return "", fmt.Errorf("dir could not make the directory: %s %s: %w", DirMode, path, err)
		}
	}
	return path, nil
}
