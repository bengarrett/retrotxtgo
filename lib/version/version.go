package version

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	gookit "github.com/gookit/color"
	gap "github.com/muesli/go-app-paths"
	"gopkg.in/yaml.v3"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/online"
	"retrotxt.com/retrotxt/lib/str"
)

var (
	// ErrCacheYaml set cache yaml error.
	ErrCacheYaml = errors.New("set cache cannot marshal yaml")
	// ErrCacheData set cache data path.
	ErrCacheData = errors.New("set cache cannot create a data path")
	// ErrCacheSave set cache save.
	ErrCacheSave = errors.New("set cache cannot save data")
)

// Build and version information.
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

// Output the version data.
type Output struct {
	App       string `json:"version"`
	Copyright string `json:"copyright"`
	Date      string `json:"date"`
	Exe       string `json:"binary"`
	Git       string `json:"gitCommit"`
	GoVer     string `json:"goVersion"`
	License   string `json:"license"`
	OS        string `json:"os"`
	URL       string `json:"url"`
}

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
	fmt.Fprintf(&b, "Build commit:\t%s\n", o.Git)
	fmt.Fprintf(&b, "Build date:\t%s\n", o.Date)
	if update {
		fmt.Fprint(&b, newRelease())
	}
	return b.String()
}

func (o *Output) json() (data []byte) {
	data, err := json.MarshalIndent(&o, "", "    ")
	if err != nil {
		logs.Fatal("version could not marshal", "json", err)
	}
	return data
}

func (o *Output) jsonMin() (data []byte) {
	data, err := json.Marshal(&o)
	if err != nil {
		logs.Fatal("version could not marshal", "json", err)
	}
	return data
}

// Version details in semantic syntax.
type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	const alpha, beta = "α", "β"
	if !v.valid() {
		return ""
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

// Cache of version data.
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
		logs.Log(err)
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
					logs.Log(err)
				}
			}
		}
	}
	if comp := compare(B.Version, ver); comp {
		return true, ver
	}
	return false, ver
}

// Print formats and prints the RetroTxt version and binary compile information.
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

// cacheGet returns the stored Github API ETag HTTP header and release version.
func cacheGet() (etag, ver string) {
	cf, err := home().DataPath(cacheFile)
	if err != nil {
		logs.Log(err)
		return
	}
	if _, err = os.Stat(cf); os.IsNotExist(err) {
		return
	}
	f, err := ioutil.ReadFile(cf)
	if err != nil {
		logs.Log(err)
	}
	var cache Cache
	if err = yaml.Unmarshal(f, &cache); err != nil {
		logs.Log(err)
	}
	// if either value is missing, delete the broken cache
	if cache.Etag == "" || cache.Ver == "" {
		err = os.Remove(cf)
		logs.Log(err)
		return "", ""
	}
	return cache.Etag, cache.Ver
}

// cacheSet saves the Github API ETag HTTP header and release version.
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
	if _, err := filesystem.Save(f, out...); err != nil {
		return fmt.Errorf("%s: %w", err, ErrCacheSave)
	}
	return nil
}

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

// copyrightYear uses the build date as a the year of copyright.
func copyrightYear(date string) string {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return date
	}
	return t.Local().Format("2006")
}

// digits returns only the digits and decimal point values from a string.
func digits(s string) string {
	reg := regexp.MustCompile("[^0-9/.]+")
	return reg.ReplaceAllString(s, "")
}

// home directory path determined by the host operating system.
func home() *gap.Scope {
	return gap.NewScope(gap.User, "retrotxt")
}

// localBuild date of this binary executable.
func localBuild(date string) string {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return date
	}
	return t.Local().Format("2006 Jan 2, 15:04 MST")
}

func marshal() Output {
	v := Output{
		Copyright: fmt.Sprintf("Copyright © %s Ben Garrett", copyrightYear(B.Date)),
		URL:       fmt.Sprintf("https://%s/go", B.Domain),
		App:       Semantic(B.Version).String(),
		GoVer:     semanticGo(),
		OS:        fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Exe:       binary(),
		Date:      localBuild(B.Date),
		Git:       B.Commit,
		License:   "LGPL-3.0 [https://www.gnu.org/licenses/lgpl-3.0.html]",
	}
	if a := arch(runtime.GOARCH); a != "" {
		if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
			v.OS += " [Apple M1 CPU]"
		} else {
			v.OS += fmt.Sprintf(" [%s CPU]", a)
		}
	}
	return v
}

func newRelease() *bytes.Buffer {
	s := "A newer edition of RetroTxt is available!\n" +
		"Learn more at https://retrotxt.com/go"
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
