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
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version information for RetroTxt",
	Run: func(cmd *cobra.Command, args []string) {
		printVer()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVer() {
	fmt.Printf("RetroTxt\thttps://%s\n", Www)
	fmt.Printf("Version:\t%s (pre-alpha)\n", Ver)
	fmt.Printf("Go version:\t%s\n", runtime.Version())
	fmt.Printf("OS/Arch:\t%s/%s\n", runtime.GOOS, runtime.GOARCH)
	bin, err := os.Executable()
	if err == nil {
		fmt.Printf("Binary:\t\t%s\n", bin)
	}
}
