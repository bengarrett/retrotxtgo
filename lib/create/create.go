// Package create makes HTML and other web resources from a text file.
package create

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bengarrett/bbs"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/static"
	"github.com/spf13/viper"
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
	layout    Layout // layout flag interpretation
	Port      uint   // Port for HTTP server
	FontEmbed bool   // embed the font as Base64 data
	Test      bool   // unit test mode
	Layout    string // Layout of the HTML
	Syntax    string // Syntax and color theming printing HTML
	tmpl      string // template filename
	pack      string // template package name
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

// PageData temporarily holds template data used for the HTML layout.
type PageData struct {
	ExternalEmbed    bool
	FontEmbed        bool
	MetaGenerator    bool
	MetaNoTranslate  bool
	MetaRetroTxt     bool
	BuildVersion     string
	CacheRefresh     string
	Comment          string
	FontFamily       string
	MetaAuthor       string
	MetaColorScheme  string
	MetaDesc         string
	MetaKeywords     string
	MetaReferrer     string
	MetaRobots       string
	MetaThemeColor   string
	PageTitle        string
	PreText          string
	SauceTitle       string
	SauceAuthor      string
	SauceGroup       string
	SauceDescription string
	SauceWidth       uint
	SauceLines       uint
	CSSEmbed         template.CSS
	HTMLEmbed        template.HTML
	ScriptEmbed      template.JS
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

	zipName = "retrotxt.zip"
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
func (args *Args) Create(b *[]byte) error {
	var err error
	args.layout, err = layout(args.Layout)
	if err != nil {
		return err
	}
	switch {
	case args.Save.AsFiles:
		if err := args.saveAssets(b); err != nil {
			// --overwrite hint
			if errors.As(err, &ErrFileExist) {
				fmt.Println(logs.Hint("create [filenames] --overwrite", ErrFileExist))
				fmt.Println(str.Info() + "Use the overwrite flag to replace any existing files.")
				os.Exit(logs.OSErrCode)
			}
			return nil
		}
	case args.Save.Compress:
		const noDestination = ""
		args.zipAssets(noDestination, b)
	default:
		// print to terminal
		if err := args.Stdout(b); err != nil {
			return err
		}
	}
	return nil
}

func (args *Args) saveAssets(b *[]byte) error {
	skip := func(c chan error) {
		c <- nil
	}
	if args.Save.Destination == "" {
		dir := []string{viper.GetString("save-directory")}
		var err error
		if args.Save.Destination, err = destination(dir...); err != nil {
			logs.FatalMark(args.Save.Destination, logs.ErrFileSaveD, err)
		}
	}

	r := bytes.NewReader(*b)
	args.Source.BBSType = bbs.Find(r)

	ch, cnt := make(chan error), 0

	go args.saveHTML(b, ch)

	if useCSS(args.layout) {
		cnt++
		go args.saveStyles(ch)
	}
	if usePCBoard(args.Source.BBSType) {
		cnt += 2
		go args.saveBBS(ch)
		go args.savePCBoard(ch)
	}
	if useFontCSS(args.layout) {
		cnt++
		go args.saveFont(ch)
	}
	if useJS(args.layout) {
		cnt++
		go args.saveJS(ch)
	}
	if useIcon(args.layout) {
		cnt++
		go args.saveFavIcon(ch)
	}

	const optionalCh = 6
	skips := optionalCh - cnt
	for i := 0; i < skips; i++ {
		go skip(ch)
	}
	return check(<-ch, <-ch, <-ch, <-ch, <-ch, <-ch, <-ch)
}

func check(ch ...error) error {
	var errs error
	for _, err := range ch {
		errs = appendErr(errs, err)
	}
	return errs
}

func appendErr(errs, err error) error {
	// handle first error
	if errs == nil {
		return err
	}
	// skip duplicate errors
	if !errors.Is(err, ErrFileExist) && errors.As(errs, &err) {
		return errs
	}
	return fmt.Errorf("%s;%w", errs, err)
}

func useCSS(l Layout) bool {
	switch l {
	case Standard, Compact:
		return true
	case Inline, None:
		return false
	}
	return false
}

func useFontCSS(l Layout) bool {
	switch l {
	case Standard, Compact:
		return true
	case Inline, None:
		return false
	}
	return false
}

func useIcon(l Layout) bool {
	switch l {
	case Standard, Compact:
		return true
	case Inline, None:
		return false
	}
	return false
}

func usePCBoard(b bbs.BBS) bool {
	switch b {
	case bbs.PCBoard:
		return true
	case bbs.ANSI, bbs.Celerity, bbs.Renegade, bbs.Telegard, bbs.WWIVHash, bbs.WWIVHeart, bbs.Wildcat:
		return false
	default:
		return false
	}
}

func useJS(l Layout) bool {
	return false
}

// zipAssets compresses all assets into a single zip archive.
// An empty destination directory argument will save the zip file to the user working directory.
func (args *Args) zipAssets(destDir string, b *[]byte) {
	defer func() {
		dir := args.Save.Destination
		m, err := filepath.Match(filepath.Join(os.TempDir(), "*"), dir)
		if err != nil {
			logs.FatalMark("*", logs.ErrTmpSaveD, err)
		}
		if m {
			if err = os.RemoveAll(dir); err != nil {
				logs.FatalMark(dir, logs.ErrTmpRMD, err)
			}
		}
	}()
	var err error
	args.Save.Destination, err = ioutil.TempDir(os.TempDir(), "*-zip")
	if err != nil {
		logs.FatalMark("temporary", logs.ErrFileSaveD, err)
	}
	if err = args.saveAssets(b); err != nil {
		fmt.Println(logs.SprintWrap(logs.ErrFileSave, err))
		return
	}
	name := zipName
	if destDir != "" {
		name = filepath.Join(destDir, zipName)
	}
	zip := filesystem.Zip{
		Name:      name,
		Root:      args.Save.Destination,
		Comment:   "",
		Overwrite: args.Save.OW,
		Quiet:     args.Test,
	}
	if err = zip.Create(); err != nil {
		logs.FatalMark(name, logs.ErrZipFile, err)
	}
}

// Stdout creates and prints the HTML template.
func (args *Args) Stdout(b *[]byte) error {
	// html
	html, err := args.marshalTextTransform(b)
	if err != nil {
		return fmt.Errorf("stdout: %w", err)
	}
	// font css
	ff := args.FontFamily.Value
	f := Family(ff).String()
	if f == "" {
		return fmt.Errorf("create.saveFontCSS %q: %w", ff, ErrFont)
	}
	font, err := FontCSS(f, args.Source.Encoding, args.FontEmbed)
	if err != nil {
		return err
	}
	// print assets
	if errj := args.printJS(&static.Scripts); errj != nil {
		return errj
	}
	if errc := args.printCSS(&static.CSSStyles); errc != nil {
		return errc
	}
	if errf := args.printFontCSS(f, &font); errf != nil {
		return errf
	}
	// always print the HTML
	fmt.Printf("\nHTML file: %s\n\n", htmlFn.write())
	if err = str.Highlight(html.String(), "html", args.Syntax, true); err != nil {
		return fmt.Errorf("stdout html highlight: %w", err)
	}
	return nil
}

func colorSyntax(s string) bool {
	switch s {
	case "", none:
		return false
	}
	return str.Valid(s)
}

func (args *Args) printCSS(b *[]byte) error {
	if !useCSS(args.layout) {
		return nil
	}
	fmt.Printf("\nCSS file: %s\n\n", cssFn.write())
	if !colorSyntax(args.Syntax) {
		fmt.Println(string(*b))
		return nil
	}
	if err := str.Highlight(string(*b), "css", args.Syntax, true); err != nil {
		return fmt.Errorf("stdout css highlight: %w", err)
	}
	return nil
}

func (args *Args) printFontCSS(name string, b *[]byte) error {
	if !useFontCSS(args.layout) {
		return nil
	}
	fmt.Printf("\nCSS for %s font file: %s\n\n", name, fontFn.write())
	if err := str.Highlight(string(*b), "css", args.Syntax, true); err != nil {
		return fmt.Errorf("stdout font css highlight: %w", err)
	}
	return nil
}

func (args *Args) printJS(b *[]byte) error {
	if !useJS(args.layout) {
		return nil
	}
	fmt.Printf("\nJS file: %s\n\n", jsFn.write())
	if !colorSyntax(args.Syntax) {
		fmt.Println(string(*b))
		return nil
	}
	if err := str.Highlight(string(*b), "js", args.Syntax, true); err != nil {
		return fmt.Errorf("stdout js highlight: %w", err)
	}
	return nil
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

// destination determines if user supplied arguments are a valid file or directory destination.
func destination(args ...string) (path string, err error) {
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

// layout parses possible --layout argument values.
func layout(name string) (Layout, error) {
	switch name {
	case standard, "s":
		return Standard, nil
	case inline, "i":
		return Inline, nil
	case compact, "c":
		return Compact, nil
	case none, "n":
		return None, nil
	}
	return 0, logs.ErrTmplName
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
