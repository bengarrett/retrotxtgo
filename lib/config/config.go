// Package config handles the user configations.
package config

import (
	"fmt"
	"os"
	"sort"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/viper"
)

const (
	filemode  os.FileMode = 0660
	namedFile             = "config.yaml"
)

// Hints for configuration values.
type Hints map[string]string

// Tip provides a brief help on the config file configurations.
func Tip() Hints {
	return Hints{
		get.Editor:     "text editor to launch when using " + str.Example("config edit"),
		get.FontEmbed:  "encode and embed the font as Base64 binary-to-text within the CSS",
		get.FontFamily: "specifies the font to use with the HTML",
		get.LayoutTmpl: "HTML template for the layout of CSS, JS and fonts",
		get.Author:     "defines the name of the page authors",
		get.Scheme:     "specifies one or more color schemes with which the page is compatible",
		get.Desc:       "a short and accurate summary of the content of the page",
		get.Genr:       fmt.Sprintf("include the %s version and page generation date?", meta.Name),
		get.Keywords:   "words relevant to the page content",
		get.Notlate:    "used to declare that the page should not be translated by Google Translate",
		get.Referr:     "controls the Referer HTTP header attached to requests sent from the page",
		get.Rtx:        "include a custom tag containing the meta information of the source textfile",
		get.Bot:        "behavor that crawlers from Google, Bing and other engines should use with the page",
		get.Theme:      "indicates a suggested color that user agents should use to customize the display of the page",
		get.Title:      "page title that is shown in the browser tab",
		get.SaveDir:    fmt.Sprintf("directory to store HTML assets created by %s", meta.Name),
		get.Serve:      "serve files using an internal web server with this port",
		get.Stylei:     "syntax highlighter for the config info output",
		get.Styleh:     "syntax highlighter for html previews",
	}
}

// Settings types and names to be saved as YAML.
type Settings struct {
	Editor string
	HTML   struct {
		Font struct {
			Embed  bool   `yaml:"embed"`
			Family string `yaml:"family"`
			Size   string `yaml:"size"`
		}
		Layout string `yaml:"layout"`
		Meta   struct {
			Author      string `yaml:"author"`
			ColorScheme string `yaml:"color-scheme"`
			Description string `yaml:"description"`
			Keywords    string `yaml:"keywords"`
			Referrer    string `yaml:"referrer"`
			Robots      string `yaml:"robots"`
			ThemeColor  string `yaml:"theme-color"`
			Generator   bool   `yaml:"generator"`
			Notranslate bool   `yaml:"notranslate"`
			RetroTxt    bool   `yaml:"retrotxt"`
		}
		Title string `yaml:"title"`
	}
	SaveDirectory string `yaml:"save-directory"`
	ServerPort    uint   `yaml:"serve"`
	Style         struct {
		Info string `yaml:"info"`
		HTML string `yaml:"html"`
	}
}

func cmdPath() string {
	return fmt.Sprintf("%s config", meta.Bin)
}

// Formats choices for flags.
type Formats struct {
	Info [5]string
}

// Format flag choices for the info command.
func Format() Formats {
	return Formats{
		Info: [5]string{"color", "json", "json.min", "text", "xml"},
	}
}

// Enabled returns all the Viper keys holding a value that are used.
// This will hide all unrecognized manual edits to the configuration file.
func Enabled() map[string]interface{} {
	sets := make(map[string]interface{})
	for _, key := range viper.AllKeys() {
		if d := get.Reset()[key]; d != nil {
			sets[key] = viper.Get(key)
		}
	}
	return sets
}

// Keys list all the available configuration setting names sorted alphabetically.
func Keys() []string {
	keys := make([]string, len(get.Reset()))
	i := 0
	for key := range get.Reset() {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// KeySort list all the available configuration setting names sorted by hand.
func KeySort() []string {
	all := Keys()
	keys := []string{get.FontFamily, get.Title, get.LayoutTmpl, get.FontEmbed,
		get.SaveDir, get.Serve, get.Editor, get.Styleh, get.Stylei}
	for _, key := range all {
		found := false
		for _, used := range keys {
			if key == used {
				found = true
				break
			}
		}
		if !found {
			keys = append(keys, key)
		}
	}
	return keys
}

// Marshal default values for use in a YAML configuration file.
func Marshal() (interface{}, error) {
	sc := Settings{}
	for key := range get.Reset() {
		if err := sc.marshals(key); err != nil {
			return sc, err
		}
	}
	return sc, nil
}

// marshals sets the default value for the key.
func (sc *Settings) marshals(key string) error { // nolint:gocyclo
	switch key {
	case get.Editor:
		sc.Editor = get.String(key)
	case get.FontEmbed:
		sc.HTML.Font.Embed = get.Bool(key)
	case get.FontFamily:
		sc.HTML.Font.Family = get.String(key)
	case get.LayoutTmpl:
		sc.HTML.Layout = get.String(key)
	case get.Author:
		sc.HTML.Meta.Author = get.String(key)
	case get.Scheme:
		sc.HTML.Meta.ColorScheme = get.String(key)
	case get.Desc:
		sc.HTML.Meta.Description = get.String(key)
	case get.Genr:
		sc.HTML.Meta.Generator = get.Bool(key)
	case get.Keywords:
		sc.HTML.Meta.Keywords = get.String(key)
	case get.Notlate:
		sc.HTML.Meta.Notranslate = get.Bool(key)
	case get.Referr:
		sc.HTML.Meta.Referrer = get.String(key)
	case get.Rtx:
		sc.HTML.Meta.RetroTxt = get.Bool(key)
	case get.Bot:
		sc.HTML.Meta.Robots = get.String(key)
	case get.Theme:
		sc.HTML.Meta.ThemeColor = get.String(key)
	case get.Title:
		sc.HTML.Title = get.String(key)
	case get.SaveDir:
		sc.SaveDirectory = get.String(key)
	case get.Serve:
		sc.ServerPort = get.UInt(key)
	case get.Stylei:
		sc.Style.Info = get.String(key)
	case get.Styleh:
		sc.Style.HTML = get.String(key)
	default:
		return fmt.Errorf("marshals %q: %w", key, logs.ErrConfigName)
	}
	return nil
}

// Missing returns the settings that are not found in the configuration file.
// This could be due to new features being added after the file was generated
// or because of manual file edits.
func Missing() (list []string) {
	if len(get.Reset()) == len(viper.AllSettings()) {
		return list
	}
	for key := range get.Reset() {
		if !viper.IsSet(key) {
			list = append(list, key)
		}
	}
	return list
}
