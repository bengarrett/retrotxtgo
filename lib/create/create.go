package create

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bengarrett/retrotxtgo/internal/pack"
	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/lib/version"
	"github.com/gookit/color"
	"github.com/spf13/viper"

	gap "github.com/muesli/go-app-paths"
)

type files map[string]string

// Args holds arguments and options sourced from user flags or the config file.
type Args struct {
	// Dest HTML destination either a directory or file
	Dest string
	// Enc text encoding of the source input
	Enc string
	// Author of the page metadata
	Author string
	// Scheme color metadata
	Scheme string
	// Desc description metadata
	Desc string
	// FontEmbed embeds font data into the CSS
	FontEmbed bool
	// FontFamily font name
	FontFamily string
	// Generator shows retrotxt version and page generated at time
	Generator bool
	// Keys are keyword metadata
	Keys string
	// Google notranslate
	NoTranslate bool
	// Referrer metadata
	Ref string
	// Robots metadata
	Robots string
	// Title for the page and browser tab
	Title string
	// Body text content
	Body string
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
}

// PageData temporarily holds template data used for the HTML layout.
type PageData struct {
	BuildVersion    string
	BuildDate       string
	CacheRefresh    string
	Comment         string
	FontEmbed       bool
	FontFamily      string
	MetaAuthor      string
	MetaColorScheme string
	MetaDesc        string
	MetaGenerator   bool
	MetaKeywords    string
	MetaNoTranslate bool
	MetaReferrer    string
	MetaRobots      string
	MetaThemeColor  string
	PageTitle       string
	PreText         string
}

// ColorScheme values for the content attribute of <meta name="color-scheme">
var ColorScheme = []string{"normal", "dark light", "only light"}

// FontFamily values for the CSS font-family.
var FontFamily = []string{"automatic", "mona", "vga"}

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

// Create handles the target output command arguments.
func (args *Args) Create(b *[]byte) {
	var err error
	switch {
	case args.SaveToFile:
		// use config save directory
		dir := []string{viper.GetString("save-directory")}
		if args.Dest, err = destination(dir); err != nil {
			log.Fatal(err)
		}
		err = args.savecss()
		if err != nil {
			log.Fatal(err) // TODO: logs.Fatal(error)
		}
		err = args.savefont()
		if err != nil {
			log.Fatal(err) // TODO: logs.Fatal(error)
		}
		err = args.savehtml(b)
		if err != nil {
			log.Fatal(err) // TODO: logs.Fatal(error)
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
	dir := filesystem.DirExpansion(args.Dest)
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
		logs.ChkErr(e)
	}
	return path, nil
}

// savecss creates and saves the styles stylesheet to the Dest argument.
func (args Args) savecss() error {
	name, err := args.destination("style.css")
	if err != nil {
		return err
	}
	b := pack.Get("css/styles.css")
	if len(b) == 0 {
		return fmt.Errorf("create.savecss: pack.get name is invalid: %q", args.pack)
	}
	_, err = filesystem.Save(name, b)
	if err != nil {
		return err
	}
	return nil
}

// savefont unpacks and saves the font binary to the Dest argument.
func (args Args) savefont() error {
	if err := args.savefontwoff2("vga.woff2", "font/ibm-vga8.woff2"); err != nil {
		return err
	}
	if err := args.savefontcss("font.css", "css/font-vga.css"); err != nil {
		return err
	}
	return nil
}

func (args Args) savefontcss(name, packName string) error {
	name, err := args.destination(name)
	if err != nil {
		return err
	}
	b := pack.Get(packName)
	if len(b) == 0 {
		return fmt.Errorf("create.savefontcss: pack.get name is invalid: %q", args.pack)
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

// SaveHTML creates and saves the html template to the Dest argument.
func (args Args) savehtml(b *[]byte) error {
	name, err := args.destination("index.html")
	if err != nil {
		return err
	}
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	tmpl, err := args.newTemplate()
	if err != nil {
		return err
	}
	d, err := args.pagedata(b)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(file, d); err != nil {
		return err
	}
	return nil
}

// Stdout creates and sends the html template to stdout.
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
		if err = str.Highlight(buf.String(), "html", args.Syntax); err != nil {
			return err
		}
	}
	return err
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
	for key := range createTemplates() {
		s = append(s, key)
	}
	return s
}

// createTemplates creates a map of the template filenames used in conjunction with the layout flag.
func createTemplates() files {
	f := make(files)
	f["body"] = "body-content"
	f["full"] = "standard"
	f["mini"] = "standard"
	f["pre"] = "pre-content"
	f["standard"] = "standard"
	return f
}

// newTemplate creates and parses a new template file.
// The argument test is used internally.
func (args Args) newTemplate() (*template.Template, error) {
	if err := args.templateCache(); err != nil {
		return nil, err
	}
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
	f := createTemplates()[args.Layout]
	if f == "" {
		return fmt.Errorf("create.templatecache: layout does not exist: %q", args.Layout)
	}
	args.tmpl, err = scope.DataPath(f + ".html")
	return err
}

func (args *Args) templatePack() error {
	f := createTemplates()[args.Layout]
	if f == "" {
		return fmt.Errorf("create.templatepack: package and layout does not exist: %q", args.Layout)
	}
	args.pack = fmt.Sprintf("html/%s.html", f)
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
	// templates are found in the dir static/html/*.html
	switch args.Layout {
	case "full", "standard":
		p.FontEmbed = viper.GetBool("html.font.embed")
		p.FontFamily = viper.GetString("html.font.family")
		p.MetaAuthor = viper.GetString("html.meta.author")
		p.MetaColorScheme = viper.GetString("html.meta.color-scheme")
		p.MetaDesc = viper.GetString("html.meta.description")
		p.MetaGenerator = viper.GetBool("html.meta.generator")
		p.MetaKeywords = viper.GetString("html.meta.keywords")
		p.MetaNoTranslate = viper.GetBool("html.meta.notranslate")
		p.MetaReferrer = viper.GetString("html.meta.referrer")
		p.MetaRobots = viper.GetString("html.meta.robots")
		p.MetaThemeColor = viper.GetString("html.meta.theme-color")
		p.PageTitle = viper.GetString("html.title")
		// generate data
		t := time.Now().UTC()
		p.BuildDate = t.Format(time.RFC3339)
		p.BuildVersion = version.B.Version
		p.Comment = "encoding: CP-437; linefeed: crlf; length: 100; width: 80; filename: somefile.txt" // TODO: make functional
	case "mini":
		p.PageTitle = viper.GetString("html.title")
		p.MetaGenerator = false
	}
	// convert bytes into utf8
	var name = args.Enc
	if name == "" {
		name = "cp437"
	}
	runes, err := convert.Text(name, b)
	logs.Check("create.pagedata.chars", err)
	p.PreText = string(runes)
	return p, nil
}
