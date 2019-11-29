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
	"net/http"

	"github.com/bengarrett/retrotxtgo/filesystem"
	"github.com/spf13/cobra"
	"gopkg.in/gookit/color.v1"
)

// TODO: move the create.go --> have --host 8080

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: color.Primary.Sprint("Serve a text file on the in-built web server"),
	//Args: cobra.ExactArgs(1), // uncomment for Args(1) - filepath
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")

		body, err := filesystem.ReadAllBytes("textfiles/hi.txt")
		if err != nil {
			log.Fatal(err)
		}

		tmpl := template.Must(template.ParseFiles("static/html/standard.html"))
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			data := PageData{
				PreText:   string(body),
				PageTitle: "Test layout",
			}
			tmpl.Execute(w, data)
		})
		fs := http.FileServer(http.Dir("static/"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))

		fmt.Println("serving Test layout on http://localhost:8080\nCtrl-C to quit")
		http.ListenAndServe(":8080", nil)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
