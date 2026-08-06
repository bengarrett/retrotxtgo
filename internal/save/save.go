package save

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DirMode     os.FileMode = 0o700
	FileMode    os.FileMode = 0o660
	LogFileMode os.FileMode = 0o600
)

var ErrEmpty = errors.New("cannot be empty")

// Save writes bytes to the specified filename.
// It returns the number of bytes written, the absolute path of the file, and any error.
func Save(name string, b ...byte) (written int, path string, err error) {
	const format = "save %s %s: %w"
	if name == "" {
		return 0, "", fmt.Errorf(format, "named path", "", ErrEmpty)
	}
	path, err = filepath.Abs(name)
	if err != nil {
		return 0, "", fmt.Errorf(format, "filepath abs", name, err)
	}

	const perm = 0o755
	if err := os.MkdirAll(filepath.Dir(path), perm); err != nil {
		return 0, path, fmt.Errorf(format, "mkdirall", name, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return 0, path, fmt.Errorf(format, "create", name, err)
	}

	defer func() {
		if cErr := file.Close(); cErr != nil && err == nil {
			err = fmt.Errorf(format, "close", name, cErr)
		}
	}()

	written, err = file.Write(b)
	if err != nil {
		return written, path, fmt.Errorf(format, "write", name, err)
	}
	return written, path, nil
}
