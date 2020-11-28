// Package create makes HTML and other web resources from a text file.
package create

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"retrotxt.com/retrotxt/internal/pack"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"

	"github.com/spf13/viper"
)

/*
arguments not working.
--controls

--swap-chars=124?

--compress # maybe this should be used as an alt to --save?

todo:
fails!: create static/text/cp437-crlf.txt # no codepage and eof
fails!: create static/text/cp437-crlf.txt --encode=cp437 # no eof
add eof bool flag (default on) ?
*/

// Args holds arguments and options sourced from user flags or the config file.
type Args struct {
	Source struct {
		Encoding   encoding.Encoding // Original encoding of the text source
		HiddenBody string            // Pre-text content override, accessible by a hidden flag
		Name       string            // Text source, usually a file or pack name
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
	layout    Layout // layout flag interpretation
	Port      uint   // Port for HTTP server
	FontEmbed bool
	test      bool   // unit test mode
	Layout    string // Layout of the HTML
	Syntax    string // Syntax and color theming printing HTML
	tmpl      string // template filename
	pack      string // template package name
}

// Meta data embedded into the webpage.
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

// PageData temporarily holds template data used for the HTML layout.
type PageData struct {
	ExternalEmbed   bool
	FontEmbed       bool
	MetaGenerator   bool
	MetaNoTranslate bool
	MetaRetroTxt    bool
	BuildVersion    string
	BuildDate       string
	CacheRefresh    string
	Comment         string
	FontFamily      string
	MetaAuthor      string
	MetaColorScheme string
	MetaDesc        string
	MetaKeywords    string
	MetaReferrer    string
	MetaRobots      string
	MetaThemeColor  string
	PageTitle       string
	PreText         string
	CSSEmbed        template.CSS
	ScriptEmbed     template.JS
}

// Layout are HTML template variations.
type Layout int

const (
	// use 0 as an error placeholder.
	_ Layout = iota
	// Standard template with external CSS, JS, fonts.
	Standard
	// Inline template with CSS and JS embedded.
	Inline
	// Compact template with external CSS, JS, fonts and no meta-tags.
	Compact
	// None template, just print the generated HTML.
	None
)

func (l Layout) String() string {
	return [...]string{unknown, standard, inline, compact, none}[l]
}

const (
	none     = "none"
	compact  = "compact"
	inline   = "inline"
	standard = "standard"
	unknown  = "unknown"
)

var (
	// ErrName unknown font.
	ErrName = errors.New("font name is not known")
	// ErrPack font not found.
	ErrPack = errors.New("font pack is not found")
	// ErrEmptyName filename is empty.
	ErrEmptyName = errors.New("filename is empty")
	// ErrReqOW require overwrite flag.
	ErrReqOW = errors.New("include an -o flag to overwrite")
	// ErrPackGet invalid pack name.
	ErrPackGet = errors.New("pack.get name is invalid")
	// ErrUnknownFF unknown font family.
	ErrUnknownFF = errors.New("unknown font family")
	// ErrNilByte nil byte value.
	ErrNilByte = errors.New("cannot convert a nil byte value")
	// ErrTmplDir temp file is a dir.
	ErrTmplDir = errors.New("the path to the template file is a directory")
	// ErrNoLayout layout missing.
	ErrNoLayout = errors.New("layout does not exist")
)

// ColorScheme values for the content attribute of <meta name="color-scheme">.
func ColorScheme() [3]string {
	return [...]string{"normal", "dark light", "only light"}
}

// Referrer values for the content attribute of <meta name="referrer">.
func Referrer() [8]string {
	return [...]string{"no-referrer", "origin", "no-referrer-when-downgrade",
		"origin-when-cross-origin", "same-origin", "strict-origin", "strict-origin-when-cross-origin", "unsafe-URL"}
}

// Robots values for the content attribute of <meta name="robots">.
func Robots() [9]string {
	return [...]string{"index", "noindex", "follow", "nofollow", none, "noarchive", "nosnippet", "noimageindex", "nocache"}
}

// Layouts are the names of the HTML templates.
func Layouts() []string {
	return []string{Standard.String(), Inline.String(), Compact.String(), None.String()}
}

// Pack is the packed name of the HTML template.
func (l Layout) Pack() string {
	return [...]string{unknown, standard, standard, standard, none}[l]
}

// Create handles the target output command arguments.
func (args *Args) Create(b *[]byte) {
	var err error
	args.layout = layout(args.Layout)
	switch {
	case args.Save.Compress:
		fmt.Println("TODO!")
	case args.Save.AsFiles:
		// use config save directory
		// otherwise assume Destination path is a temporary --serve location
		if args.Save.Destination == "" {
			dir := []string{viper.GetString("save-directory")}
			if args.Save.Destination, err = destination(dir...); err != nil {
				logs.Fatal("save to directory failure", fmt.Sprintf("%s", dir), err)
			}
		}
		ch := make(chan error)
		go args.saveCSS(ch)
		go args.saveFont(ch)
		go args.saveHTML(b, ch)
		go args.saveJS(ch)
		go args.saveFavIcon(ch)
		err1, err2, err3, err4, err5 := <-ch, <-ch, <-ch, <-ch, <-ch
		const errS, errCode = "could not save file", 1
		if err1 != nil {
			logs.Println(errS, "", err1)
			os.Exit(errCode)
		}
		if err2 != nil {
			logs.Println(errS, "", err2)
			os.Exit(errCode)
		}
		if err3 != nil {
			logs.Println(errS, "", err3)
			os.Exit(errCode)
		}
		if err4 != nil {
			logs.Println(errS, "", err4)
			os.Exit(errCode)
		}
		if err5 != nil {
			logs.Println(errS, "", err5)
			os.Exit(errCode)
		}
	default:
		// print to terminal
		if err = args.Stdout(b); err != nil {
			logs.Fatal("print to stdout", "", err)
		}
	}
}

// Stdout creates and prints the html template.
func (args *Args) Stdout(b *[]byte) error {
	// html
	buf, err := args.marshalTextTransform(b)
	if err != nil {
		return fmt.Errorf("stdout: %w", err)
	}
	// js
	js := pack.Get("js/scripts.js")
	// css
	css := pack.Get("css/styles.css")
	// font
	ff := args.FontFamily.Value
	f := Family(ff).String()
	if f == "" {
		return fmt.Errorf("create.saveFontCSS %q: %w", ff, ErrUnknownFF)
	}
	font, err := FontCSS(f, args.FontEmbed)
	if err != nil {
		return err
	}
	const (
		fJS   = "\nJS file: %s\n"
		fCSS  = "\nCSS file: %s\n"
		fFont = "\nFont %q file: %s\n"
		fHTML = "\nHTML file: %s\n"
	)
	var noSyntax = func() {
		fmt.Printf(fJS, nameJS)
		fmt.Println(string(js))
		fmt.Printf(fCSS, nameCSS)
		fmt.Println(string(css))
		fmt.Printf(fFont, f, nameFont)
		fmt.Println(string(font))
		fmt.Printf(fHTML, nameHTML)
		fmt.Println(buf.String())
	}
	// stdout
	const format = "%s\n\n%s"
	switch args.Syntax {
	case "", none:
		noSyntax()
	default:
		if !str.Valid(args.Syntax) {
			fmt.Printf("unknown style %q, so using none\n", args.Syntax)
			noSyntax()
			return nil
		}
		fmt.Printf(fJS, nameJS)
		if err = str.Highlight(string(js), "js", args.Syntax, true); err != nil {
			return fmt.Errorf("stdout js highlight: %w", err)
		}
		fmt.Printf(fCSS, nameCSS)
		if err = str.Highlight(string(css), "css", args.Syntax, true); err != nil {
			return fmt.Errorf("stdout css highlight: %w", err)
		}
		fmt.Printf(fFont, f, nameFont)
		if err = str.Highlight(string(font), "css", args.Syntax, true); err != nil {
			return fmt.Errorf("stdout font css highlight: %w", err)
		}
		fmt.Printf(fHTML, nameHTML)
		if err = str.Highlight(buf.String(), "html", args.Syntax, true); err != nil {
			return fmt.Errorf("stdout html highlight: %w", err)
		}
	}
	return nil
}

// Normalize runes to bytes by making adjustments to text control codes.
func Normalize(e encoding.Encoding, r ...rune) (b []byte) {
	s := ""
	switch e {
	case charmap.CodePage037, charmap.CodePage1047, charmap.CodePage1140:
		var err error
		s, _, err = transform.String(replaceNELs(), string(r))
		if err != nil {
			s = string(r)
		}
	default:
		s = string(r)
	}
	return []byte(s)
}

// destination determines if user supplied arguments are a valid file or directory destination.
func destination(args ...string) (path string, err error) {
	if len(args) == 0 {
		return path, nil
	}
	dir := strings.Join(args, " ")
	dir = filepath.Clean(dir)
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

func dirs(dir string) (path string, err error) {
	switch dir {
	case "~":
		return os.UserHomeDir()
	case ".":
		return os.Getwd()
	case "\\", "/":
		return filepath.Abs(dir)
	}
	if err != nil {
		return "", fmt.Errorf("parse directory error: %q: %w", dir, err)
	}
	return "", nil
}

func layout(name string) (l Layout) {
	switch name {
	case standard, "s":
		return Standard
	case inline, "i":
		return Inline
	case compact, "c":
		return Compact
	case none, "n":
		return None
	}
	return l
}

// replace EBCDIC newline with linefeed.
func replaceNELs() runes.Transformer {
	return runes.Map(func(r rune) rune {
		if r == filesystem.NextLine {
			return filesystem.Linefeed
		}
		return r
	})
}
