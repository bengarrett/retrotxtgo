package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bengarrett/bbs"
	"github.com/bengarrett/retrotxtgo/lib/create/internal/layout"
	"golang.org/x/text/encoding"
)

// Args holds arguments and options sourced from user flags and the config file.
type Args struct {
	Source struct {
		Encoding   encoding.Encoding // Original encoding of the text source
		HiddenBody string            // Pre-text content override, accessible by a hidden flag
		Name       string            // Text source, usually a file or pack name
		BBSType    bbs.BBS           // Optional BBS or ANSI text format
	}
	Save struct {
		AsFiles     bool   // Save assets as files
		Cache       bool   // Cache, when false will always unpack a new .gohtml template
		Compress    bool   // Compress and store all assets into an archive
		OW          bool   // OW overwrite any existing files when saving
		Destination string // Destination HTML destination either a directory or file
	}
	Title struct {
		Flag  bool
		Value string
	}
	FontFamily struct {
		Flag  bool
		Value string
	}
	Metadata  Meta
	SauceData struct {
		Use         bool
		Title       string
		Author      string
		Group       string
		Description string
		Width       uint
		Lines       uint
	}
	Port      uint   // Port for HTTP server
	FontEmbed bool   // embed the font as Base64 data
	Test      bool   // unit test mode
	Layout    string // Layout of the HTML
	Syntax    string // Syntax and color theming printing HTML
	// internals
	Layouts layout.Layout // layout flag interpretation
	Tmpl    string        // template filename
	Pack    string        // template package name
}

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
func Destination(args ...string) (path string, err error) {
	if len(args) == 0 {
		return path, nil
	}
	dir := filepath.Clean(strings.Join(args, " "))
	if len(dir) == 1 {
		return dirs(dir)
	}
	part := strings.Split(dir, string(os.PathSeparator))
	if len(part) > 1 {
		part[0], err = dirs(part[0])
		if err != nil {
			return path, fmt.Errorf("destination arguments: %w", err)
		}
	}
	return strings.Join(part, string(os.PathSeparator)), nil
}

// dirs parses and expand special directory characters.
func dirs(dir string) (path string, err error) {
	const (
		homeDir    = "~"
		currentDir = "."
	)
	switch dir {
	case homeDir:
		return os.UserHomeDir()
	case currentDir:
		return os.Getwd()
	case string(os.PathSeparator):
		return filepath.Abs(dir)
	}
	if err != nil {
		return "", fmt.Errorf("parse directory error: %q: %w", dir, err)
	}
	return "", nil
}
