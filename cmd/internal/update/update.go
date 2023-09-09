package update

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/pkg/filesystem"
	"github.com/bengarrett/retrotxtgo/pkg/logs"
	"github.com/bengarrett/retrotxtgo/pkg/online"
	"github.com/bengarrett/retrotxtgo/pkg/str"
	gap "github.com/muesli/go-app-paths"
	"gopkg.in/yaml.v3"
)

var (
	ErrCacheYaml = errors.New("set cache cannot marshal yaml")
	ErrCacheData = errors.New("set cache cannot create a data path")
	ErrCacheSave = errors.New("set cache cannot save data")
)

// Cache the release data.
type Cache struct {
	Etag    string `yaml:"etag"`
	Version string `yaml:"version"`
}

const cacheFile = "api.github.cache"

// String the new release notification box and text.
func String(old, current string) *bytes.Buffer {
	s := fmt.Sprintf("%s%s%s\n%s%s\n%s → %s",
		"A newer edition of ", meta.Name, " is available!",
		"Learn more at ", meta.URL, meta.Semantic(old), current)
	return str.Border(s)
}

// Check the GitHub for the newest release tag.
// The returned string will only contain the newest available release tag
// if the local program version is out of date.
func Check() (string, error) {
	if meta.App.Version == meta.GoBuild {
		return "", nil
	}
	etag, tag := CacheGet()
	cache, data, err := online.Endpoint(online.ReleaseAPI, etag)
	if err != nil {
		return "", err
	}
	if !cache { //nolint:nestif
		tag = fmt.Sprint(data["tag_name"])
		if tag == "" {
			return "", nil
		}
		if fmt.Sprintf("%T", data["etag"]) == "string" {
			if data["etag"].(string) != "" {
				if err = CacheSet(data["etag"].(string), tag); err != nil {
					return "", err
				}
			}
		}
	}
	if comp := Compare(meta.App.Version, tag); comp {
		return tag, nil
	}
	return "", nil
}

// CacheGet reads the stored GitHub API, HTTP ETag header and release version.
func CacheGet() (etag, version string) {
	cf, err := home().DataPath(cacheFile)
	if err != nil {
		logs.Save(err)
		return
	}
	if _, err = os.Stat(cf); os.IsNotExist(err) {
		return
	}
	f, err := os.ReadFile(cf)
	if err != nil {
		logs.Save(err)
	}
	var cache Cache
	if err = yaml.Unmarshal(f, &cache); err != nil {
		logs.Save(err)
	}
	// if either value is missing, delete the broken cache
	if cache.Etag == "" || cache.Version == "" {
		err = os.Remove(cf)
		logs.Save(err)
		return "", ""
	}
	return cache.Etag, cache.Version
}

// CacheSet saves the Github API, ETag HTTP header and release version.
func CacheSet(etag, version string) error {
	if etag == "" || version == "" {
		return nil
	}
	cache := Cache{
		Etag:    etag,
		Version: version,
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

// Compare the release version of this program against the latest version hosted on GitHub.
func Compare(current, fetched string) bool {
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

// home is the user's home directory.
func home() *gap.Scope {
	return gap.NewScope(gap.User, meta.Dir)
}
