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

	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/lib/version"
	"github.com/gookit/color"
	"github.com/spf13/viper"
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
	// HTTP server
	HTTP bool
	// Port for HTTP server
	Port uint
	// Test mode
	Test bool
	// SaveToFile save to a file or print to stdout
	SaveToFile bool
	// OW overwrite any existing files when saving
	OW bool
}

// PageData temporarily holds template data used for the HTML layout.
type PageData struct {
	BuildVersion    string
	BuildDate       string
	CacheRefresh    string
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

// Referrer values for the content attribute of <meta name="referrer">
var Referrer = []string{"no-referrer", "origin", "no-referrer-when-downgrade",
	"origin-when-cross-origin", "same-origin", "strict-origin", "strict-origin-when-cross-origin", "unsafe-URL"}

// Robots values for the content attribute of <meta name="robots">
var Robots = []string{"index", "noindex", "follow", "nofollow", "none", "noarchive", "nosnippet", "noimageindex", "nocache"}

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

// Cmd handles the target output command arguments.
func (args *Args) Cmd(data []byte, value string) {
	var err error
	switch {
	case args.SaveToFile:
		// use config save directory
		dir := []string{viper.GetString("save-directory")}
		if args.Dest, err = Dest(dir); err != nil {
			log.Fatal(err)
		}
		err = args.Save(&data)
	case !args.HTTP:
		// print to terminal
		err = args.Stdout(&data)
	}
	if err != nil {
		log.Fatal(err)
	}
}

// Dest determines if user supplied arguments are a valid file or directory destination.
func Dest(args []string) (path string, err error) {
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

// Save creates and saves the html template to the Dest argument.
func (args Args) Save(data *[]byte) error {
	name := filesystem.DirExpansion(args.Dest)
	stat, err := os.Stat(name)
	if err != nil {
		return fmt.Errorf("%s %q", err, name)
	}
	if stat.IsDir() {
		name = filepath.Join(name, "index.html")
	}
	color.OpFuzzy.Printf("Saving to %s\n", name)
	if !args.OW && !os.IsNotExist(err) {
		e := logs.Err{Issue: "html file exists", Arg: name, Msg: errors.New("include an -o flag to overwrite")}
		logs.ChkErr(e)
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
	d, err := args.pagedata(data)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(file, d); err != nil {
		return err
	}
	return nil
}

// Stdout creates and sends the html template to stdout.
func (args Args) Stdout(data *[]byte) error {
	tmpl, err := args.newTemplate()
	if err != nil {
		return err
	}
	d, err := args.pagedata(data)
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

// filename creates a filepath for the template filenames.
func (args Args) filename() (path string, err error) {
	base := "static/html/"
	if args.Test {
		base = filepath.Join("../../", base)
	}
	f := createTemplates()[args.Layout]
	if f == "" {
		return path, errors.New("filename: invalid-layout")
	}
	path = filepath.Join(base, f+".html")
	return path, err
}

// newTemplate creates and parses a new template file.
// The argument test is used internally.
func (args Args) newTemplate() (*template.Template, error) {
	fn, err := args.filename()
	if err != nil {
		return nil, err
	}
	fn, err = filepath.Abs(fn)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return nil, fmt.Errorf("create newTemplate: %s", err)
	}
	t := template.Must(template.ParseFiles(fn))
	return t, nil
}

// pagedata creates the meta and page template data.
func (args Args) pagedata(b *[]byte) (p PageData, err error) {
	// templates are found in the dir static/html/*.html
	switch args.Layout {
	case "full", "standard":
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
	case "mini":
		p.PageTitle = viper.GetString("html.title")
		p.MetaGenerator = false
	}
	// convert bytes into utf8
	var name = args.Enc
	if name == "" {
		name = "cp437"
	}
	runes, err := convert.Chars(name, b)
	logs.Check("create.pagedata.chars", err)
	p.PreText = string(runes)
	return p, nil
}
