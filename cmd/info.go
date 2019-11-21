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
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/net/html/charset"
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
		path := "textfiles/hix.txt"
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		f, _ := os.Open(path)
		b1 := make([]byte, 500)
		n1, err := f.Read(b1)
		f.Close()
		c, name, certain := charset.DetermineEncoding(b1, "text/plain")
		fmt.Printf("1: %v\t 2: %v\t 3: %v 4: %v\n", c, name, certain, n1)
		//		DetermineEncoding(content []byte, contentType string) (e encoding.Encoding, name string, certain bool)
		fmt.Printf("Filename:\t\t%s\n", stat.Name())
		fmt.Printf("Size:\t\t\t%v bytes\n", stat.Size())
		fmt.Printf("Modified:\t\t%v\n", stat.ModTime())
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
