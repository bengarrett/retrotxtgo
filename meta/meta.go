// Package meta handles the metadata generated through the go builder using ldflags.
package meta

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Release metadata.
type Release struct {
	// Version of this program.
	Version string
	// GitHub commit checksum.
	Commit string
	// Date in the RFC3339 format.
	Date string
	// Built by name (goreleaser).
	BuiltBy string
}

// Version using semantic syntax values.
type Version struct {
	Major int // Major for incompatible API changes.
	Minor int // Minor for functionality in a backwards compatible manner.
	Patch int // Patch for backwards compatible bug fixes.
}

// App contains the version release and build metadata.
var App = Release{}

const (
	// Alpha Greek character.
	Alpha = "α"
	// Beta Greek character.
	Beta = "β"
	// GoBuild version when no ldflags are in use.
	GoBuild = "0.0.0"
	// Placeholder string when no ldflags are in use.
	Placeholder = "unset"
	// Bin is the binary filename of this program.
	Bin = "retrotxt"
	// Dir is the sub-directory name used for configuration and temporary paths.
	Dir = "retrotxt"
	// Name of this program.
	Name = "Retrotxt"
	// URL for this program's website.
	URL = "https://retrotxt.com/go"
	// Copyright © year range.
	Copyright = "2020-24"
)

// String returns the release version string.
func String() string {
	return Semantic(App.Version).String()
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
		num, err := strconv.Atoi(Digits(v))
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

// Digits returns only the digits and decimal point values from a string.
func Digits(s string) string {
	r := regexp.MustCompile("[^0-9/.]+")
	return r.ReplaceAllString(s, "")
}

// String returns the semantic version string.
// If the version is invalid, it returns a placeholder, unset string.
func (v Version) String() string {
	if !v.Valid() {
		return Placeholder
	}
	p := ""
	switch {
	case v.Major == 0 && v.Minor == 0 && v.Patch == 0:
		p = "x"
	case v.Major == 0 && v.Minor == 0:
		p = Alpha
	case v.Major == 0:
		p = Beta
	}
	return fmt.Sprintf("%s%d.%d.%d", p, v.Major, v.Minor, v.Patch)
}

// Valid reports whether the sematic versioning values are valid.
// Values that are less than zero are considered invalid.
func (v Version) Valid() bool {
	if v.Major < 0 && v.Minor < 0 && v.Patch < 0 {
		return false
	}
	return true
}
