package create

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
	// Destination HTML destination either a directory or file
	Destination string
	// Encoding text encoding of the source input
	Encoding string
	// Body text content
	Body string
	// Cache when false will always unpack a new .gohtml template
	Cache bool
	// Layout of the HTML
	Layout string
	// Syntax and color theming printing HTML
	Syntax string
	// Port for HTTP server
	Port uint
	// Test mode
	Test bool
	// SaveToFile save to a file or print to stdout
	SaveToFile bool
	// OW overwrite any existing files when saving
	OW bool
	// template filename
	tmpl string
	// template package name
	pack string

	// Flag values, command arguments and change state

	FilenameVal        string
	TitleVal           string
	Title              bool
	MetaAuthorVal      string
	MetaAuthor         bool
	MetaDescriptionVal string
	MetaDescription    bool
	MetaGeneratorVal   bool
	MetaColorSchemeVal string
	MetaColorScheme    bool
	MetaKeywordsVal    string
	MetaKeywords       bool
	MetaNoTranslateVal bool
	MetaReferrerVal    string
	MetaReferrer       bool
	MetaRetroTxtVal    bool
	MetaRobotsVal      string
	MetaRobots         bool
	MetaThemeColorVal  string
	MetaThemeColor     bool
	FontFamilyVal      string
	FontFamily         bool
	FontEmbedVal       bool
}

// PageData temporarily holds template data used for the HTML layout.
type PageData struct {
	BuildVersion    string
	BuildDate       string
	CacheRefresh    string
	Comment         string
	CSSEmbed        template.CSS
	ExternalEmbed   bool
	FontEmbed       bool
	FontFamily      string
	MetaAuthor      string
	MetaColorScheme string
	MetaDesc        string
	MetaGenerator   bool
	MetaKeywords    string
	MetaNoTranslate bool
	MetaReferrer    string
	MetaRetroTxt    bool
	MetaRobots      string
	MetaThemeColor  string
	PageTitle       string
	PreText         string
	ScriptEmbed     template.JS
}

// ColorScheme values for the content attribute of <meta name="color-scheme">
var ColorScheme = []string{"normal", "dark light", "only light"}

// Referrer values for the content attribute of <meta name="referrer">
var Referrer = []string{"no-referrer", "origin", "no-referrer-when-downgrade",
	"origin-when-cross-origin", "same-origin", "strict-origin", "strict-origin-when-cross-origin", "unsafe-URL"}

// Robots values for the content attribute of <meta name="robots">
var Robots = []string{"index", "noindex", "follow", "nofollow", "none", "noarchive", "nosnippet", "noimageindex", "nocache"}

var scope = gap.NewScope(gap.User, "retrotxt")

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
		return "", err
	}
	return path, err
}

// TmplsPacks are a map of the named packed templates.
type TmplsPacks map[string]string

// Templates creates a map of the template packed names used in conjunction with the layout flag.
func Templates() TmplsPacks {
	f := make(TmplsPacks, 5)
	f["standard"] = "standard" // standard template with external CSS, JS, fonts
	f["inline"] = "standard"   // standard template with CSS and JS embedded
	f["compact"] = "standard"  // standard template with external CSS, JS, fonts and no meta-tags
	f["none"] = "none"         // no template, just print the generated HTML
	return f
}

func (t TmplsPacks) String() string {
	return str.UnderlineKeys(t.Keys())
}

// Keys method returns the keys of templs as a sorted slice.
func (t TmplsPacks) Keys() (s []string) {
	for key := range t {
		s = append(s, key)
	}
	sort.Strings(s)
	return s
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
			if args.Destination, err = destination(dir); err != nil {
				log.Fatal(err)
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
			log.Fatal(err1)
		}
		if err2 != nil {
			log.Fatal(err2)
		}
		if err3 != nil {
			log.Fatal(err3)
		}
		if err4 != nil {
			log.Fatal(err4)
		}
		if err5 != nil {
			log.Fatal(err5)
		}
	default:
		// print to terminal
		err = args.Stdout(b)
		if err != nil {
			log.Fatal(err) // TODO: logs.Fatal(error)
		}
	}
}

func (args *Args) destination(name string) (string, error) {
	dir := filesystem.DirExpansion(args.Destination)
	path := dir
	stat, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("%s %q", err, dir)
	}
	if stat.IsDir() {
		path = filepath.Join(dir, name)
	}
	color.OpFuzzy.Printf("Saving to %s\n", path)
	stat, err = os.Stat(path)
	if !args.OW && !os.IsNotExist(err) {
		e := logs.Err{Issue: "file exists", Arg: path, Msg: errors.New("include an -o flag to overwrite")}
		fmt.Println(e)
	}
	return path, nil
}

// savecss creates and saves the styles stylesheet to the Destination argument.
func (args Args) savecss(c chan error) {
	if args.Layout != "standard" {
		c <- nil
		return
	}
	name, err := args.destination("styles.css")
	if err != nil {
		c <- err
	}
	b := pack.Get("css/styles.css")
	if len(b) == 0 {
		c <- fmt.Errorf("create.savecss: pack.get name is invalid: %q", args.pack)
	}
	_, err = filesystem.Save(name, b)
	if err != nil {
		c <- err
	}
	c <- nil
}

func (args Args) savefavicon(c chan error) {
	if args.Layout != "standard" {
		c <- nil
		return
	}
	name, err := args.destination("favicon.ico")
	if err != nil {
		c <- err
	}
	b := pack.Get("img/retrotxt_16.png")
	if len(b) == 0 {
		c <- fmt.Errorf("create.savefavicon: pack.get name is invalid: %q", args.pack)
	}
	_, err = filesystem.Save(name, b)
	if err != nil {
		c <- err
	}
	c <- nil
}

// savefont unpacks and saves the font binary to the Destination argument.
func (args Args) savefont(c chan error) {
	// TODO: handle different font names
	if !args.FontEmbedVal {
		if err := args.savefontwoff2("vga.woff2", "font/ibm-vga8.woff2"); err != nil {
			c <- err
		}
	}
	if args.Layout == "standard" {
		if err := args.savefontcss("font.css", "css/font.css"); err != nil {
			c <- err
		}
	}
	c <- nil
}

func (args Args) savefontcss(name, packName string) error {
	name, err := args.destination(name)
	if err != nil {
		return err
	}
	f := SelectFont(args.FontFamilyVal)
	if f == "" {
		return errors.New("create.savefontcss: unknown font name: " + name)
	}
	b, err := FontCSS(f, args.FontEmbedVal)
	if err != nil {
		return err
	}
	_, err = filesystem.Save(name, b)
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
		return fmt.Errorf("create.savefontwoff2: pack.get name is invalid: %q", args.pack)
	}
	_, err = filesystem.Save(name, b)
	if err != nil {
		return err
	}
	return nil
}

func (args Args) savejs(c chan error) {
	if args.Layout != "standard" {
		c <- nil
		return
	}
	name, err := args.destination("scripts.js")
	if err != nil {
		c <- err
	}
	b := pack.Get("js/scripts.js")
	if len(b) == 0 {
		c <- fmt.Errorf("create.savejs: pack.get name is invalid: %q", args.pack)
	}
	_, err = filesystem.Save(name, b)
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
	defer file.Close()
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
	c <- nil
}

// Stdout creates and prints the html template.
func (args Args) Stdout(b *[]byte) error {
	tmpl, err := args.newTemplate()
	if err != nil {
		return err
	}
	d, err := args.pagedata(b)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, d); err != nil {
		return err
	}
	switch args.Syntax {
	case "", "none":
		fmt.Printf("%s", buf.String())
	default:
		if !str.IsStyle(args.Syntax) {
			fmt.Printf("unknown style %q, so using none\n", args.Syntax)
			fmt.Printf("%s", buf.String())
			return nil
		}
		if err = str.Highlight(buf.String(), "html", args.Syntax); err != nil {
			return err
		}
	}
	return nil
}

// destination determines if user supplied arguments are a valid file or directory destination.
func destination(args []string) (path string, err error) {
	if len(args) == 0 {
		return path, err
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
			return path, err
		}
	}
	path = strings.Join(part, string(os.PathSeparator))
	return path, err
}

// Layouts returns the options permitted by the layout flag.
func Layouts() (s []string) {
	for key := range Templates() {
		s = append(s, key)
	}
	return s
}

func layout(name string) string {
	if name == "" {
		return ""
	}
	for key, val := range Templates() {
		if name == key {
			return val
		}
		if name == key[:1] {
			return val
		}
	}
	return ""
}

// newTemplate creates and parses a new template file.
// The argument test is used internally.
func (args Args) newTemplate() (*template.Template, error) {
	if err := args.templateCache(); err != nil {
		return nil, err
	}
	if !args.Cache {
		if err := args.templateSave(); err != nil {
			return nil, err
		}
		println("template cache saved to:", args.tmpl)
	} else {
		if s, err := os.Stat(args.tmpl); os.IsNotExist(err) {
			if err := args.templateSave(); err != nil {
				return nil, err
			}
			println("template cache saved to:", args.tmpl)
		} else if err != nil {
			return nil, err
		} else if s.IsDir() {
			return nil, fmt.Errorf("create.newtemplate: template file is a directory: %q", args.tmpl)
		}
	}
	// to avoid a potential panic, Stat again incase os.IsNotExist() is true
	s, err := os.Stat(args.tmpl)
	if err != nil {
		return nil, err
	}
	if err = args.templatePack(); err != nil {
		return nil, err
	}
	b, err := args.templateData()
	if s.Size() != int64(len(*b)) {
		if err != nil {
			return nil, err
		}
		if _, err := filesystem.Save(args.tmpl, *b); err != nil {
			return nil, err
		}
	}
	t := template.Must(template.ParseFiles(args.tmpl))
	return t, nil
}

// filename creates a filepath for the template filenames.
func (args *Args) templateCache() (err error) {
	l := layout(args.Layout)
	if l == "" {
		return fmt.Errorf("create.templatecache: layout does not exist: %q", args.Layout)
	}
	args.tmpl, err = scope.DataPath(l + ".gohtml")
	return err
}

func (args *Args) templatePack() error {
	l := layout(args.Layout)
	if l == "" {
		return fmt.Errorf("create.templatepack: package and layout does not exist: %q", args.Layout)
	}
	args.pack = fmt.Sprintf("html/%s.gohtml", l)
	return nil
}

func (args Args) templateData() (*[]byte, error) {
	b := pack.Get(args.pack)
	if len(b) == 0 {
		return nil, fmt.Errorf("create.templatedata: pack.get name is invalid: %q", args.pack)
	}
	return &b, nil
}

func (args Args) templateSave() error {
	err := args.templatePack()
	if err != nil {
		return err
	}
	b, err := args.templateData()
	if err != nil {
		return err
	}
	if _, err := filesystem.Save(args.tmpl, *b); err != nil {
		return err
	}
	return nil
}

// pagedata creates the meta and page template data.
func (args Args) pagedata(b *[]byte) (p PageData, err error) {
	if b == nil {
		return PageData{}, errors.New("create.pagedata: cannot convert b <nil>")
	}
	// templates are found in the dir static/html/*.gohtml
	switch args.Layout {
	case "inline":
		p.ExternalEmbed = true
		m := minify.New()
		m.AddFunc("text/css", css.Minify)
		m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
		// styles
		s := bytes.TrimSpace(pack.Get("css/styles.css"))
		// font
		f, err := FontCSS(args.FontFamilyVal, args.FontEmbedVal)
		if err != nil {
			return p, err
		}
		f = bytes.TrimSpace(f)
		// merge
		c := [][]byte{s, []byte("/* font */"), f}
		b := bytes.Join(c, []byte("\n\n"))
		// compress & embed
		b, err = m.Bytes("text/css", b)
		if err != nil {
			panic(err)
		}
		p.CSSEmbed = template.CSS(string(b))
		js := pack.Get("js/scripts.js")
		js, err = m.Bytes("application/javascript", js)
		if err != nil {
			panic(err)
		}
		p.ScriptEmbed = template.JS(string(js))
		fallthrough
	case "standard":
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
	case "compact": // disables all meta tags
		p.PageTitle = args.pageTitle()
		p.MetaGenerator = false
	case "none":
		// do nothing
	default:
		return PageData{}, errors.New("create.pagedata: unknown layout: " + args.Layout)
	}
	// check encoding
	var conv = convert.Args{Encoding: args.Encoding}
	if args.Encoding == "" {
		conv.Encoding = "cp437"
	}
	// convert bytes into utf8
	runes, err := conv.Text(b)
	if p.MetaRetroTxt {
		p.Comment = args.comment(conv, b, runes)
	}
	logs.Check("create.pagedata.chars", err)
	p.PreText = string(runes)
	fmt.Println(args)
	return p, nil
}

func (args Args) comment(c convert.Args, old *[]byte, new []rune) string {
	e, nl, l, w, f := "", "", 0, 0, "n/a"
	b := []byte(string(new))
	// to handle EBCDIC cases, both raw bytes and utf8 runes need newline scans.
	nlr := filesystem.Newlines([]rune(string(*old)), false)
	nl = filesystem.Newline(nlr, false)
	nnl := filesystem.Newlines(new, true)
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
