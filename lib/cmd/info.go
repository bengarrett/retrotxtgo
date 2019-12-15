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
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/aofei/mimesniffer"
	"github.com/bengarrett/retrotxtgo/lib/codepage"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	humanize "github.com/labstack/gommon/bytes"
	"github.com/mattn/go-runewidth"
	"github.com/mozillazg/go-slugify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/gookit/color.v1"
)

const infoFormats string = "color, json, json.min, text, xml"

//Detail of a file
type Detail struct {
	Bytes     int64
	CharCount int
	Name      string
	MD5       string
	Mime      string
	Modified  time.Time
	SHA256    string
	Slug      string
	Size      string
	Utf8      bool
}

var (
	infoFmt  string
	fileName string
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information on a text file",
	Run: func(cmd *cobra.Command, args []string) {
		// only show Usage() with no errors if no flags .NFlags() are set
		if fileName == "" && cmd.Flags().NFlag() == 0 {
			fmt.Printf("%s\n\n", cmd.Short)
			_ = cmd.Usage()
			os.Exit(0)
		}
		if fileName == "" {
			_ = cmd.Usage()
			FileMissingErr()
		}
		f, err := details(fileName)
		Check(ErrorFmt{"file is invalid", fileName, err})
		switch viper.GetString("info.format") {
		case "color", "c":
			fmt.Printf("%s", infoText(true, f))
		case "json", "j":
			fmt.Printf("%s\n", infoJSON(true, f))
		case "json.min", "jm":
			fmt.Printf("%s\n", infoJSON(false, f))
		case "text":
			fmt.Printf("%s", infoText(false, f))
		case "xml", "x":
			fmt.Printf("%s\n", infoXML(f))
		default:
			CheckFlag(ErrorFmt{"format", infoFmt, fmt.Errorf(infoFormats)})
		}
	},
}

func init() {
	InitDefaults()
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&fileName, "name", "n", "", cp("text file to analyse")+" (required)\n")
	infoCmd.Flags().StringVarP(&infoFmt, "format", "f", viper.GetString("info.format"), "output format \noptions: "+ci(infoFormats))
	_ = viper.BindPFlag("info.format", infoCmd.Flags().Lookup("format"))
	_ = infoCmd.MarkFlagFilename("file")
	_ = infoCmd.MarkFlagRequired("file")
	infoCmd.Flags().SortFlags = false
}

func infoJSON(indent bool, f Detail) []byte {
	var j []byte
	var err error
	switch indent {
	case true:
		j, err = json.MarshalIndent(f, "", "    ")
	default:
		j, err = json.Marshal(f)
	}
	Check(ErrorFmt{"could not create", "json", err})
	return j
}

func infoText(c bool, f Detail) string {
	color.Enable = c
	var info = func(t string) string {
		return cinf(fmt.Sprintf("%s\t", t))
	}
	var hr = func() string {
		return fmt.Sprintf("\t%s\n", cf(strings.Repeat("\u2015", 26)))
	}
	var data = []struct {
		d, v string
	}{
		{d: "filename", v: f.Name},
		{d: "UTF-8", v: fmt.Sprintf("%v", f.Utf8)},
		{d: "characters", v: fmt.Sprintf("%v", f.CharCount)},
		{d: "size", v: f.Size},
		{d: "modified", v: fmt.Sprintf("%v", f.Modified.Format(FileDate))},
		{d: "MD5 checksum", v: f.MD5},
		{d: "SHA256 checksum", v: f.SHA256},
		{d: "MIME type", v: f.Mime},
		{d: "slug", v: f.Slug},
	}
	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 0, 8, 0, '\t', 0)
	fmt.Fprint(w, hr())
	for _, x := range data {
		fmt.Fprintf(w, "\t %s\t  %s\n", x.d, info(x.v))
	}
	fmt.Fprint(w, hr())
	w.Flush()
	return buf.String()
}

func infoXML(f Detail) []byte {
	type xmldetail struct {
		XMLName   xml.Name  `xml:"file"`
		ID        string    `xml:"id,attr"`
		Name      string    `xml:"name"`
		Mime      string    `xml:"content>mime"`
		Utf8      bool      `xml:"content>utf8"`
		Bytes     int64     `xml:"size>bytes"`
		Size      string    `xml:"size>value"`
		CharCount int       `xml:"size>character-count"`
		MD5       string    `xml:"checksum>md5"`
		SHA256    string    `xml:"checksum>sha256"`
		Modified  time.Time `xml:"modified"`
	}
	x := xmldetail{}
	x.Bytes = f.Bytes
	x.CharCount = f.CharCount
	x.ID = f.Slug
	x.MD5 = f.MD5
	x.Mime = f.Mime
	x.Modified = f.Modified
	x.Name = f.Name
	x.SHA256 = f.SHA256
	x.Size = f.Size
	x.Utf8 = f.Utf8
	xmlData, err := xml.MarshalIndent(x, "", "\t")
	Check(ErrorFmt{"could not create", "xml", err})
	return xmlData
}

func details(name string) (Detail, error) {
	d := Detail{}
	// Get the file details
	stat, err := os.Stat(name)
	if err != nil {
		return d, err
	}
	// Read file content
	data, err := filesystem.ReadAllBytes(name)
	if err != nil {
		return d, err
	}
	return parse(data, stat)
}

func parse(data []byte, stat os.FileInfo) (Detail, error) {
	md5sum := md5.Sum(data)
	sha256 := sha256.Sum256(data)
	mime := mimesniffer.Sniff(data)
	// create a table of data
	d := Detail{}
	d.Bytes = stat.Size()
	d.CharCount = runewidth.StringWidth(string(data))
	d.Name = stat.Name()
	d.MD5 = fmt.Sprintf("%x", md5sum)
	d.Modified = stat.ModTime()
	d.Slug = slugify.Slugify(stat.Name())
	d.SHA256 = fmt.Sprintf("%x", sha256)
	d.Utf8 = codepage.UTF8(data)
	if stat.Size() < 1000 {
		d.Size = fmt.Sprintf("%v bytes", stat.Size())
	} else {
		d.Size = fmt.Sprintf("%v (%v bytes)", humanize.Format(stat.Size()), stat.Size())
	}
	if strings.Contains(mime, ";") {
		d.Mime = strings.Split(mime, ";")[0]
	} else {
		d.Mime = mime
	}
	return d, nil
}
