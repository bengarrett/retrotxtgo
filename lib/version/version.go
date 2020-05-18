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

// Build and version information
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

// B holds build and version information.
var B = Build{
	Commit:  "",
	Date:    "",
	Domain:  "retrotxt.com",
	Version: "",
}

// JSON formats the RetroTxt version and binary compile infomation.
func JSON(indent bool) (data []byte) {
	var err error
	switch indent {
	case true:
		data, err = json.MarshalIndent(information(), "", "    ")
	default:
		data, err = json.Marshal(information())
	}
	logs.Check(logs.Err{Issue: "could not create", Arg: "json", Msg: err})
	return data
}

// Sprint formats the RetroTxt version and binary compile infomation.
func Sprint(color bool) (text string) {
	c.Enable = color
	i := information()
	text = fmt.Sprintf(logs.Cp("RetroTxt\t%s [%s]\n"), i["copyright"], i["url"]) +
		fmt.Sprintf(logs.Cinf("Version:\t%s\n"), i["app ver"]) +
		fmt.Sprintf("Go version:\t%s\n", i["go ver"]) +
		fmt.Sprintf("\nBinary:\t\t%s\n", i["exe"]) +
		fmt.Sprintf("OS/Arch:\t%s\n", i["os"]) +
		fmt.Sprintf("Build commit:\t%s\n", i["git"]) +
		fmt.Sprintf("Build date:\t%s\n", i["date"])
	return text
}

// arch humanises some common Go architecture targets.
func arch(goarch string) string {
	a := map[string]string{
		"386":   "32-bit Intel/AMD",
		"amd64": "64-bit Intel/AMD",
		"arm":   "32-bit ARM",
		"arm64": "64-bit ARM",
		"ppc64": "64-bit PowerPC",
	}
	return a[goarch]
}

// binary is the location of this program executable.
func binary() string {
	bin, err := os.Executable()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return bin
}

// semantic go version.
func semantic() string {
	ver := runtime.Version()
	if len(ver) > 2 && ver[:2] == "go" {
		return ver[2:]
	}
	return ver
}

// information and version details of retrotxt.
func information() versionInfo {
	v := versionInfo{
		"copyright": fmt.Sprintf("Copyright © 2020 Ben Garrett"),
		"url":       fmt.Sprintf("https://%s/go", B.Domain),
		"app ver":   B.Version,
		"go ver":    semantic(),
		"os":        fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		"exe":       binary(),
		"date":      localBuild(B.Date),
		"git":       B.Commit,
		"license":   fmt.Sprintf("LGPL-3.0 [https://www.gnu.org/licenses/lgpl-3.0.html]"),
	}
	if a := arch(runtime.GOARCH); a != "" {
		v["os"] += fmt.Sprintf(" [%s CPU]", a)
	}
	v["app ver"] += " (pre-alpha)"
	return v
}

// localBuild date of this binary executable.
func localBuild(date string) string {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return date
	}
	return t.Local().Format("2006 Jan 2, 15:04 MST")
}