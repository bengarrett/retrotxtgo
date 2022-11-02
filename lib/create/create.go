// Package create makes HTML and other web resources from a text file.
package create

import (
	"github.com/bengarrett/bbs"
	"github.com/bengarrett/retrotxtgo/lib/create/internal/assets"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
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
	Metadata  assets.Meta
	SauceData SAUCE
	Port      uint   // Port for HTTP server
	FontEmbed bool   // embed the font as Base64 data
	Test      bool   // unit test mode
	Layout    string // Layout of the HTML
	Syntax    string // Syntax and color theming printing HTML
	// internals
	Tmpl string // template filename
	Pack string // template package name
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

// ColorScheme values for the content attribute of <meta name="color-scheme">.
func ColorScheme() [3]string {
	return [...]string{"normal", "dark light", "only light"}
}

// Referrer values for the content attribute of <meta name="referrer">.
func Referrer() [8]string {
	return [...]string{
		"no-referrer", "origin", "no-referrer-when-downgrade",
		"origin-when-cross-origin", "same-origin", "strict-origin",
		"strict-origin-when-cross-origin", "unsafe-URL",
	}
}

// Robots values for the content attribute of <meta name="robots">.
func Robots() [9]string {
	return [...]string{
		"index", "noindex", "follow", "nofollow", "none",
		"noarchive", "nosnippet", "noimageindex", "nocache",
	}
}

// Normalize runes into bytes by making adjustments to text control codes.
func Normalize(e encoding.Encoding, r ...rune) []byte {
	switch e {
	case charmap.CodePage037, charmap.CodePage1047, charmap.CodePage1140:
		s, _, err := transform.String(replaceNELs(), string(r))
		if err != nil {
			return []byte(string(r))
		}
		return []byte(s)
	}
	return []byte(string(r))
}

// replaceNELs replace EBCDIC newlines with Unicode linefeeds.
func replaceNELs() runes.Transformer {
	return runes.Map(func(r rune) rune {
		if r == filesystem.NextLine {
			return filesystem.Linefeed
		}
		return r
	})
}
