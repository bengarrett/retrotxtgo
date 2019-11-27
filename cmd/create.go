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
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/filesystem"
	"github.com/spf13/cobra"
)

//

type files map[string]string

var (
	HTMLLayout      string
	MetaAuthor      string
	MetaColorScheme string
	MetaDesc        string
	MetaGenerator   bool
	MetaKeywords    string
	MetaReferrer    string
	MetaThemeColor  string
	PageTitle       string
	PreText         string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [FILE]",
	Short: "Create a HTML document from a text file",
	Example: `  retrotxtgo create textfile.txt -t "Text file" -d "Some random text file"
  retrotxtgo create ~/Downloads/textfile.txt --layout mini`,
	//todo make --layout a cobra.Command.ValidArgs
	//replace with a -f --file=""
	// ValidArgs: strings.Split(layOpts(), " "),
	// Args:      cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var data []byte
		var err error
		b := cmd.Flags().Lookup("body")
		if b.Changed == true {
			data = []byte(fmt.Sprintf("%s", b.Value))
		} else {
			if len(args) == 0 {
				err = fmt.Errorf("the [FILE] argument pointing to a textfile is required and missing")
				log.Fatal(err)
				os.Exit(1)
			}
			data, err = read(args[0])
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		}
		err = toStdout(data, false)
		if err != nil {
			log.Fatal(err) //Todo: handle differently, ie no log
			os.Exit(1)
		}
	},
}

// docs: https://godoc.org/github.com/spf13/pflag
func init() {
	d := LayoutDefault()
	rootCmd.AddCommand(createCmd)
	// main flags
	createCmd.Flags().StringVarP(&HTMLLayout, "layout", "l", "standard", "output HTML layout\noptions: "+layOpts())
	createCmd.Flags().StringVarP(&PageTitle, "title", "t", d.PageTitle, "defines the page title that is shown in a browser title bar or tab")
	createCmd.Flags().StringVarP(&MetaDesc, "meta-description", "d", d.MetaDesc, "a short and accurate summary of the content of the page")
	createCmd.Flags().StringVarP(&MetaAuthor, "meta-author", "a", d.MetaAuthor, "defines the name of the page authors")
	// minor flags
	createCmd.Flags().BoolVarP(&MetaGenerator, "meta-generator", "g", d.MetaGenerator, "include the RetroTxt version and page generation date")
	createCmd.Flags().StringVar(&MetaColorScheme, "meta-color-scheme", d.MetaColorScheme, "specifies one or more color schemes with which the page is compatible")
	createCmd.Flags().StringVar(&MetaKeywords, "meta-keywords", d.MetaKeywords, "words relevant to the page content")
	createCmd.Flags().StringVar(&MetaReferrer, "meta-referrer", d.MetaReferrer, "controls the Referer HTTP header attached to requests sent from the page")
	createCmd.Flags().StringVar(&MetaThemeColor, "meta-theme-color", d.MetaThemeColor, "indicates a suggested color that user agents should use to customize the display of the page")
	// hidden flags
	createCmd.Flags().StringVarP(&PreText, "body", "b", "", "override and inject string content into the body element")
	// flag options
	createCmd.Flags().MarkHidden("body")
	createCmd.Flags().SortFlags = false
}

// layOpts lists the options permitted by the layout flag.
func layOpts() string {
	h := layTemplates()
	s := []string{}
	for key := range h {
		s = append(s, key)
	}
	return strings.Join(s, ", ")
}

// layTemplates creates a map of the template filenames used in conjunction with the layout flag.
func layTemplates() files {
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
	f := layTemplates()[HTMLLayout]
	if f == "" {
		return "", fmt.Errorf("invalid flag value for --layout %s\nvalid options: %v", HTMLLayout, layOpts())
	}
	path += f + ".html"
	return path, nil
}

// pagedata creates the meta and page template data.
// todo handle all arguments
func pagedata(data []byte) PageData {
	var p PageData
	switch HTMLLayout {
	case "full":
	case "standard":
		p = LayoutDefault()
		p.MetaAuthor = MetaAuthor
		p.MetaGenerator = MetaGenerator
	case "mini":
		p.MetaGenerator = false
	}
	p.PreText = string(data)
	return p
}

// toStdout creates and sends the html template to stdout.
// The argument test is used internally.
func toStdout(data []byte, test bool) error {
	fn, err := filename(test)
	if err != nil {
		return err
	}
	t := template.Must(template.ParseFiles(fn))
	err = t.Execute(os.Stdout, pagedata(data))
	if err != nil {
		return err
	}
	return nil
}
