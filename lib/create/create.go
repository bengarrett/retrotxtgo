// Package create makes HTML and other web resources from a text file.
package create

import (
	"bufio"
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
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"retrotxt.com/retrotxt/internal/pack"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/humanize"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
	"retrotxt.com/retrotxt/lib/version"

	gap "github.com/muesli/go-app-paths"
)

// Args holds arguments and options sourced from user flags or the config file.
type Args struct {
	Source struct {
		Encoding   encoding.Encoding // Original encoding of the text source
		HiddenBody string            // Pre-text content override, accessible by a hidden flag
		Name       string            // Text source, usually a file or pack name
	}
	Output struct {
		Cache       bool   // Cache, when false will always unpack a new .gohtml template
		Compress    bool   // TODO: Compress and store all files into an archive
		OW          bool   // OW overwrite any existing files when saving
		SaveToFile  bool   // TODO: SaveToFile save to a file or print to stdout
		Destination string // TODO: Destination HTML destination either a directory or file
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
	return [...]string{unknown, standard, "inline", "compact", none}[l]
}

const (
	none     = "none"
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
	case args.Output.SaveToFile:
		// use config save directory
		// otherwise assume Destination path is a temporary --serve location
		if args.Output.Destination == "" {
			dir := []string{viper.GetString("save-directory")}
			if args.Output.Destination, err = destination(dir...); err != nil {
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

func (args *Args) destination(name string) (string, error) {
	dir := filesystem.DirExpansion(args.Output.Destination)
	path := dir
	stat, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("args destination directory failed %q: %w", dir, err)
	}
	if stat.IsDir() {
		path = filepath.Join(dir, name)
	}
	_, err = os.Stat(path)
	if !args.Output.OW && !os.IsNotExist(err) {
		switch name {
		case "favicon.ico", "scripts.js", "vga.woff2":
			// existing static files can be ignored
			return path, nil
		}
		logs.Println("file exists", path, ErrReqOW)
		return path, nil
	}
	if os.IsNotExist(err) {
		empty := []byte{}
		if _, _, err = filesystem.Save(path, empty...); err != nil {
			return "", fmt.Errorf("args destination path failed %q: %w", path, err)
		}
	}
	return path, nil
}

// saveCSS creates and saves the styles stylesheet to the Destination argument.
func (args *Args) saveCSS(c chan error) {
	switch args.layout {
	case Standard:
	case Compact, Inline, None:
		c <- nil
	}
	name, err := args.destination("styles.css")
	if err != nil {
		c <- err
	}
	b := pack.Get("css/styles.css")
	if len(b) == 0 {
		c <- fmt.Errorf("create.saveCSS %q: %w", args.pack, ErrPackGet)
	}
	nn, _, err := filesystem.Save(name, b...)
	if err != nil {
		c <- err
	}
	bytesStats(name, nn)
	c <- nil
}

func (args *Args) saveFavIcon(c chan error) {
	switch args.layout {
	case Standard:
	case Compact, Inline, None:
		c <- nil
	}
	name, err := args.destination("favicon.ico")
	if err != nil {
		c <- err
	}
	b := pack.Get("img/retrotxt_16.png")
	if len(b) == 0 {
		c <- fmt.Errorf("create.saveFavIcon %q: %w", args.pack, ErrPackGet)
	}
	nn, _, err := filesystem.Save(name, b...)
	if err != nil {
		c <- err
	}
	bytesStats(name, nn)
	c <- nil
}

// saveFont unpacks and saves the font binary to the Destination argument.
func (args *Args) saveFont(c chan error) {
	if !args.FontEmbed {
		f := Family(args.FontFamily.Value)
		if f.String() == "" {
			c <- fmt.Errorf("save font, could not save %q: %w", args.FontFamily.Value, ErrUnknownFF)
			return
		}
		if err := args.saveFontWoff2(f.File(), "font/"+f.File()); err != nil {
			c <- err
		}
	}
	switch args.layout {
	case Standard:
		if err := args.saveFontCSS("font.css"); err != nil {
			c <- err
		}
	case Compact, Inline, None:
	}
	c <- nil
}

func (args *Args) saveFontCSS(name string) error {
	name, err := args.destination(name)
	if err != nil {
		return err
	}
	f := Family(args.FontFamily.Value).String()
	if f == "" {
		return fmt.Errorf("create.saveFontCSS %q: %w", name, ErrUnknownFF)
	}
	b, err := FontCSS(f, args.FontEmbed)
	if err != nil {
		return err
	}
	nn, _, err := filesystem.Save(name, b...)
	if err != nil {
		return err
	}
	bytesStats(name, nn)
	return nil
}

func (args *Args) saveFontWoff2(name, packName string) error {
	name, err := args.destination(name)
	if err != nil {
		return err
	}
	b := pack.Get(packName)
	if len(b) == 0 {
		return fmt.Errorf("create.saveFontWoff2 %q: %w", args.pack, ErrPackGet)
	}
	nn, _, err := filesystem.Save(name, b...)
	if err != nil {
		return err
	}
	bytesStats(name, nn)
	return nil
}

func (args *Args) saveJS(c chan error) {
	switch args.layout {
	case Standard:
	case Compact, Inline, None:
		c <- nil
		return
	}
	name, err := args.destination("scripts.js")
	if err != nil {
		c <- err
	}
	b := pack.Get("js/scripts.js")
	if len(b) == 0 {
		c <- fmt.Errorf("create.saveJS %q: %w", args.pack, ErrPackGet)
	}
	nn, _, err := filesystem.Save(name, b...)
	if err != nil {
		c <- err
	}
	bytesStats(name, nn)
	c <- nil
}

// SaveHTML creates and saves the html template to the Destination argument.
func (args *Args) saveHTML(b *[]byte, c chan error) {
	name, err := args.destination("index.html")
	if err != nil {
		c <- err
	}
	if name == "" {
		c <- ErrEmptyName
	}
	file, err := os.Create(name)
	if err != nil {
		c <- err
	}
	defer func() {
		cerr := file.Close()
		c <- cerr
	}()
	buf, err := args.marshalTextTransform(b)
	if err != nil {
		c <- err
	}
	w := bufio.NewWriter(file)
	nn, err := w.Write(buf.Bytes())
	if err != nil {
		c <- err
	}
	bytesStats(name, nn)
	if err := w.Flush(); err != nil {
		c <- err
	}
	c <- file.Close()
}

func bytesStats(name string, nn int) {
	const kB = 1000
	h := humanize.Decimal(int64(nn), language.AmericanEnglish)
	color.OpFuzzy.Printf("saved to %s", name)
	switch {
	case nn == 0:
		color.OpFuzzy.Printf("saved to %s (zero-byte file)", name)
	case nn < kB:
		color.OpFuzzy.Printf(", %s", h)
	default:
		color.OpFuzzy.Printf(", %s (%d)", h, nn)
	}
	fmt.Print("\n")
}

func (args *Args) marshalTextTransform(b *[]byte) (buf bytes.Buffer, err error) {
	tmpl, err := args.newTemplate()
	if err != nil {
		return buf, fmt.Errorf("stdout new template failure: %w", err)
	}
	d, err := args.marshal(b)
	if err != nil {
		return buf, fmt.Errorf("stdout meta and pagedata failure: %w", err)
	}
	if err = tmpl.Execute(&buf, d); err != nil {
		return buf, fmt.Errorf("stdout template execute failure: %w", err)
	}
	return buf, nil
}

// Stdout creates and prints the html template.
func (args *Args) Stdout(b *[]byte) error {
	buf, err := args.marshalTextTransform(b)
	if err != nil {
		return fmt.Errorf("stdout: %w", err)
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
		if err = str.Highlight(buf.String(), "html", args.Syntax, true); err != nil {
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
func (args *Args) newTemplate() (*template.Template, error) {
	if err := args.templateCache(); err != nil {
		return nil, fmt.Errorf("using existing template cache: %w", err)
	}
	if !args.Output.Cache {
		if err := args.templateSave(); err != nil {
			return nil, fmt.Errorf("creating a new template: %w", err)
		}
	} else {
		s, err := os.Stat(args.tmpl)
		switch {
		case os.IsNotExist(err):
			if err2 := args.templateSave(); err2 != nil {
				return nil, fmt.Errorf("saving to the template: %w", err2)
			}
		case err != nil:
			return nil, err
		case s.IsDir():
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
		if _, _, err := filesystem.Save(args.tmpl, *b...); err != nil {
			return nil, fmt.Errorf("saving template: %q: %w", args.tmpl, err)
		}
	}
	t := template.Must(template.ParseFiles(args.tmpl))
	return t, nil
}

// filename creates a filepath for the template filenames.
func (args *Args) templateCache() (err error) {
	l := args.layout.Pack()
	if l == "" {
		return fmt.Errorf("template cache %q: %w", args.layout, ErrNoLayout)
	}
	args.tmpl, err = gap.NewScope(gap.User, "retrotxt").DataPath(l + ".gohtml")
	if err != nil {
		return fmt.Errorf("template cache path: %q: %w", args.tmpl, err)
	}
	return nil
}

func (args *Args) templatePack() error {
	l := args.layout.Pack()
	if l == "" {
		return fmt.Errorf("template pack %q: %w", args.layout, ErrNoLayout)
	}
	args.pack = fmt.Sprintf("html/%s.gohtml", l)
	return nil
}

func (args *Args) templateData() (*[]byte, error) {
	b := pack.Get(args.pack)
	if len(b) == 0 {
		return nil, fmt.Errorf("template data %q: %w", args.pack, ErrPackGet)
	}
	return &b, nil
}

func (args *Args) templateSave() error {
	err := args.templatePack()
	if err != nil {
		return fmt.Errorf("template save pack error: %w", err)
	}
	b, err := args.templateData()
	if err != nil {
		return fmt.Errorf("template save data error: %w", err)
	}
	if _, _, err := filesystem.Save(args.tmpl, *b...); err != nil {
		return fmt.Errorf("template save error: %w", err)
	}
	return nil
}

func layout(name string) (l Layout) {
	switch name {
	case standard, "s":
		return Standard
	case "inline", "i":
		return Inline
	case "compact", "c":
		return Compact
	case none, "n":
		return None
	}
	return l
}

func (args *Args) marshalCompact(p *PageData) PageData {
	p.PageTitle = args.pageTitle()
	p.MetaGenerator = false
	return *p
}

func (args *Args) marshalInline(b *[]byte) (p PageData, err error) {
	if b == nil {
		return PageData{}, fmt.Errorf("pagedata: %w", ErrNilByte)
	}
	p.ExternalEmbed = true
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	// styles
	s := bytes.TrimSpace(pack.Get("css/styles.css"))
	// font
	var f []byte
	f, err = FontCSS(args.FontFamily.Value, args.FontEmbed)
	if err != nil {
		return p, fmt.Errorf("pagedata font error: %w", err)
	}
	f = bytes.TrimSpace(f)
	// merge
	c := [][]byte{s, []byte("/* font */"), f}
	*b = bytes.Join(c, []byte("\n\n"))
	// compress & embed
	*b, err = m.Bytes("text/css", *b)
	if err != nil {
		return p, fmt.Errorf("pagedata minify css: %w", err)
	}
	p.CSSEmbed = template.CSS(string(*b))
	jsp := pack.Get("js/scripts.js")
	jsp, err = m.Bytes("application/javascript", jsp)
	if err != nil {
		return p, fmt.Errorf("pagedata minify javascript: %w", err)
	}
	p.ScriptEmbed = template.JS(string(jsp))
	return p, nil
}

func (args *Args) marshalStandard(p *PageData) PageData {
	p.FontEmbed = args.FontEmbed
	p.FontFamily = args.fontFamily()
	p.MetaAuthor = args.metaAuthor()
	p.MetaColorScheme = args.metaColorScheme()
	p.MetaDesc = args.metaDesc()
	p.MetaGenerator = args.Metadata.Generator
	p.MetaKeywords = args.metaKeywords()
	p.MetaNoTranslate = args.Metadata.NoTranslate
	p.MetaReferrer = args.metaReferrer()
	p.MetaRobots = args.metaRobots()
	p.MetaRetroTxt = args.Metadata.RetroTxt
	p.MetaThemeColor = args.metaThemeColor()
	p.PageTitle = args.pageTitle() + "honk!"
	// generate data
	t := time.Now().UTC()
	p.BuildDate = t.Format(time.RFC3339)
	p.BuildVersion = version.B.Version
	return *p
}

// TODO: reorder etc.

// ReplaceNELs todo: placeholder todo.
func ReplaceNELs() runes.Transformer {
	return runes.Map(func(r rune) rune {
		if r == filesystem.NextLine {
			return filesystem.Linefeed
		}
		return r
	})
}

// Normalize runes to bytes by making adjustments to text control codes.
func Normalize(e encoding.Encoding, r ...rune) (b []byte) {
	s := ""
	switch e {
	case charmap.CodePage037, charmap.CodePage1047, charmap.CodePage1140:
		s, _, _ = transform.String(ReplaceNELs(), string(r))
	default:
		s = string(r)
	}
	return []byte(s)
}

// marshal transforms bytes into UTF-8, creates the page meta and template data.
func (args *Args) marshal(b *[]byte) (p PageData, err error) {
	if b == nil {
		return PageData{}, fmt.Errorf("pagedata: %w", ErrNilByte)
	}
	// templates are found in the dir static/html/*.gohtml
	switch args.layout {
	case Inline:
		if p, err = args.marshalInline(b); err != nil {
			return p, err
		}
		p = args.marshalStandard(&p)
	case Standard:
		p = args.marshalStandard(&p)
	case Compact: // disables all meta tags
		p = args.marshalCompact(&p)
	case None:
		// do nothing
	default:
		return PageData{}, fmt.Errorf("pagedata %s: %w", args.layout, ErrNoLayout)
	}
	// convert bytes into utf8
	r := bytes.Runes(*b)
	p.PreText = string(r)
	if p.MetaRetroTxt {
		lb := filesystem.NEL() // TODO: replace
		p.Comment = args.comment(lb, r...)
	}
	return p, nil
}

func (args *Args) comment(lb filesystem.LB, r ...rune) string {
	l, w, f := 0, 0, "n/a"
	b, lbs, e := []byte(string(r)),
		filesystem.LineBreak(lb, false),
		args.Source.Encoding
	l, err := filesystem.Lines(bytes.NewReader(b), filesystem.NEL())
	if err != nil {
		l = -1
	}
	w, err = filesystem.Columns(bytes.NewReader(b), filesystem.NEL())
	if err != nil {
		w = -1
	}
	if args.Source.Name != "" {
		f = args.Source.Name
	}
	return fmt.Sprintf("encoding: %s; line break: %s; length: %d; width: %d; name: %s", e, lbs, l, w, f)
}

func (args *Args) fontFamily() string {
	if args.FontFamily.Flag {
		return args.FontFamily.Value
	}
	return viper.GetString("html.font.family")
}

func (args *Args) metaAuthor() string {
	if args.Metadata.Author.Flag {
		return args.Metadata.Author.Value
	}
	return viper.GetString("html.meta.author")
}

func (args *Args) metaColorScheme() string {
	if args.Metadata.ColorScheme.Flag {
		return args.Metadata.ColorScheme.Value
	}
	return viper.GetString("html.meta.color-scheme")
}

func (args *Args) metaDesc() string {
	if args.Metadata.Description.Flag {
		return args.Metadata.Description.Value
	}
	return viper.GetString("html.meta.description")
}

func (args *Args) metaKeywords() string {
	if args.Metadata.Keywords.Flag {
		return args.Metadata.Keywords.Value
	}
	return viper.GetString("html.meta.keywords")
}

func (args *Args) metaReferrer() string {
	if args.Metadata.Referrer.Flag {
		return args.Metadata.Referrer.Value
	}
	return viper.GetString("html.meta.referrer")
}

func (args *Args) metaRobots() string {
	if args.Metadata.Referrer.Flag {
		return args.Metadata.Referrer.Value
	}
	return viper.GetString("html.meta.robots")
}

func (args *Args) metaThemeColor() string {
	if args.Metadata.ThemeColor.Flag {
		return args.Metadata.ThemeColor.Value
	}
	return viper.GetString("html.meta.theme-color")
}

func (args *Args) pageTitle() string {
	if args.Title.Flag {
		return args.Title.Value
	}
	return viper.GetString("html.title")
}
