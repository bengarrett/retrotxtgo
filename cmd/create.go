/*
Copyright © 2019 Ben Garrett <code.by.ben@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/bengarrett/retrotxtgo/filesystem"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/gookit/color.v1"
)

type files map[string]string

// create command flag
var (
	createFileName  string
	createStyles    string
	htmlLayout      string
	metaAuthor      string
	metaColorScheme string
	metaDesc        string
	metaGenerator   bool
	metaKeywords    string
	metaReferrer    string
	metaThemeColor  string
	pageTitle       string
	preText         string
	saveToFiles     string
	serverFiles     bool
	serverPort      int
)

// createCmd makes create usage examples
var exampleCmd = func() string {
	s := string(os.PathSeparator)
	e := `  retrotxtgo create -n textfile.txt -t "Text file" -d "Some random text file"`
	e += fmt.Sprintf("\n  retrotxtgo create --name ~%sDownloads%stextfile.txt --layout mini --save .%shtml", s, s, s)
	return color.Info.Sprint(e)
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use: "create",
	//Aliases: []string{"new"},
	Short: cp("Create a HTML document from a text file"),
	//Long: `` // used by help create
	Example: exampleCmd(),
	Run: func(cmd *cobra.Command, args []string) {
		htmlLayout = viper.GetString("create.layout")
		var data []byte
		var err error
		// --body="" is a hidden flag to test without providing a FILE
		b := cmd.Flags().Lookup("body")
		switch b.Changed {
		case true:
			data = []byte(fmt.Sprintf("%s", b.Value))
		default:
			if createFileName == "" {
				fmt.Printf("%s\n\n", cmd.Short)
				cmd.Usage()
				os.Exit(0)
			}
			data, err = filesystem.Read(createFileName)
			Check(ErrorFmt{"file is invalid", createFileName, err})
		}
		// check for a --save flag to save to files
		// otherwise output is sent to stdout
		s := cmd.Flags().Lookup("save")
		switch {
		case s.Changed == true:
			err = writeFile(data, fmt.Sprintf("%s", s.Value), false)
		case viper.GetString("create.save-directory") != "":
			err = writeFile(data, fmt.Sprintf("%s", viper.GetString("create.save-directory")), false)
		case serverFiles == false:
			err = writeStdout(data, false)
		}
		if err != nil {
			if err.Error() == errors.New("invalid-layout").Error() {
				CheckFlag(ErrorFmt{"layout", fmt.Sprintf("%s", htmlLayout), fmt.Errorf(createLayouts())})
			}
			Check(ErrorFmt{"create error", ">", err})
		}
		// check for a --server flag to serve the HTML
		if serverFiles == true {
			// viper.GetInt() doesn't work as expected
			port, err := strconv.Atoi(viper.GetString("create.server-port"))
			if err != nil {
				port = serverPort
			}
			if err = serveFile(data, port, false); err == nil {
				Check(ErrorFmt{"server problem", "HTTP", err})
			}
		}
	},
}

func init() {
	InitDefaults()
	homedir := func() string {
		s := "\n" + ci("--save ~") + " saves to the home or user directory"
		d, err := os.UserHomeDir()
		if err != nil {
			return s
		}
		return s + " at " + cf(d)
	}
	curdir := func() string {
		s := "\n" + ci("--save .") + " saves to the current working directory"
		d, err := os.Getwd()
		if err != nil {
			return s
		}
		return s + " at " + cf(d)
	}
	def := func(s string) string {
		return viper.GetString(s)
	}
	rootCmd.AddCommand(createCmd)
	// required flags
	createCmd.Flags().StringVarP(&createFileName, "name", "n", "", cp("text file to parse")+" (required)\n")
	// main flags
	createCmd.Flags().StringVarP(&htmlLayout, "layout", "l", def("create.layout"), "output HTML layout\noptions: "+ci(createLayouts()))
	viper.BindPFlag("create.layout", createCmd.Flags().Lookup(("layout")))
	createCmd.Flags().StringVarP(&createStyles, "syntax-style", "c", "lovelace", "HTML syntax highligher, use "+ci("none")+" to disable")
	createCmd.Flags().StringVarP(&pageTitle, "title", "t", def("create.title"), "defines the page title that is shown in a browser title bar or tab")
	viper.BindPFlag("create.title", createCmd.Flags().Lookup("title"))
	createCmd.Flags().StringVarP(&metaDesc, "meta-description", "d", def("create.meta.description"), "a short and accurate summary of the content of the page")
	viper.BindPFlag("create.meta.description", createCmd.Flags().Lookup("meta-description"))
	createCmd.Flags().StringVarP(&metaAuthor, "meta-author", "a", def("create.meta.author"), "defines the name of the page authors")
	viper.BindPFlag("create.meta.author", createCmd.Flags().Lookup("meta-author"))
	// minor flags
	createCmd.Flags().BoolVarP(&metaGenerator, "meta-generator", "g", viper.GetBool("create.meta.generator"), "include the RetroTxt version and page generation date")
	viper.BindPFlag("create.meta.generator", createCmd.Flags().Lookup("meta-generator"))
	createCmd.Flags().StringVar(&metaColorScheme, "meta-color-scheme", def("create.meta.color-scheme"), "specifies one or more color schemes with which the page is compatible")
	viper.BindPFlag("create.meta.color-scheme", createCmd.Flags().Lookup("meta-color-scheme"))
	createCmd.Flags().StringVar(&metaKeywords, "meta-keywords", def("create.meta.keywords"), "words relevant to the page content")
	viper.BindPFlag("create.meta.keywords", createCmd.Flags().Lookup("meta-keywords"))
	createCmd.Flags().StringVar(&metaReferrer, "meta-referrer", def("create.meta.referrer"), "controls the Referer HTTP header attached to requests sent from the page")
	viper.BindPFlag("create.meta.referrer", createCmd.Flags().Lookup("meta-referrer"))
	createCmd.Flags().StringVar(&metaThemeColor, "meta-theme-color", def("create.meta.theme-color"), "indicates a suggested color that user agents should use to customize the display of the page")
	viper.BindPFlag("create.meta.theme-color", createCmd.Flags().Lookup("meta-theme-color"))
	// output flags
	// todo: when using save-directory config setting, there is no way to stdout using flags
	// instead add an output flag with print, file|save
	createCmd.Flags().StringVarP(&saveToFiles, "save", "s", def("create.save-directory"), "save HTML as files to store this directory"+homedir()+curdir())
	viper.BindPFlag("create.save-directory", createCmd.Flags().Lookup("save"))
	createCmd.Flags().BoolVarP(&serverFiles, "server", "p", false, "serve HTML over an internal web server")
	createCmd.Flags().IntVar(&serverPort, "port", viper.GetInt("create.server-port"), "port which the internet web server will listen")
	viper.BindPFlag("create.server-port", createCmd.Flags().Lookup("port"))
	// hidden flags
	createCmd.Flags().StringVarP(&preText, "body", "b", "", "override and inject string content into the body element")
	// flag options
	createCmd.Flags().MarkHidden("body")
	createCmd.Flags().SortFlags = false
}

// createLayouts lists the options permitted by the layout flag.
func createLayouts() string {
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
func filename(test bool) (string, error) {
	path := "static/html/"
	if test {
		path = "../" + path
	}
	f := createTemplates()[htmlLayout]
	if f == "" {
		return "", errors.New("invalid-layout")
	}
	path += f + ".html"
	return path, nil
}

// pagedata creates the meta and page template data.
// todo handle all arguments
func pagedata(data []byte) PageData {
	var p PageData
	switch htmlLayout {
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

// newTemplate creates and parses a new template file.
// The argument test is used internally.
func newTemplate(test bool) (*template.Template, error) {
	fn, err := filename(test)
	if err != nil {
		return nil, err
	}
	t := template.Must(template.ParseFiles(fn))
	return t, nil
}

// serveFile creates and serves the html template on a local HTTP web server.
// The argument test is used internally.
func serveFile(data []byte, port int, test bool) error {
	t, err := newTemplate(test)
	if err != nil {
		return err
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, pagedata(data))
	})
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Printf("Web server is available at %s\n", color.Primary.Sprintf("http://localhost:%v", port))
	color.Info.Println("Press Ctrl+C to stop")
	if err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		return err
	}
	return nil
}

// writeFile creates and saves the html template to the name file.
// The argument test is used internally.
func writeFile(data []byte, name string, test bool) error {
	p := name
	s, err := os.Stat(name)
	if err != nil {
		return err
	}
	if s.IsDir() {
		p = path.Join(p, "index.html")
	}
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	t, err := newTemplate(test)
	if err != nil {
		return err
	}
	if err = t.Execute(f, pagedata(data)); err != nil {
		return err
	}
	return nil
}

// writeStdout creates and sends the html template to stdout.
// The argument test is used internally.
func writeStdout(data []byte, test bool) error {
	t, err := newTemplate(test)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, pagedata(data)); err != nil {
		return err
	}
	switch createStyles {
	case "", "none":
		fmt.Printf(buf.String())
	default:
		if err = quick.Highlight(os.Stdout, buf.String(), "html", "terminal256", createStyles); err != nil {
			return err
		}
	}
	return nil
}
