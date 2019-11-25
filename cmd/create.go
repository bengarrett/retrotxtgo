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

	"github.com/bengarrett/retrotxtgo/filesystem"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a HTML document from a text file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		data, err := filesystem.ReadAllBytes("textfiles/hi.txt")
		if err != nil {
			log.Fatal(err)
		}
		tail, err := filesystem.TailBytes("textfiles/hi.txt", -10)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\n%v\n\t\t%v\n", string(data), string(tail))
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
