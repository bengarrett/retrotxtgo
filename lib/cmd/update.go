package cmd

import (
	"bytes"
	"errors"
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

var (
	ErrCacheYaml = errors.New("set cache cannot marshal yaml")
	ErrCacheData = errors.New("set cache cannot create a data path")
	ErrCacheSave = errors.New("set cache cannot save data")
)

// Cache the version data.
type Cache struct {
	Etag string `yaml:"etag"`
	Ver  string `yaml:"version"`
}

const cacheFile = "api.github.cache"

// chkRelease checks to see if the active executable matches the version hosted on GitHub.
// The ver value contains the result returned from the GitHub releases/latest API.
func chkRelease() (ok bool, ver string) {
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
	if comp := compare(meta.App.Version, ver); comp {
		return true, ver
	}
	return false, ver
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

// Home returns the user's home directory determined by the host operating system.
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
