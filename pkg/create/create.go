// Package create makes HTML and other web resources from a text file.
package create

import (
	"github.com/bengarrett/bbs"
	"github.com/bengarrett/retrotxtgo/pkg/create/internal/assets"
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
	FontFamily struct {
		Flag  bool
		Value string
	}
	Metadata  assets.Meta
	SauceData SAUCE
	Syntax    string // Syntax and color theming printing HTML
}

type SAUCE struct {
	Use         bool
	Title       string
	Author      string
	Group       string
	Description string
	Width       uint
	Lines       uint
}
