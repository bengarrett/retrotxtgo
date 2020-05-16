package version

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	c "github.com/gookit/color"
)

// Build ..
type Build struct {
	// Commit git SHA
	Commit string
	// Date in RFC3339
	Date string
	// Domain name for the website
	Domain string
	// Version of RetroTxt
	Version string
}

type versionInfo map[string]string

// Inf ...
var Inf = Build{
	Commit:  "",
	Date:    "",
	Domain:  "retrotxt.com",
	Version: "",
}

// JSON ..
func JSON(indent bool) (data []byte) {
	var err error
	switch indent {
	case true:
		data, err = json.MarshalIndent(info(), "", "    ")
	default:
		data, err = json.Marshal(info())
	}
	logs.Check(logs.Err{"could not create", "json", err})
	return data
}

// Sprint formats the RetroTxt version and binary compile infomation.
func Sprint(color bool) (text string) {
	c.Enable = color
	i := info()
	text = fmt.Sprintf(logs.Cp("RetroTxt\t%s [%s]\n"), i["copyright"], i["url"]) +
		fmt.Sprintf(logs.Cinf("Version:\t%s\n"), i["app ver"]) +
		fmt.Sprintf("Go version:\t%s\n", i["go ver"]) +
		fmt.Sprintf("\nBinary:\t\t%s\n", i["exe"]) +
		fmt.Sprintf("OS/Arch:\t%s\n", i["os"]) +
		fmt.Sprintf("Build commit:\t%s\n", i["git"]) +
		fmt.Sprintf("Build date:\t%s\n", i["date"])
	return text
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

func goVer() string {
	ver := runtime.Version()
	if len(ver) > 2 && ver[:2] == "go" {
		return ver[2:]
	}
	return ver
}

func info() versionInfo {
	v := versionInfo{
		"copyright": fmt.Sprintf("Copyright © 2020 Ben Garrett"),
		"url":       fmt.Sprintf("https://%s/go", Inf.Domain),
		"app ver":   Inf.Version,
		"go ver":    goVer(),
		"os":        fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		"exe":       binary(),
		"date":      locBuildDate(Inf.Date),
		"git":       Inf.Commit,
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
