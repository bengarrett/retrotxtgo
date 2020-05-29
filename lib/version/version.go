package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/online"
	"github.com/bengarrett/retrotxtgo/lib/str"
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

// Version details in semantic syntax
type Version struct {
	Major int
	Minor int
	Patch int
}

type versionInfo map[string]string

// B holds build and version information.
var B = Build{
	Commit:  "",
	Date:    "",
	Domain:  "retrotxt.com",
	Version: "",
}

// Digits returns only the digits and decimal point values from a string.
func Digits(s string) string {
	reg, err := regexp.Compile("[^0-9/.]+")
	if err != nil {
		return ""
	}
	return reg.ReplaceAllString(s, "")
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
	logs.ChkErr(logs.Err{Issue: "could not create", Arg: "json", Msg: err})
	return data
}

// NewRelease checks to see if the active executable matches the version hosted on GitHub.
// The ver value contains the result returned from the GitHub releases/latest API.
func NewRelease() (ok bool, ver string) {
	g, err := online.Endpoint(online.ReleaseAPI)
	if err != nil {
		logs.LogCont(err)
		return false, ver
	}
	ver = fmt.Sprint(g["tag_name"])
	if ver == "" {
		return false, ver
	}
	if c := compare(B.Version, ver); c {
		return true, ver
	}
	return false, ver
}

// Print formats and prints the RetroTxt version and binary compile infomation.
func Print(format string) (ok bool) {
	switch format {
	case "color", "c", "":
		print(Sprint(true))
	case "json", "j":
		fmt.Printf("%s\n", JSON(true))
	case "json.min", "jm":
		fmt.Printf("%s\n", JSON(false))
	case "text", "t":
		print(Sprint(false))
	default:
		return false
	}
	return true
}

// Semantic breaks down a semantic version string into major, minor and patch integers.
func Semantic(ver string) (ok bool, version Version) {
	if ver == "" {
		return false, version
	}
	vers := strings.Split(ver, ".")
	var nums [3]int
	for i, v := range vers {
		if v == "" {
			v = "0"
		}
		num, err := strconv.Atoi(Digits(v))
		if err != nil {
			return false, version
		}
		nums[i] = num
	}
	version.Major = nums[0]
	version.Minor = nums[1]
	version.Patch = nums[2]
	return true, version
}

// Sprint formats the RetroTxt version and binary compile infomation.
func Sprint(color bool) (text string) {
	c.Enable = color
	i := information()
	new, ver := NewRelease()
	var b bytes.Buffer
	fmt.Fprintf(&b, str.Cp("RetroTxt\t%s [%s]\n"), i["copyright"], i["url"])
	fmt.Fprintf(&b, str.Cinf("Version:\t%s"), i["app ver"])
	if new {
		fmt.Fprintf(&b, str.Cinf("  current: %s"), ver)
	}
	fmt.Fprintf(&b, "\nGo version:\t%s\n", i["go ver"])
	fmt.Fprintf(&b, "\nBinary:\t\t%s\n", i["exe"])
	fmt.Fprintf(&b, "OS/Arch:\t%s\n", i["os"])
	fmt.Fprintf(&b, "Build commit:\t%s\n", i["git"])
	fmt.Fprintf(&b, "Build date:\t%s\n", i["date"])
	if new {
		fmt.Fprint(&b, newRelease())
	}
	return b.String()
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

func compare(current, fetched string) bool {
	ok, c := Semantic(current)
	if !ok {
		return false
	}
	ok, f := Semantic(fetched)
	if !ok {
		return false
	}
	if c.Major < f.Major {
		return true
	}
	if c.Minor < f.Minor {
		return true
	}
	if c.Patch < f.Patch {
		return true
	}
	return false
}

// information and version details of retrotxt.
func information() versionInfo {
	v := versionInfo{
		"copyright": fmt.Sprintf("Copyright © 2020 Ben Garrett"),
		"url":       fmt.Sprintf("https://%s/go", B.Domain),
		"app ver":   B.Version,
		"go ver":    semanticGo(),
		"os":        fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		"exe":       binary(),
		"date":      localBuild(B.Date),
		"git":       B.Commit,
		"license":   fmt.Sprintf("LGPL-3.0 [https://www.gnu.org/licenses/lgpl-3.0.html]"),
	}
	if a := arch(runtime.GOARCH); a != "" {
		v["os"] += fmt.Sprintf(" [%s CPU]", a)
	}
	v["app ver"] = fmt.Sprintf("α%s", v["app ver"])
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

func newRelease() *bytes.Buffer {
	s := "A newer edition of RetroTxt is available!\n" +
		"Learn more at https://github.com/bengarrett/retrotxtgo"
	return str.Border(s)
}

// semantic go version.
func semanticGo() string {
	ver := runtime.Version()
	if len(ver) > 2 && ver[:2] == "go" {
		return ver[2:]
	}
	return ver
}
