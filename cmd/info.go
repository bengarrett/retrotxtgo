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
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/InVisionApp/tabular"
	"github.com/aofei/mimesniffer"
	"github.com/bengarrett/retrotxtgo/encoding"
	"github.com/bengarrett/retrotxtgo/filesystem"
	"github.com/labstack/gommon/bytes"
	"github.com/mattn/go-runewidth"
	"github.com/mozillazg/go-slugify"
	"github.com/spf13/cobra"
)

//Detail of a file
type Detail struct {
	Bytes     int64
	CharCount int
	Name      string
	MD5       string
	Mime      string
	Modified  time.Time
	Slug      string
	Size      string
	Utf8      bool
}

var (
	//Output format flag
	Output string
)

// paths := [5]string{
// 	"textfiles/hi.txt",
// 	"/Users/ben/Downloads/impure74/jp!xqtrd.asc",
// 	"/Users/ben/Downloads/impure74/impure74.ans",
// 	"/Users/ben/Downloads/bbh/hx_joker2019.ans",
// 	"/Users/ben/Downloads/bbh/hx_jack.ans",
// }

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info FILE",
	Short: "Information on a text file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := details(args[0])
		if err != nil {
			log.Fatal(err)
		}
		switch Output {
		case "json":
			jsonData, err := json.MarshalIndent(f, "", "    ")
			if err != nil {
				log.Println(err)
			}
			fmt.Println(string(jsonData))
		case "json.min":
			jsonData, err := json.Marshal(f)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(string(jsonData))
		case "table":
		case "":
			table(f)
		default:
			fmt.Printf("invalid flag, --output %s\noptions: json, json.min, table\n", Output)
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&Output, "output", "o", "", "output format (default is table)\noptions: json, json.min, table")
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
	checksum := md5.Sum(data)
	mime := mimesniffer.Sniff(data)
	// create a table of data
	d.Bytes = stat.Size()
	d.CharCount = runewidth.StringWidth(string(data))
	d.Name = stat.Name()
	d.MD5 = fmt.Sprintf("%x", checksum)
	d.Modified = stat.ModTime()
	d.Slug = slugify.Slugify(stat.Name())
	d.Utf8 = encoding.IsUTF8(data)
	if stat.Size() < 1000 {
		d.Size = fmt.Sprintf("%v bytes", stat.Size())
	} else {
		d.Size = fmt.Sprintf("%v (%v bytes)", bytes.Format(stat.Size()), stat.Size())
	}
	if strings.Contains(mime, ";") {
		d.Mime = strings.Split(mime, ";")[0]
	} else {
		d.Mime = mime
	}
	return d, nil
}

func table(d Detail) {
	tab := tabular.New()
	tab.Col("det", "Details", 14)
	tab.Col("val", "", 10)
	var data = []struct {
		d, v string
	}{
		{d: "Filename", v: d.Name},
		{d: "MIME type", v: d.Mime},
		{d: "UTF-8", v: fmt.Sprintf("%v", d.Utf8)},
		{d: "Characters", v: fmt.Sprintf("%v", d.CharCount)},
		{d: "Size", v: d.Size},
		{d: "Modified", v: fmt.Sprintf("%v", d.Modified)},
		{d: "MD5 checksum", v: d.MD5},
		{d: "Slug", v: d.Slug},
	}
	format := tab.Print("*")
	for _, x := range data {
		fmt.Printf(format, x.d, x.v)
	}
}
