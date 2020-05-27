package create

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding/charmap"
)

// func X( data []byte, options)

type files map[string]string

// Args holds arguments and options.
type Args struct {
	HTMLLayout  string
	ServerFiles bool
	ServerPort  int
	Styles      string
	Test        bool
}

// PageData holds template data used by the HTML layouts.
type PageData struct {
	BuildVersion    string
	BuildDate       time.Time
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

// File creates and saves the html template to the named file.
func (args Args) File(data []byte, name string) error {
	if name == "~" {
		// allow the use of ~ as the home directory on Windows
		u, err := user.Current()
		if err != nil {
			return err
		}
		name = u.HomeDir
	}
	stat, err := os.Stat(name)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		name = path.Join(name, "index.html")
	}
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	tmpl, err := args.newTemplate(args.Test)
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

// Save ...
func (args Args) Save(data []byte, value string, changed bool) {
	var err error
	switch {
	case changed:
		err = args.File(data, value)
	case viper.GetString("create.save-directory") != "":
		err = args.File(data, viper.GetString("create.save-directory"))
	case !args.ServerFiles:
		err = args.Stdout(data, false)
	}
	if err != nil {
		// TODO: handle errors
	}
}

// Serve ...
func (args Args) Serve(data []byte) {
	p := uint(args.ServerPort)
	if !prompt.PortValid(p) {
		// viper.GetInt() doesn't work as expected
		port, err := strconv.Atoi(viper.GetString("create.server-port"))
		if err != nil {
			logs.Check("create serve port", err)
		}
		p = uint(port)
	}
	if err := args.serveFile(data, p, false); err != nil {
		logs.ChkErr(logs.Err{Issue: "server problem", Arg: "HTTP", Msg: err})
	}
}

// serveFile creates and serves the html template on a local HTTP web server.
// The argument test is used internally.
func (args Args) serveFile(data []byte, port uint, test bool) error {
	t, err := args.newTemplate(test)
	if err != nil {
		return err
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		d, err := args.pagedata(data)
		if err != nil {
			logs.ChkErr(logs.Err{Issue: "servefile encoding", Arg: "http", Msg: err})
		}
		if err = t.Execute(w, d); err != nil {
			logs.ChkErr(logs.Err{Issue: "servefile", Arg: "http", Msg: err})
		}
	})
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Printf("Web server is available at %s\n",
		str.Cp(fmt.Sprintf("http://localhost:%v", port)))
	println(str.Cinf("Press Ctrl+C to stop"))
	if err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		return err
	}
	return nil
}

// Stdout creates and sends the html template to stdout.
// The argument test is used internally.
func (args Args) Stdout(data []byte, test bool) error {
	tmpl, err := args.newTemplate(test)
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
	switch args.Styles {
	case "", "none":
		fmt.Printf("%s", buf.String())
	default:
		if err = str.Highlight(buf.String(), "html", args.Styles); err != nil {
			return err
		}
	}
	return err
}

// Options returns the options permitted by the layout flag.
func Options() (s []string) {
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
// The argument test is used internally.
func (args Args) filename(test bool) (path string, err error) {
	base := "static/html/"
	if test {
		base = filepath.Join("../../", base)
	}
	f := createTemplates()[args.HTMLLayout]
	if f == "" {
		return path, errors.New("filename: invalid-layout")
	}
	path = filepath.Join(base, f+".html")
	return path, err
}

// newTemplate creates and parses a new template file.
// The argument test is used internally.
func (args Args) newTemplate(test bool) (*template.Template, error) {
	fn, err := args.filename(test)
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
// todo handle all arguments
func (args Args) pagedata(data []byte) (p PageData, err error) {
	switch args.HTMLLayout {
	case "full", "standard":
		p.MetaAuthor = viper.GetString("create.meta.author")
		p.MetaColorScheme = viper.GetString("create.meta.color-scheme")
		p.MetaDesc = viper.GetString("create.meta.description")
		p.MetaGenerator = viper.GetBool("create.meta.generator")
		p.MetaKeywords = viper.GetString("create.meta.keywords")
		p.MetaReferrer = viper.GetString("create.meta.referrer")
		p.MetaThemeColor = viper.GetString("create.meta.theme-color")
		p.PageTitle = viper.GetString("create.title")
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

func transform(m *charmap.Charmap, p []byte) (runes int, encoded []byte, err error) {
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
