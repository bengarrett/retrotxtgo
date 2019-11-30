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
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/bengarrett/retrotxtgo/filesystem"
	"github.com/spf13/cobra"
	"gopkg.in/gookit/color.v1"
)

type files map[string]string

// create command flag
var (
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
)

// createCmd makes create usage examples
var exampleCmd = func() string {
	s := string(os.PathSeparator)
	e := `  retrotxtgo create textfile.txt -t "Text file" -d "Some random text file"`
	e += fmt.Sprintf("\n  retrotxtgo create ~%sDownloads%stextfile.txt --layout mini", s, s)
	e += fmt.Sprintf("\n  retrotxtgo create textfile.txt -s .%shtml", s)
	return color.Info.Sprint(e)
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create FILE",
	Short: color.Primary.Sprint("Create a HTML document from a text file"),
	//Long: `` // used by help create
	Example: exampleCmd(),
	Run: func(cmd *cobra.Command, args []string) {
		var data []byte
		var err error
		// --body="" is a hidden flag to test without providing a FILE
		b := cmd.Flags().Lookup("body")
		switch b.Changed {
		case true:
			data = []byte(fmt.Sprintf("%s", b.Value))
		default:
			if len(args) == 0 {
				FileMissingErr()
			}
			data, err = read(args[0])
			if err != nil {
				h := ErrorFmt{"invalid FILE", args[0], err}
				h.GoErr()
			}
		}
		// check for a --save flag to save to a file
		// otherwise output is sent to stdout
		s := cmd.Flags().Lookup("save")
		switch s.Changed {
		case true:
			err = writeFile(data, fmt.Sprintf("%s", s.Value), false)
		default:
			err = writeStdout(data, false)
		}
		if err != nil {
			if err.Error() == errors.New("invalid-layout").Error() {
				h := ErrorFmt{"layout", fmt.Sprintf("%s", htmlLayout), fmt.Errorf(createLayouts())}
				h.FlagErr()
			}
			h := ErrorFmt{"create error", ">", err}
			h.GoErr()
		}
	},
}

func init() {
	homedir := func() string {
		s := "\n--save ~ saves to the home or user directory"
		d, err := os.UserHomeDir()
		if err != nil {
			return s
		}
		return s + " at " + d
	}
	curdir := func() string {
		s := "\n--save . saves to the current working directory"
		d, err := os.Getwd()
		if err != nil {
			return s
		}
		return s + " at " + d
	}

	d := LayoutDefault()
	rootCmd.AddCommand(createCmd)
	// main flags
	createCmd.Flags().StringVarP(&htmlLayout, "layout", "l", "standard", "output HTML layout\noptions: "+createLayouts())
	createCmd.Flags().StringVarP(&pageTitle, "title", "t", d.PageTitle, "defines the page title that is shown in a browser title bar or tab")
	createCmd.Flags().StringVarP(&metaDesc, "meta-description", "d", d.MetaDesc, "a short and accurate summary of the content of the page")
	createCmd.Flags().StringVarP(&metaAuthor, "meta-author", "a", d.MetaAuthor, "defines the name of the page authors")
	// minor flags
	createCmd.Flags().BoolVarP(&metaGenerator, "meta-generator", "g", d.MetaGenerator, "include the RetroTxt version and page generation date")
	createCmd.Flags().StringVar(&metaColorScheme, "meta-color-scheme", d.MetaColorScheme, "specifies one or more color schemes with which the page is compatible")
	createCmd.Flags().StringVar(&metaKeywords, "meta-keywords", d.MetaKeywords, "words relevant to the page content")
	createCmd.Flags().StringVar(&metaReferrer, "meta-referrer", d.MetaReferrer, "controls the Referer HTTP header attached to requests sent from the page")
	createCmd.Flags().StringVar(&metaThemeColor, "meta-theme-color", d.MetaThemeColor, "indicates a suggested color that user agents should use to customize the display of the page")
	// hidden flags
	createCmd.Flags().StringVarP(&preText, "body", "b", "", "override and inject string content into the body element")
	createCmd.Flags().StringVarP(&saveToFiles, "save", "s", "", "save HTML as files to store this directory"+homedir()+curdir())
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

// read opens and returns the content of the name file.
func read(name string) ([]byte, error) {
	// check name is file not anything else
	data, err := filesystem.ReadAllBytes(name)
	if err != nil {
		return nil, err
	}
	return data, nil
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
		p = LayoutDefault()
		p.MetaAuthor = metaAuthor
		p.MetaColorScheme = metaColorScheme
		p.MetaDesc = metaDesc
		p.MetaGenerator = metaGenerator
		p.MetaKeywords = metaKeywords
		p.MetaReferrer = metaReferrer
		p.MetaThemeColor = metaThemeColor
		p.PageTitle = pageTitle
	case "mini":
		p.PageTitle = pageTitle
		p.MetaGenerator = false
	}
	p.PreText = string(data)
	return p
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
	fn, err := filename(test)
	if err != nil {
		return err
	}
	t := template.Must(template.ParseFiles(fn))
	if err = t.Execute(f, pagedata(data)); err != nil {
		return err
	}
	return nil
}

// writeStdout creates and sends the html template to stdout.
// The argument test is used internally.
func writeStdout(data []byte, test bool) error {
	fn, err := filename(test)
	if err != nil {
		return err
	}
	t := template.Must(template.ParseFiles(fn))
	if err = t.Execute(os.Stdout, pagedata(data)); err != nil {
		return err
	}
	return nil
}
