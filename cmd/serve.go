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
	"net/http"

	"github.com/bengarrett/retrotxtgo/filesystem"
	"github.com/spf13/cobra"
)

//PageData contains template data used by layout.html
type PageData struct {
	PageTitle string
	BodyText  string
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	//Args: cobra.ExactArgs(1), // uncomment for Args(1) - filepath
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
		tmpl := template.Must(template.ParseFiles("layout.html"))
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			data := PageData{
				BodyText:  filesystem.Read("textfiles/hi.txt"),
				PageTitle: "Test layout",
			}
			tmpl.Execute(w, data)
		})
		fs := http.FileServer(http.Dir("static/"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		http.ListenAndServe(":80", nil)

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
