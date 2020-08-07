package create

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/spf13/viper"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
	"retrotxt.com/retrotxt/internal/pack"
	"retrotxt.com/retrotxt/lib/convert"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
	"retrotxt.com/retrotxt/lib/version"

	gap "github.com/muesli/go-app-paths"
)

// Args holds arguments and options sourced from user flags or the config file.
type Args struct {
	Destination        string // Destination HTML destination either a directory or file
	Encoding           string // Encoding text encoding of the source input
	Body               string // Body text content
	Layout             string // Layout of the HTML
	Syntax             string // Syntax and color theming printing HTML
	tmpl               string // template filename
	pack               string // template package name
	FilenameVal        string
	TitleVal           string
	MetaAuthorVal      string
	MetaDescriptionVal string
	MetaColorSchemeVal string
	MetaKeywordsVal    string
	MetaReferrerVal    string
	MetaRobotsVal      string
	MetaThemeColorVal  string
	FontFamilyVal      string
	Port               uint // Port for HTTP server
	Cache              bool // Cache when false will always unpack a new .gohtml template
	Test               bool // Test mode
	SaveToFile         bool // SaveToFile save to a file or print to stdout
	OW                 bool // OW overwrite any existing files when saving
	Compress           bool // Compress and store all files into an archive
	Title              bool
	MetaAuthor         bool
	MetaDescription    bool
	MetaGeneratorVal   bool
	MetaColorScheme    bool
	MetaKeywords       bool
	MetaNoTranslateVal bool
	MetaReferrer       bool
	MetaRetroTxtVal    bool
	MetaRobots         bool
	MetaThemeColor     bool
	FontFamily         bool
	FontEmbedVal       bool
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

const (
	none = "none"
	std  = "standard"
)

// ColorScheme values for the content attribute of <meta name="color-scheme">.
var ColorScheme = [...]string{"normal", "dark light", "only light"}

// Referrer values for the content attribute of <meta name="referrer">.
var Referrer = [...]string{"no-referrer", "origin", "no-referrer-when-downgrade",
	"origin-when-cross-origin", "same-origin", "strict-origin", "strict-origin-when-cross-origin", "unsafe-URL"}

// Robots values for the content attribute of <meta name="robots">.
var Robots = [...]string{"index", "noindex", "follow", "nofollow", none, "noarchive", "nosnippet", "noimageindex", "nocache"}

// Layout are HTML template variations.
type Layout int

const (
	// Standard template with external CSS, JS, fonts.
	Standard Layout = iota
	// Inline template with CSS and JS embedded.
	Inline
	// Compact template with external CSS, JS, fonts and no meta-tags.
	Compact
	// None template, just print the generated HTML.
	None
)

func (l Layout) String() string {
	layouts := [...]string{std, "inline", "compact", none}
	if l < Standard || l > None {
		return ""
	}
	return layouts[l]
}

// Layouts are the names of the HTML templates.
func Layouts() []string {
	return []string{Standard.String(), Inline.String(), Compact.String(), None.String()}
}

// Pack is the packed name of the HTML template.
func (l Layout) Pack() string {
	packs := [...]string{std, std, std, none}
	if l < Standard || l > None {
		return ""
	}
	return packs[l]
}

// Create handles the target output command arguments.
func (args *Args) Create(b *[]byte) {
	var err error
	switch {
	case args.SaveToFile:
		// use config save directory
		// otherwise assume Destination path is a temporary --serve location
		if args.Destination == "" {
			dir := []string{viper.GetString("save-directory")}
			if args.Destination, err = destination(dir...); err != nil {
				logs.Fatal("save to directory failure", fmt.Sprintf("%s", dir), err)
			}
		}
		ch := make(chan error)
		go args.savecss(ch)
		go args.savefont(ch)
		go args.savehtml(b, ch)
		go args.savejs(ch)
		go args.savefavicon(ch)
		err1, err2, err3, err4, err5 := <-ch, <-ch, <-ch, <-ch, <-ch
		if err1 != nil {
			logs.Println("save file 1", "", err1)
		}
		if err2 != nil {
			logs.Println("save file 2", "", err2)
		}
		if err3 != nil {
			logs.Println("save file 3", "", err3)
		}
		if err4 != nil {
			logs.Println("save file 4", "", err4)
		}
		if err5 != nil {
			logs.Println("save file 5", "", err5)
		}
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
			os.Exit(1)
		}
	default:
		// print to terminal
		if err = args.Stdout(b); err != nil {
			logs.Fatal("print to stdout", "", err)
		}
	}
}

var (
	ErrReqOW     = errors.New("include an -o flag to overwrite")
	ErrPackGet   = errors.New("pack.get name is invalid")
	ErrUnknownFF = errors.New("unknown font family")
	ErrNilByte   = errors.New("cannot convert a nil byte value")
	ErrTmplDir   = errors.New("the path to the template file is a directory")
	ErrNoLayout  = errors.New("layout does not exist")
)

func (args *Args) destination(name string) (string, error) {
	dir := filesystem.DirExpansion(args.Destination)
	path := dir
	stat, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("args destination directory failed %q: %w", dir, err)
	}
	if stat.IsDir() {
		path = filepath.Join(dir, name)
	}
	_, err = os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("args destination path failed %q: %w", path, err)
	}
	if !args.OW && !os.IsNotExist(err) {
		switch name {
		case "favicon.ico", "scripts.js",
			"vga.woff2":
			// existing static files can be ignored
			return path, nil
		}
		logs.Println("file exists", path, ErrReqOW)
	} else {
		color.OpFuzzy.Printf("saving to %s\n", path)
	}
	return path, nil
}

// savecss creates and saves the styles stylesheet to the Destination argument.
func (args Args) savecss(c chan error) {
	switch args.Layout {
	case std, "s":
	default:
		c <- nil
		return
	}
	name, err := args.destination("styles.css")
	if err != nil {
		c <- err
	}
	b := pack.Get("css/styles.css")
	if len(b) == 0 {
		c <- fmt.Errorf("create.savecss %q: %w", args.pack, ErrPackGet)
	}
	_, err = filesystem.Save(name, b...)
	if err != nil {
		c <- err
	}
	c <- nil
}

func (args Args) savefavicon(c chan error) {
	switch args.Layout {
	case std, "s":
	default:
		c <- nil
		return
	}
	name, err := args.destination("favicon.ico")
	if err != nil {
		c <- err
	}
	b := pack.Get("img/retrotxt_16.png")
	if len(b) == 0 {
		c <- fmt.Errorf("create.savefavicon %q: %w", args.pack, ErrPackGet)
	}
	_, err = filesystem.Save(name, b...)
	if err != nil {
		c <- err
	}
	c <- nil
}

// savefont unpacks and saves the font binary to the Destination argument.
func (args Args) savefont(c chan error) {
	if !args.FontEmbedVal {
		f := Family(args.FontFamilyVal)
		if f.String() == "" {
			c <- fmt.Errorf("save font, could not save %q: %w", args.FontFamilyVal, ErrUnknownFF)
			return
		}
		if err := args.savefontwoff2(f.File(), "font/"+f.File()); err != nil {
			c <- err
		}
	}
	switch args.Layout {
	case std, "s":
		if err := args.savefontcss("font.css"); err != nil {
			c <- err
		}
	}
	c <- nil
}

func (args Args) savefontcss(name string) error {
	name, err := args.destination(name)
	if err != nil {
		return err
	}
	f := Family(args.FontFamilyVal).String()
	if f == "" {
		return fmt.Errorf("create.savefontcss %q: %w", name, ErrUnknownFF)
	}
	b, err := FontCSS(f, args.FontEmbedVal)
	if err != nil {
		return err
	}
	_, err = filesystem.Save(name, b...)
	if err != nil {
		return err
	}
	return nil
}

func (args Args) savefontwoff2(name, packName string) error {
	name, err := args.destination(name)
	if err != nil {
		return err
	}
	b := pack.Get(packName)
	if len(b) == 0 {
		return fmt.Errorf("create.savefontwoff2 %q: %w", args.pack, ErrPackGet)
	}
	_, err = filesystem.Save(name, b...)
	if err != nil {
		return err
	}
	return nil
}

func (args Args) savejs(c chan error) {
	switch args.Layout {
	case std, "s":
	default:
		c <- nil
		return
	}
	name, err := args.destination("scripts.js")
	if err != nil {
		c <- err
	}
	b := pack.Get("js/scripts.js")
	if len(b) == 0 {
		c <- fmt.Errorf("create.savejs %q: %w", args.pack, ErrPackGet)
	}
	_, err = filesystem.Save(name, b...)
	if err != nil {
		c <- err
	}
	c <- nil
}

// SaveHTML creates and saves the html template to the Destination argument.
func (args Args) savehtml(b *[]byte, c chan error) {
	name, err := args.destination("index.html")
	if err != nil {
		c <- err
	}
	file, err := os.Create(name)
	if err != nil {
		c <- err
	}
	defer func() {
		cerr := file.Close()
		c <- cerr
	}()
	tmpl, err := args.newTemplate()
	if err != nil {
		c <- err
	}
	d, err := args.pagedata(b)
	if err != nil {
		c <- err
	}
	if err = tmpl.Execute(file, d); err != nil {
		c <- err
	}
	c <- file.Close()
}

// Stdout creates and prints the html template.
func (args Args) Stdout(b *[]byte) error {
	tmpl, err := args.newTemplate()
	if err != nil {
		return fmt.Errorf("stdout new template failure: %w", err)
	}
	d, err := args.pagedata(b)
	if err != nil {
		return fmt.Errorf("stdout meta and pagedata failure: %w", err)
	}
	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, d); err != nil {
		return fmt.Errorf("stdout template execute failure: %w", err)
	}
	switch args.Syntax {
	case "", none:
		fmt.Printf("%s", buf.String())
	default:
		if !str.Valid(args.Syntax) {
			fmt.Printf("unknown style %q, so using none\n", args.Syntax)
			fmt.Printf("%s", buf.String())
			return nil
		}
		if err = str.Highlight(buf.String(), "html", args.Syntax); err != nil {
			return fmt.Errorf("stdout html highlight: %w", err)
		}
	}
	return nil
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
		path, err = os.UserHomeDir()
	case ".":
		path, err = os.Getwd()
	case "\\", "/":
		path, err = filepath.Abs(dir)
	}
	if err != nil {
		return "", fmt.Errorf("parse directory error: %q: %w", dir, err)
	}
	return path, nil
}

// newTemplate creates and parses a new template file.
// The argument test is used internally.
func (args Args) newTemplate() (*template.Template, error) {
	if err := args.templateCache(); err != nil {
		return nil, fmt.Errorf("using existing template cache: %w", err)
	}
	if !args.Cache {
		if err := args.templateSave(); err != nil {
			return nil, fmt.Errorf("creating a new template: %w", err)
		}
	} else {
		if s, err := os.Stat(args.tmpl); os.IsNotExist(err) {
			if err := args.templateSave(); err != nil {
				return nil, fmt.Errorf("saving to the template: %w", err)
			}
		} else if err != nil {
			return nil, err
		} else if s.IsDir() {
			return nil, fmt.Errorf("new template %q: %w", args.tmpl, ErrTmplDir)
		}
	}
	// to avoid a potential panic, Stat again in case os.IsNotExist() is true
	s, err := os.Stat(args.tmpl)
	if err != nil {
		return nil, fmt.Errorf("could not access file: %q: %w", args.tmpl, err)
	}
	if err = args.templatePack(); err != nil {
		return nil, fmt.Errorf("template pack: %w", err)
	}
	b, err := args.templateData()
	if s.Size() != int64(len(*b)) {
		if err != nil {
			return nil, fmt.Errorf("template data: %w", err)
		}
		if _, err := filesystem.Save(args.tmpl, *b...); err != nil {
			return nil, fmt.Errorf("saving template: %q: %w", args.tmpl, err)
		}
	}
	t := template.Must(template.ParseFiles(args.tmpl))
	return t, nil
}

// filename creates a filepath for the template filenames.
func (args *Args) templateCache() (err error) {
	l := layout(args.Layout).Pack()
	if l == "" {
		return fmt.Errorf("template cache %q: %w", args.Layout, ErrNoLayout)
	}
	args.tmpl, err = gap.NewScope(gap.User, "retrotxt").DataPath(l + ".gohtml")
	if err != nil {
		return fmt.Errorf("template cache path: %q: %w", args.tmpl, err)
	}
	return nil
}

func (args *Args) templatePack() error {
	l := layout(args.Layout).Pack()
	if l == "" {
		return fmt.Errorf("template pack %q: %w", args.Layout, ErrNoLayout)
	}
	args.pack = fmt.Sprintf("html/%s.gohtml", l)
	return nil
}

func (args Args) templateData() (*[]byte, error) {
	b := pack.Get(args.pack)
	if len(b) == 0 {
		return nil, fmt.Errorf("template data %q: %w", args.pack, ErrPackGet)
	}
	return &b, nil
}

func (args Args) templateSave() error {
	err := args.templatePack()
	if err != nil {
		return fmt.Errorf("template save pack error: %w", err)
	}
	b, err := args.templateData()
	if err != nil {
		return fmt.Errorf("template save data error: %w", err)
	}
	if _, err := filesystem.Save(args.tmpl, *b...); err != nil {
		return fmt.Errorf("template save error: %w", err)
	}
	return nil
}

func layout(name string) Layout {
	switch name {
	case std, "s":
		return Standard
	case "inline", "i":
		return Inline
	case "compact", "c":
		return Compact
	case none, "n":
		return None
	}
	return -1
}

// pagedata creates the meta and page template data.
func (args Args) pagedata(b *[]byte) (p PageData, err error) {
	if b == nil {
		return PageData{}, fmt.Errorf("pagedata: %w", ErrNilByte)
	}
	// templates are found in the dir static/html/*.gohtml
	switch layout(args.Layout) {
	case Inline:
		p.ExternalEmbed = true
		m := minify.New()
		m.AddFunc("text/css", css.Minify)
		m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
		// styles
		s := bytes.TrimSpace(pack.Get("css/styles.css"))
		// font
		f, err := FontCSS(args.FontFamilyVal, args.FontEmbedVal)
		if err != nil {
			return p, fmt.Errorf("pagedata font error: %w", err)
		}
		f = bytes.TrimSpace(f)
		// merge
		c := [][]byte{s, []byte("/* font */"), f}
		b := bytes.Join(c, []byte("\n\n"))
		// compress & embed
		b, err = m.Bytes("text/css", b)
		if err != nil {
			return p, fmt.Errorf("pagedata minify css: %w", err)
		}
		p.CSSEmbed = template.CSS(string(b))
		js := pack.Get("js/scripts.js")
		js, err = m.Bytes("application/javascript", js)
		if err != nil {
			return p, fmt.Errorf("pagedata minify javascript: %w", err)
		}
		p.ScriptEmbed = template.JS(string(js))
		fallthrough
	case Standard:
		p.FontEmbed = args.FontEmbedVal
		p.FontFamily = args.fontFamily()
		p.MetaAuthor = args.metaAuthor()
		p.MetaColorScheme = args.metaColorScheme()
		p.MetaDesc = args.metaDesc()
		p.MetaGenerator = args.MetaGeneratorVal
		p.MetaKeywords = args.metaKeywords()
		p.MetaNoTranslate = args.MetaNoTranslateVal
		p.MetaReferrer = args.metaReferrer()
		p.MetaRobots = args.metaRobots()
		p.MetaRetroTxt = args.MetaRetroTxtVal
		p.MetaThemeColor = args.metaThemeColor()
		p.PageTitle = args.pageTitle()
		// generate data
		t := time.Now().UTC()
		p.BuildDate = t.Format(time.RFC3339)
		p.BuildVersion = version.B.Version
	case Compact: // disables all meta tags
		p.PageTitle = args.pageTitle()
		p.MetaGenerator = false
	case None:
		// do nothing
	default:
		return PageData{}, fmt.Errorf("pagedata %s: %w", args.Layout, ErrNoLayout)
	}
	// check encoding
	var conv = convert.Args{Encoding: args.Encoding}
	if args.Encoding == "" {
		conv.Encoding = "cp437"
	}
	// convert bytes into utf8
	runes, err := conv.Text(b)
	if err != nil {
		return p, fmt.Errorf("pagedata convert text bytes to utf8 failure: %w", err)
	}
	if p.MetaRetroTxt {
		p.Comment = args.comment(conv, b, runes...)
	}
	p.PreText = string(runes)
	return p, nil
}

func (args Args) comment(c convert.Args, old *[]byte, new ...rune) string {
	e, nl, l, w, f := "", "", 0, 0, "n/a"
	b := []byte(string(new))
	// to handle EBCDIC cases, both raw bytes and utf8 runes need newline scans.
	nlr := filesystem.Newlines(false, []rune(string(*old))...)
	nl = filesystem.Newline(nlr, false)
	nnl := filesystem.Newlines(true, new...)
	e = convert.Humanize(c.Encoding)
	l, err := filesystem.Lines(bytes.NewReader(b), nnl)
	if err != nil {
		l = -1
	}
	w, err = filesystem.Columns(bytes.NewReader(b), nnl)
	if err != nil {
		w = -1
	}
	if args.FilenameVal != "" {
		f = args.FilenameVal
	}
	return fmt.Sprintf("encoding: %s; newline: %s; length: %d; width: %d; name: %s", e, nl, l, w, f)
}

func (args Args) fontFamily() string {
	if args.FontFamily {
		return args.FontFamilyVal
	}
	return viper.GetString("html.font.family")
}

func (args Args) metaAuthor() string {
	if args.MetaAuthor {
		return args.MetaAuthorVal
	}
	return viper.GetString("html.meta.author")
}

func (args Args) metaColorScheme() string {
	if args.MetaColorScheme {
		return args.MetaColorSchemeVal
	}
	return viper.GetString("html.meta.color-scheme")
}

func (args Args) metaDesc() string {
	if args.MetaDescription {
		return args.MetaDescriptionVal
	}
	return viper.GetString("html.meta.description")
}

func (args Args) metaKeywords() string {
	if args.MetaKeywords {
		return args.MetaKeywordsVal
	}
	return viper.GetString("html.meta.keywords")
}

func (args Args) metaReferrer() string {
	if args.MetaReferrer {
		return args.MetaReferrerVal
	}
	return viper.GetString("html.meta.referrer")
}

func (args Args) metaRobots() string {
	if args.MetaRobots {
		return args.MetaRobotsVal
	}
	return viper.GetString("html.meta.robots")
}

func (args Args) metaThemeColor() string {
	if args.MetaThemeColor {
		return args.MetaThemeColorVal
	}
	return viper.GetString("html.meta.theme-color")
}

func (args Args) pageTitle() string {
	if args.Title {
		return args.TitleVal
	}
	return viper.GetString("html.title")
}
