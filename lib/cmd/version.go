package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

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
		CheckFlag(versionPrint(viper.GetString("version.format")))
	},
}

func init() {
	InitDefaults()
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVarP(&versionFmt, "format", "f",
		viper.GetString("version.format"), "output format \noptions: "+versionFormats)
	_ = viper.BindPFlag("version.format", versionCmd.Flags().Lookup("format"))
}

func versionPrint(format string) (err ErrorFmt) {
	switch format {
	case "color", "c", "":
		print(versionText(true))
	case "json", "j":
		fmt.Printf("%s\n", versionJSON(true))
	case "json.min", "jm":
		fmt.Printf("%s\n", versionJSON(false))
	case "text", "t":
		print(versionText(false))
	default:
		return ErrorFmt{"format", format, fmt.Errorf(versionFormats)}
	}
	return err
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
		return fmt.Sprintf("error: %v", err)
	}
	return bin
}

func info() versionInfo {
	v := versionInfo{
		"copyright": fmt.Sprintf("Copyright © 2020 Ben Garrett"),
		"url":       fmt.Sprintf("https://%s/go", Www),
		"app ver":   BuildVer,
		"go ver":    goVer(),
		"os":        fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		"exe":       binary(),
		"date":      locBuildDate(BuildDate),
		"git":       BuildCommit,
		"license":   fmt.Sprintf("LGPL-3.0 [https://www.gnu.org/licenses/lgpl-3.0.html]"),
	}
	if a := arch(runtime.GOARCH); a != "" {
		v["os"] += fmt.Sprintf(" [%s CPU]", a)
	}
	v["app ver"] += " (pre-alpha)"
	return v
}

func locBuildDate(date string) string {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return date
	}
	return t.Local().Format("2006 Jan 2, 15:04 MST")
}

func goVer() string {
	ver := runtime.Version()
	if len(ver) > 2 && ver[:2] == "go" {
		return ver[2:]
	}
	return ver
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

func versionText(colr bool) (text string) {
	color.Enable = colr
	i := info()
	text = fmt.Sprintf(cp("RetroTxt\t%s [%s]\n"), i["copyright"], i["url"]) +
		fmt.Sprintf(cinf("Version:\t%s\n"), i["app ver"]) +
		fmt.Sprintf("Go version:\t%s\n", i["go ver"]) +
		fmt.Sprintf("\nBinary:\t\t%s\n", i["exe"]) +
		fmt.Sprintf("OS/Arch:\t%s\n", i["os"]) +
		fmt.Sprintf("Build commit:\t%s\n", i["git"]) +
		fmt.Sprintf("Build date:\t%s\n", i["date"])
	return text
}
