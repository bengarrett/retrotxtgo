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
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/aofei/mimesniffer"
	"github.com/bengarrett/retrotxtgo/encoding"
	"github.com/bengarrett/retrotxtgo/sauce"
	"github.com/labstack/gommon/bytes"
	"github.com/mattn/go-runewidth"
	"github.com/mozillazg/go-slugify"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("info called")
		//path := "textfiles/hi.txt"
		path := "/Users/ben/Downloads/impure74/jp!xqtrd.asc"
		//path := "/Users/ben/Downloads/impure74/impure74.ans"
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		// Open file for reading
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		// Create new hasher, which is a writer interface
		hasher := md5.New()
		_, err = io.Copy(hasher, file)
		if err != nil {
			log.Fatal(err)
		}
		// Hash and print. Pass nil since
		// the data is not coming in as a slice argument
		// but is coming through the writer interface
		sum := hasher.Sum(nil)

		// Open file for reading
		file, err = os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Filename:\t\t%s\n", stat.Name())
		if stat.Size() < 1000 {
			fmt.Printf("Size:\t\t\t%v bytes\n", stat.Size())
		} else {
			fmt.Printf("Size:\t\t\t%v (%v bytes)\n", bytes.Format(stat.Size()), stat.Size())
		}
		fmt.Printf("Modified:\t\t%v\n", stat.ModTime())
		fmt.Printf("UTF-8 encoded:\t\t%v\n", encoding.IsUTF8(path))
		fmt.Printf("MD5 checksum:\t\t%x\n", sum)
		fmt.Printf("Slug:\t\t\t%v\n", slugify.Slugify(stat.Name()))
		fmt.Printf("Width:\t\t\t%v (not working)\n\n", runewidth.StringWidth(string(data)))
		fmt.Printf("MIME type:\t\t%v\n", mimesniffer.Sniff(data)) // todo slice first 512bytes

		fmt.Printf("SAUCE metadata:\t\t%v\n", sauce.Scan(data))

		s := sauce.Get(data)

		fmt.Printf("\nversion:\t\t%v\ntitle:\t\t%q\nauthor:\t\t%q\n", s.Version, s.Title, s.Author)
		fmt.Printf("group:\t\t%q\ndate:\t\t%s\nlsd:\t\t%s\n", s.Group, s.Date, s.LSDate)
		fmt.Printf("file size:\t%v\n", s.FileSize)
		fmt.Printf("data type:\t%q\n", s.DataType)
		fmt.Printf("file type:\t%q\n", s.FileType)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
