package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/online"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	gap "github.com/muesli/go-app-paths"
	yaml "gopkg.in/yaml.v2"
)

// Cache the version data.
type Cache struct {
	Etag string `yaml:"etag"`
	Ver  string `yaml:"version"`
}

const cacheFile = "api.github.cache"

// chkRelease returns true if the latest GitHub repo tag is newer than the program version.
// The returned v string is the latest GitHub repo tag.
func chkRelease() (newRel bool, v string) {
	if meta.App.Version == meta.GoBuild {
		return false, ""
	}
	etag, v := cacheGet()
	cache, data, err := online.Endpoint(online.ReleaseAPI, etag)
	if err != nil {
		logs.Save(err)
		return false, v
	}
	if !cache {
		v = fmt.Sprint(data["tag_name"])
		if v == "" {
			return false, v
		}
		if fmt.Sprintf("%T", data["etag"]) == "string" {
			if data["etag"].(string) != "" {
				if err = cacheSet(data["etag"].(string), v); err != nil {
					logs.Save(err)
				}
			}
		}
	}
	if comp := compare(meta.App.Version, v); comp {
		return true, v
	}
	return false, v
}

// cacheGet reads the stored GitHub API, HTTP ETag header and release version.
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

// cacheSet saves the Github API, ETag HTTP header and release version.
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
	if _, _, err := filesystem.Write(f, out...); err != nil {
		return fmt.Errorf("%s: %w", err, ErrCacheSave)
	}
	return nil
}

// compare the version value of this executable program against the latest version hosted on GitHub.
func compare(current, fetched string) bool {
	cur := meta.Semantic(current)
	if !cur.Valid() {
		return false
	}
	f := meta.Semantic(fetched)
	if !f.Valid() {
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

// home returns the user home directory.
func home() *gap.Scope {
	return gap.NewScope(gap.User, meta.Dir)
}

// newRelease notification box and text.
func newRelease(old, newest string) *bytes.Buffer {
	s := fmt.Sprintf("%s%s%s\n%s%s\n%s → %s",
		"A newer edition of ", meta.Name, " is available!",
		"Learn more at ", meta.URL, meta.Semantic(old), newest)
	return str.Border(s)
}
