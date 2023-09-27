// Package update provides the ability to check GitHub for the newest release tag.
package update

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/logs"
	"github.com/bengarrett/retrotxtgo/pkg/online"
	"github.com/bengarrett/retrotxtgo/term"
	gap "github.com/muesli/go-app-paths"
	"gopkg.in/yaml.v3"
)

var (
	ErrCacheYaml = errors.New("set cache cannot marshal yaml")
	ErrCacheData = errors.New("set cache cannot create a data path")
	ErrCacheSave = errors.New("set cache cannot save data")
)

// Cache of the newest production release.
// The Etag string is used to check if the release has changed,
// by requesting a tiny HTTP header from the GitHub API instead
// of the full response.
type Cache struct {
	Etag    string `yaml:"etag"`    // http etag header
	Version string `yaml:"version"` // semantic version
}

const cacheFile = "api.github.cache"

// Notice writes a new release notification box and text.
func Notice(w io.Writer, old, current string) {
	if w == nil {
		w = io.Discard
	}
	s := fmt.Sprintf("%s%s%s\n%s%s\n%s → %s",
		"A newer edition of ", meta.Name, " is available!",
		"Learn more at ", meta.URL, meta.Semantic(old), current)
	term.Border(w, s)
}

// Check the GitHub for the newest release tag.
// The returned string will only contain the newest available release tag
// if the local program version is out of date.
func Check() (string, error) {
	if meta.App.Version == meta.GoBuild {
		return "", nil
	}
	cache := CacheGet()
	etag, tag := cache.Etag, cache.Version
	c, data, err := online.Endpoint(online.ReleaseAPI, etag)
	if err != nil {
		return "", err
	}
	if !c {
		tag = fmt.Sprint(data["tag_name"])
		if tag == "" {
			return "", nil
		}
		if s, ok := data["etag"].(string); ok && s != "" {
			if err := CacheSet(s, tag); err != nil {
				return "", err
			}
		}
	}
	if comp := Compare(meta.App.Version, tag); comp {
		return tag, nil
	}
	return "", nil
}

// CacheGet reads and returns the locally cached GitHub API.
func CacheGet() Cache {
	cf, err := home().DataPath(cacheFile)
	if err != nil {
		logs.Sprint(err)
		return Cache{}
	}
	if _, err = os.Stat(cf); os.IsNotExist(err) {
		return Cache{}
	}
	f, err := os.ReadFile(cf)
	if err != nil {
		logs.Sprint(err)
	}
	var cache Cache
	if err := yaml.Unmarshal(f, &cache); err != nil {
		logs.Sprint(err)
	}
	// if either value is missing, delete the broken cache
	if cache.Etag == "" || cache.Version == "" {
		err = os.Remove(cf)
		logs.Sprint(err)
		return Cache{}
	}
	return cache
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
		return fmt.Errorf("%w: %w", err, ErrCacheYaml)
	}
	f, err := home().DataPath(cacheFile)
	if err != nil {
		return fmt.Errorf("%q: %w", cacheFile, ErrCacheData)
	}
	if _, _, err := fsys.Write(f, out...); err != nil {
		return fmt.Errorf("%w: %w", err, ErrCacheSave)
	}
	return nil
}

// Compare reports whether the release version of this program
// matches the latest version hosted on GitHub.
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
