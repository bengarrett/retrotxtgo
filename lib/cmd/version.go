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
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/gookit/color.v1"
)

type versionInfo map[string]string

const versionFormats string = "color, json, json.min, text"

var versionFmt string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use: "version",
	//Aliases: []string{"ver", "v"},
	Short: "Version information for RetroTxt",
	Long: `Version information for Retrotxt

The shown ` + cc("RetroTxt URL") + ` is the weblink to the application Github page.

` + cc("Version") + ` number reflects ` + ci("[major].[minor].[patch]") + `.
* Major is a generational iteration that may break backwards compatibility.
* Minor changes are increased whenever new features are added.
* Patch reflect hot fixes or bug corrections.

` + cc("Go version") + ` reports the edition of Go used to build this application.
` + cc("OS/Arch") + ` reports both the operating system and CPU architecture.

` + cc("Binary") + ` should return the path name of this program. It maybe inaccurate
if it is launched through an operating system symlink.`,
	Run: func(cmd *cobra.Command, args []string) {
		f := viper.GetString("version.format")
		switch f {
		case "color", "c", "":
			versionText(true)
		case "json", "j":
			fmt.Printf("%s\n", versionJSON(true))
		case "json.min", "jm":
			fmt.Printf("%s\n", versionJSON(false))
		case "text", "t":
			versionText(false)
		default:
			CheckFlag(ErrorFmt{"format", f, fmt.Errorf(versionFormats)})
		}
	},
}

func init() {
	InitDefaults()
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVarP(&versionFmt, "format", "f", viper.GetString("version.format"), "output format \noptions: "+versionFormats)
	_ = viper.BindPFlag("version.format", versionCmd.Flags().Lookup("format"))
}

func versionJSON(indent bool) []byte {
	var j []byte
	var err error
	switch indent {
	case true:
		j, err = json.MarshalIndent(info(), "", "    ")
	default:
		j, err = json.Marshal(info())
	}
	Check(ErrorFmt{"could not create", "json", err})
	return j
}

func versionText(c bool) {
	color.Enable = c
	i := info()
	fmt.Printf(cp("RetroTxt\t%s [%s]\n"), i["copyright"], i["url"])
	fmt.Printf(cinf("Version:\t%s\n"), i["app ver"])
	fmt.Printf("Go version:\t%s\n", i["go ver"])
	fmt.Printf("\nBinary:\t\t%s\n", i["exe"])
	fmt.Printf("OS/Arch:\t%s\n", i["os"])
	fmt.Printf("Build commit:\t%s\n", i["git"])
	fmt.Printf("Build date:\t%s\n", i["date"])
}

func arch(v string) string {
	a := map[string]string{
		"386":   "32-bit Intel/AMD",
		"amd64": "64-bit Intel/AMD",
		"arm":   "32-bit ARM",
		"arm64": "64-bit ARM",
		"ppc64": "64-bit PowerPC",
	}
	return a[v]
}

func binary() string {
	bin, err := os.Executable()
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return bin
}

func info() versionInfo {
	v := versionInfo{
		"copyright": fmt.Sprintf("Copyright © 2019 Ben Garrett"),
		"url":       fmt.Sprintf("https://%s/go", Www),
		"app ver":   GoBuildVer,
		"go ver":    runtime.Version(),
		"os":        fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		"exe":       binary(),
		"date":      GoBuildDate,
		"git":       GoBuildGitCommit,
		"license":   fmt.Sprintf("LGPL-3.0 [https://www.gnu.org/licenses/lgpl-3.0.html]"),
	}
	if a := arch(runtime.GOARCH); a != "" {
		v["os"] += fmt.Sprintf(" [%s CPU]", a)
	}
	if GoBuildGitCount != "" {
		v["app ver"] += fmt.Sprintf(".%s", GoBuildGitCount)
	} else {
		v["app ver"] += ".?"
	}
	v["app ver"] += " (pre-alpha)"
	if GoBuildDate == "" {
		v["date"] = "n/a"
	} else {
		v["date"] += " UTC"
	}
	if GoBuildGitCommit == "" {
		v["git"] = "n/a"
	}
	return v
}
