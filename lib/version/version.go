// Package version is the release and build information.
package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/online"
	"github.com/bengarrett/retrotxtgo/lib/str"
	gookit "github.com/gookit/color"
	gap "github.com/muesli/go-app-paths"
	"gopkg.in/yaml.v3"
)

// Data for the release information.
type Data struct {
	// Built by (usually goreleaser).
	BuiltBy string
	// Git commit SHA hash.
	GitHash string
	// Date in RFC3339.
	Date string
	// Domain name for the website.
	Domain string
	// Version of RetroTxt.
	Version string
}

// Release placeholder data.
var Release = Data{ //nolint:gochecknoglobals
	BuiltBy: "go build",
	GitHash: "unset",
	Date:    time.Now().Format("2006 Jan 2, 15:04 MST"),
	Domain:  "retrotxt.com",
	Version: "",
}

// Output the version data.
type Output struct {
	App       string `json:"version"`
	By        string `json:"buildBy"`
	Copyright string `json:"copyright"`
	Date      string `json:"date"`
	Exe       string `json:"binary"`
	Git       string `json:"gitCommit"`
	GoVer     string `json:"goVersion"`
	License   string `json:"license"`
	OS        string `json:"os"`
	URL       string `json:"url"`
}

// Print the results of the version command as ANSI or plain text.
func (o *Output) String(color bool) string {
	gookit.Enable = color
	update, ver := NewRelease()
	var b bytes.Buffer
	fmt.Fprintf(&b, str.Cp("RetroTxt\t%s [%s]\n"), o.Copyright, o.URL)
	fmt.Fprintf(&b, str.Cinf("Version:\t%s"), o.App)
	if update {
		fmt.Fprintf(&b, str.Cinf("  current: %s"), ver)
	}
	fmt.Fprintf(&b, "\nGo version:\t%s\n", o.GoVer)
	fmt.Fprintf(&b, "\nBinary:\t\t%s\n", o.Exe)
	fmt.Fprintf(&b, "OS/Arch:\t%s\n", o.OS)
	fmt.Fprintf(&b, "\nBuild commit:\t%s\n", o.Git)
	fmt.Fprintf(&b, "Date:\t\t%s\n", o.Date)
	fmt.Fprintf(&b, "By:\t\t%s\n", o.By)
	if update {
		fmt.Fprint(&b, newRelease())
	}
	return b.String()
}

// Print the results of the version command as JSON.
func (o *Output) json() (data []byte) {
	data, err := json.MarshalIndent(&o, "", "    ")
	if err != nil {
		logs.ProblemMarkFatal("json", ErrMarshal, err)
	}
	return data
}

// Print the results of the version command as minified JSON.
func (o *Output) jsonMin() (data []byte) {
	data, err := json.Marshal(&o)
	if err != nil {
		logs.ProblemMarkFatal("json", ErrMarshal, err)
	}
	return data
}

// Version using semantic syntax values.
type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	const alpha, beta = "α", "β"
	if !v.valid() {
		return "unset"
	}
	p := ""
	switch {
	case v.Major == 0 && v.Minor == 0:
		p = alpha
	case v.Major == 0:
		p = beta
	}
	return fmt.Sprintf("%s%d.%d.%d", p, v.Major, v.Minor, v.Patch)
}

// Valid checks the Version syntax is correct.
func (v Version) valid() bool {
	if v.Major < 0 && v.Minor < 0 && v.Patch < 0 {
		return false
	}
	return true
}

// Cache the version data.
type Cache struct {
	Etag string `yaml:"etag"`
	Ver  string `yaml:"version"`
}

const cacheFile = "api.github.cache"

// NewRelease checks to see if the active executable matches the version hosted on GitHub.
// The ver value contains the result returned from the GitHub releases/latest API.
func NewRelease() (ok bool, ver string) {
	etag, ver := cacheGet()
	cache, data, err := online.Endpoint(online.ReleaseAPI, etag)
	if err != nil {
		logs.Save(err)
		return false, ver
	}
	if !cache {
		ver = fmt.Sprint(data["tag_name"])
		if ver == "" {
			return false, ver
		}
		if fmt.Sprintf("%T", data["etag"]) == "string" {
			if data["etag"].(string) != "" {
				if err = cacheSet(data["etag"].(string), ver); err != nil {
					logs.Save(err)
				}
			}
		}
	}
	if comp := compare(Release.Version, ver); comp {
		return true, ver
	}
	return false, ver
}

// Print and format the RetroTxt version plus the binary compile information.
func Print(format string) (ok bool) {
	m := marshal()
	switch format {
	case "color", "c", "":
		fmt.Print(m.String(true))
	case "json", "j":
		fmt.Printf("%s\n", m.json())
	case "json.min", "jm":
		fmt.Printf("%s\n", m.jsonMin())
	case "text", "t":
		fmt.Print(m.String(false))
	default:
		return false
	}
	return true
}

// Semantic breaks down a semantic version string into major, minor and patch integers.
func Semantic(ver string) Version {
	invalid := Version{-1, -1, -1}
	if ver == "" {
		return invalid
	}
	vers, nums := strings.Split(ver, "."), [3]int{}
	for i, v := range vers {
		if v == "" {
			v = "0"
		}
		num, err := strconv.Atoi(digits(v))
		if err != nil {
			return invalid
		}
		nums[i] = num
	}
	return Version{
		Major: nums[0],
		Minor: nums[1],
		Patch: nums[2],
	}
}

// Arch humanizes some common Go target architectures.
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

// Binary is the location of this executable program.
func binary() string {
	bin, err := os.Executable()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return bin
}

// CacheGet returns the stored Github API ETag HTTP header and release version.
func cacheGet() (etag, ver string) {
	cf, err := home().DataPath(cacheFile)
	if err != nil {
		logs.Save(err)
		return
	}
	if _, err = os.Stat(cf); os.IsNotExist(err) {
		return
	}
	f, err := ioutil.ReadFile(cf)
	if err != nil {
		logs.Save(err)
	}
	var cache Cache
	if err = yaml.Unmarshal(f, &cache); err != nil {
		logs.Save(err)
	}
	// if either value is missing, delete the broken cache
	if cache.Etag == "" || cache.Ver == "" {
		err = os.Remove(cf)
		logs.Save(err)
		return "", ""
	}
	return cache.Etag, cache.Ver
}

// CacheSet saves the Github API ETag HTTP header and release version.
func cacheSet(etag, ver string) error {
	if etag == "" || ver == "" {
		return nil
	}
	cache := Cache{
		Etag: etag,
		Ver:  ver,
	}
	out, err := yaml.Marshal(&cache)
	if err != nil {
		return fmt.Errorf("%s: %w", err, ErrCacheYaml)
	}
	f, err := home().DataPath(cacheFile)
	if err != nil {
		return fmt.Errorf("%q: %w", cacheFile, ErrCacheData)
	}
	if _, _, err := filesystem.Save(f, out...); err != nil {
		return fmt.Errorf("%s: %w", err, ErrCacheSave)
	}
	return nil
}

// Compare the version of this executable program against
// the latest version hosted on GitHub.
func compare(current, fetched string) bool {
	cur := Semantic(current)
	if !cur.valid() {
		return false
	}
	f := Semantic(fetched)
	if !f.valid() {
		return false
	}
	if cur.Major < f.Major {
		return true
	}
	if cur.Minor < f.Minor {
		return true
	}
	if cur.Patch < f.Patch {
		return true
	}
	return false
}

// CopyrightYear uses the build date as a the year of copyright.
func copyrightYear(date string) string {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return time.Now().Format("2006")
	}
	return t.Local().Format("2006")
}

// Digits returns only the digits and decimal point values from a string.
func digits(s string) string {
	reg := regexp.MustCompile("[^0-9/.]+")
	return reg.ReplaceAllString(s, "")
}

// Home returns the user's home directory determined by the host operating system.
func home() *gap.Scope {
	return gap.NewScope(gap.User, "retrotxt")
}

// LocalBuild date of this binary executable.
func localBuild(date string) string {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return date
	}
	return t.Local().Format("2006 Jan 2, 15:04 MST")
}

func marshal() Output {
	v := Output{
		Copyright: fmt.Sprintf("Copyright © %s Ben Garrett", copyrightYear(Release.Date)),
		URL:       fmt.Sprintf("https://%s/go", Release.Domain),
		App:       Semantic(Release.Version).String(),
		GoVer:     semanticGo(),
		OS:        fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Exe:       binary(),
		Date:      localBuild(Release.Date),
		By:        Release.BuiltBy,
		Git:       Release.GitHash,
		License:   "LGPL-3.0 [https://www.gnu.org/licenses/lgpl-3.0.html]",
	}
	if a := arch(runtime.GOARCH); a != "" {
		// as of Go v1.16, darwin/arm64 equals the Mx architecture
		if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
			v.OS += " [Apple M1 CPU]"
			return v
		}
		v.OS += fmt.Sprintf(" [%s CPU]", a)
	}
	return v
}

// NewRelease notification.
func newRelease() *bytes.Buffer {
	s := "A newer edition of RetroTxt is available!\n" +
		"Learn more at https://retrotxt.com/go"
	return str.Border(s)
}

// Semantic go version.
func semanticGo() string {
	ver := runtime.Version()
	if len(ver) > 2 && ver[:2] == "go" {
		return ver[2:]
	}
	return ver
}
