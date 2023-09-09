package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Meta data to embed into the HTML.
type Meta struct {
	Author struct {
		Flag  bool
		Value string
	}
	ColorScheme struct {
		Flag  bool
		Value string
	}
	Description struct {
		Flag  bool
		Value string
	}
	Keywords struct {
		Flag  bool
		Value string
	}
	Referrer struct {
		Flag  bool
		Value string
	}
	Robots struct {
		Flag  bool
		Value string
	}
	ThemeColor struct {
		Flag  bool
		Value string
	}
	Generator   bool
	NoTranslate bool
	RetroTxt    bool
}

// Destination determines if user supplied arguments are a valid file or directory destination.
func Destination(args ...string) (string, error) {
	if len(args) == 0 {
		return "", nil
	}
	dir := filepath.Clean(strings.Join(args, " "))
	if len(dir) == 1 {
		return dirs(dir)
	}
	part := strings.Split(dir, string(os.PathSeparator))
	if len(part) > 1 {
		var err error
		part[0], err = dirs(part[0])
		if err != nil {
			return "", fmt.Errorf("destination arguments: %w", err)
		}
	}
	return strings.Join(part, string(os.PathSeparator)), nil
}

// dirs parses and expand special directory characters.
func dirs(dir string) (string, error) {
	const (
		homeDir    = "~"
		currentDir = "."
	)
	s := ""
	var err error
	switch dir {
	case homeDir:
		s, err = os.UserHomeDir()
	case currentDir:
		s, err = os.Getwd()
	case string(os.PathSeparator):
		s, err = filepath.Abs(dir)
	}
	if err != nil {
		return "", fmt.Errorf("parse directory error: %q: %w", dir, err)
	}
	return s, nil
}
