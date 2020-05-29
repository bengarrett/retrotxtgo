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
	"unicode/utf8"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/lib/version"
	"github.com/gookit/color"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding/charmap"
)

type files map[string]string

// Args holds arguments and options sourced from user flags or the config file.
type Args struct {
	// Src source input text file
	Src string
	// Dest HTML destination either a directory or file
	Dest string
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
	// Referrer metadata
	Ref string
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
	MetaReferrer    string
	MetaThemeColor  string
	PageTitle       string
	PreText         string
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
		dir := []string{viper.GetString("create.save-directory")}
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

// Save creates and saves the html template to the named file.
func (args Args) Save(data *[]byte) error {
	name := args.Dest
	stat, err := os.Stat(name)
	if err != nil {
		return err
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
// TODO: handle all arguments
func (args Args) pagedata(data *[]byte) (p PageData, err error) {
	// templates are found in the dir static/html/*.html
	switch args.Layout {
	case "full", "standard":
		p.MetaAuthor = viper.GetString("create.meta.author")
		p.MetaColorScheme = viper.GetString("create.meta.color-scheme")
		p.MetaDesc = viper.GetString("create.meta.description")
		p.MetaGenerator = viper.GetBool("create.meta.generator")
		p.MetaKeywords = viper.GetString("create.meta.keywords")
		p.MetaReferrer = viper.GetString("create.meta.referrer")
		p.MetaThemeColor = viper.GetString("create.meta.theme-color")
		p.PageTitle = viper.GetString("create.title")
		/*
			https://webmasters.googleblog.com/2007/12/answering-more-popular-picks-meta-tags.html
			<meta name="google" value="notranslate">
			<span class="notranslate">
			<meta name="robots" content="…, …">
			viewport ?
		*/
		// generate data
		t := time.Now().UTC()
		p.BuildDate = t.Format(time.RFC3339)
		p.BuildVersion = version.B.Version
	case "mini":
		p.PageTitle = viper.GetString("create.title")
		p.MetaGenerator = false
	}
	// convert to utf8
	_, encoded, err := transform(nil, data)
	if err != nil {
		return p, err
	}
	p.PreText = string(encoded)
	return p, nil
}

func transform(m *charmap.Charmap, px *[]byte) (runes int, encoded []byte, err error) {
	p := *px
	if len(p) == 0 {
		return 0, encoded, nil
	}
	// confirm encoding is not utf8
	if utf8.Valid(p) {
		return utf8.RuneCount(p), p, nil
	}
	// use cp437 by default if text is not utf8
	// TODO: add default-unknown.encoding setting
	if m == nil {
		m = charmap.CodePage437
	}
	// convert to utf8
	if encoded, err = m.NewDecoder().Bytes(p); err != nil {
		return 0, encoded, err
	}
	return utf8.RuneCount(encoded), encoded, nil
}

/*
var encodings = []struct {
	name        string
	mib         string
	comment     string
	varName     string
	replacement byte
	mapping     string
}{
		"IBM Code Page 437",
		"PC8CodePage437",
		"",
		"CodePage437",
		encoding.ASCIISub,
		"http://source.icu-project.org/repos/icu/data/trunk/charset/data/ucm/glibc-IBM437-2.1.2.ucm",
	},
	{
		"Windows 1254",
		"Windows1254",
		"",
		"Windows1254",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-windows-1254.txt",
	},	{
		"Macintosh",
		"Macintosh",
		"",
		"Macintosh",
		encoding.ASCIISub,
		"http://encoding.spec.whatwg.org/index-macintosh.txt",
	},

*/
