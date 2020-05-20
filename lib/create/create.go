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
	"strings"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/viper"
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
	if err = tmpl.Execute(file, args.pagedata(data)); err != nil {
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
	if !logs.PortValid(p) {
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
		if err = t.Execute(w, args.pagedata(data)); err != nil {
			logs.ChkErr(logs.Err{Issue: "serveFile", Arg: "http", Msg: err})
		}
	})
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Printf("Web server is available at %s\n",
		logs.Cp(fmt.Sprintf("http://localhost:%v", port)))
	println(logs.Cinf("Press Ctrl+C to stop"))
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
	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, args.pagedata(data)); err != nil {
		return err
	}
	switch args.Styles {
	case "", "none":
		fmt.Printf("%s", buf.String())
	default:
		if err = quick.Highlight(os.Stdout, buf.String(), "html", "terminal256", args.Styles); err != nil {
			return err
		}
	}
	return err
}

// Layouts lists the options permitted by the layout flag.
func Layouts() string {
	s := []string{}
	for key := range createTemplates() {
		s = append(s, key)
	}
	return strings.Join(s, ", ")
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
func (args Args) pagedata(data []byte) (p PageData) {
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
	p.PreText = string(data)
	return p
}
