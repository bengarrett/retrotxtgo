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
	"log"
	"os"

	"github.com/bengarrett/retrotxtgo/encoding"
	"github.com/labstack/gommon/bytes"
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

		fmt.Printf("Filename:\t\t%s\n", stat.Name())
		if stat.Size() < 1000 {
			fmt.Printf("Size:\t\t\t%v bytes\n", stat.Size())
		} else {
			fmt.Printf("Size:\t\t\t%v (%v bytes)\n", bytes.Format(stat.Size()), stat.Size())
		}
		fmt.Printf("Modified:\t\t%v\n", stat.ModTime())
		fmt.Printf("UTF-8 encoded:\t\t%v\n", encoding.IsUTF8(path))
		fmt.Printf("MD5 checksum:\t\t%x\n", sum)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
