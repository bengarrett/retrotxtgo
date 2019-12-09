/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/codepage"

	"golang.org/x/text/encoding/japanese"

	"github.com/bengarrett/retrotxtgo/filesystem"

	"github.com/spf13/cobra"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
)

const viewFormats string = "color, text"

type iana struct {
	mime  string
	index string
	mib   string
	s     []string
}

var (
	viewCodePage string = "ibm437"
	viewFilename string
	viewFormat   string
	viewWidth    int
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// todo reverse scan of file looking for SAUCE00 and COMNTT
		// todo scan for unique color codes like 24-bit color
		// todo scan for new lines or character counts and hard-code the width

		if viewFilename == "" {
			// todo: remove from final
			viewFilename = "textfiles/cp-437-all-characters.txt"
		}
		// todo handle unchanged viewCodePage, where UTF8 encoding will be checked otherwise use
		encoding, err := ianaindex.IANA.Encoding(viewCodePage)
		// todo handle invalid encoding value, with notice on how to display the list
		Check(ErrorFmt{"encoding transform", viewCodePage, err})
		var d codepage.Set

		data, err := filesystem.Read(viewFilename)
		Check(ErrorFmt{"file open", viewFilename, err})

		d.Transform(data, encoding)
		d.SwapAll(true)
		fmt.Printf("\n%s\n", d.Data)

		println(codepage.Table(""))

		// fmt.Printf("\n%s\n", data)
		// // todo: make an --example that auto generates a table bytes 0 - 255 | lf every 16 characters
		// oof := []byte{48, 49, 50, 51, 52}
		// fmt.Printf("\n%v\t%d\n", oof, oof)
	},
}

var viewCodePagesCmd = &cobra.Command{
	Use:   "codepages",
	Short: "list available legacy codepages that RetroTxt can convert into UTF-8",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(codepages())
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVarP(&viewFilename, "name", "n", "", cp("text file to display")+" (required)\n")
	viewCmd.Flags().StringVarP(&viewCodePage, "codepage", "c", "cp437", "legacy character encoding used by the text file")
	viewCmd.Flags().StringVarP(&viewFormat, "format", "f", "color", "output format, options: "+ci(viewFormats))
	viewCmd.Flags().IntVarP(&viewWidth, "width", "w", 80, "document column character width")
	// override ascii 0-F + 1-F || Control characters || IBM, ASCII, IBM+
	// example flag showing CP437 table
	viewCmd.MarkFlagFilename("name")
	viewCmd.MarkFlagRequired("name")
	viewCmd.Flags().SortFlags = false
	viewCmd.AddCommand(viewCodePagesCmd)
}

// codepages returns a tabled list of supported IANA character set encodings
func codepages() string {
	// create a buffer and writer for tab formatting
	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 0, 8, 0, '\t', 0)

	var ii iana
	var err error
	fmt.Fprintln(w, cp("\nSupported legacy codepages and encodings"))
	fmt.Fprintln(w, cf(strings.Repeat("\u2015", 40)))
	fmt.Fprintf(w, "\ttitle\talias(s)\n")
	c := append(charmap.All, japanese.All...)
	for _, n := range c {
		name := fmt.Sprint(n)
		if name == "X-User-Defined" {
			continue
		}
		ii.mime, err = ianaindex.MIME.Name(n)
		if err != nil {
			continue
		}
		ii.index, _ = ianaindex.IANA.Name(n)
		ii.mib, _ = ianaindex.MIB.Name(n)
		ii.s = strings.Split(name, " ")
		// display encoding name and alias
		fmt.Fprintf(w, "\t%s\t%s", name, ci(ii.mib))
		// create common use CP aliases
		switch {
		case ii.s[0] == "IBM":
		case ii.s[0] == "Windows" && ii.s[1] == "Code":
			fmt.Fprintf(w, "\tCP%s", ii.s[3])
		case ii.s[0] == "Windows":
			fmt.Fprintf(w, "\tCP%s", ii.s[1])
		default:
			fmt.Fprintf(w, "\t")
		}
		// only show MIME if it is different to the previous aliases
		switch {
		case strings.ToLower(strings.ReplaceAll(name, "-", " ")) == strings.ToLower(strings.ReplaceAll(ii.mime, "-", " ")):
		case strings.ReplaceAll(name, "-", "") == strings.ReplaceAll(ii.mime, "-", ""):
		case ii.mib == ii.mime:
			fmt.Fprintf(w, "\t%s", cf(""))
		default:
			fmt.Fprintf(w, "\t%s", cf(ii.mime))
		}
		fmt.Fprint(w, "\n")
	}
	fmt.Fprint(w, "\ttitle\talias(s)\n")
	fmt.Fprint(w, "\n"+cinf("*")+" Code Page 437 ("+cc("CP437")+") is commonly used by MS-DOS English text and ANSI art")
	fmt.Fprint(w, "\n  ISO 8859-1 ("+cc("ISOLatin1")+") is found in legacy Internet, Unix and Amiga documents")
	fmt.Fprint(w, "\n  Windows 1252 ("+cc("windows1252")+") is found in legacy English language Windows operating systems")
	w.Flush()
	return buf.String()
}

// transform opens a text file and converts it from the Encoding legacy character set into UTF-8.
func transform(name string, e encoding.Encoding) codepage.Set {
	var s codepage.Set
	data, err := filesystem.Read(name)
	Check(ErrorFmt{"file open", name, err})
	iana, err := ianaindex.IANA.Name(e)
	if err != nil {
		iana = "unknown"
	}
	Check(ErrorFmt{"IANA index", iana, err})
	decode, err := e.NewDecoder().Bytes(data)
	Check(ErrorFmt{"encoding transform", iana, err})
	s.Data = decode
	return s
}
