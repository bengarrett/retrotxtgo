// Package meta handles the metadata generated through the go builder using ldflags.
package meta

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Release struct {
	// Version of RetroTxt.
	Version string
	// GitHub commit checksum.
	Commit string
	// Date in the RFC3339 format.
	Date string
	// Built by (goreleaser).
	BuiltBy string
}

// Version using semantic syntax values.
type Version struct {
	Major int
	Minor int
	Patch int
}

// App contains the version release and build metadata.
var App = Release{} //nolint:gochecknoglobals

const (
	Alpha       = "α"
	Beta        = "β"
	Placeholder = "unset"
)

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

// digits returns only the digits and decimal point values from a string.
func digits(s string) string {
	reg := regexp.MustCompile("[^0-9/.]+")
	return reg.ReplaceAllString(s, "")
}

func (v Version) String() string {
	if !v.Valid() {
		return Placeholder
	}
	p := ""
	switch {
	case v.Major == 0 && v.Minor == 0:
		p = Alpha
	case v.Major == 0:
		p = Beta
	}
	return fmt.Sprintf("%s%d.%d.%d", p, v.Major, v.Minor, v.Patch)
}

// Valid checks the version syntax.
func (v Version) Valid() bool {
	if v.Major < 0 && v.Minor < 0 && v.Patch < 0 {
		return false
	}
	return true
}
