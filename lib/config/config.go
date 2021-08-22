// Package config handles the user configations.
package config

import (
	"fmt"
	"os"
	"sort"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/viper"
)

const (
	editor     = "editor"
	fontEmbed  = "html.font.embed"
	fontFamily = "html.font.family"
	layoutTmpl = "html.layout"
	author     = "html.meta.author"
	scheme     = "html.meta.color-scheme"
	desc       = "html.meta.description"
	genr       = "html.meta.generator"
	keywords   = "html.meta.keywords"
	notlate    = "html.meta.notranslate"
	referr     = "html.meta.referrer"
	rtx        = "html.meta.retrotxt"
	bot        = "html.meta.robots"
	theme      = "html.meta.theme-color"
	title      = "html.title"
	saveDir    = "save-directory"
	serve      = "serve"
	stylei     = "style.info"
	styleh     = "style.html"

	filemode  os.FileMode = 0660
	namedFile             = "config.yaml"
)

// Defaults for setting keys.
type Defaults map[string]interface{}

// Reset configuration values.
// These are the default values whenever a setting is deleted,
// or when a new configuration file is created.
func Reset() Defaults {
	return Defaults{
		editor:     "",
		fontEmbed:  false,
		fontFamily: "vga",
		layoutTmpl: "standard",
		author:     "",
		scheme:     "",
		desc:       "",
		genr:       true,
		keywords:   "",
		notlate:    false,
		referr:     "",
		rtx:        true,
		bot:        "",
		theme:      "",
		title:      meta.Name,
		saveDir:    "",
		serve:      meta.WebPort,
		stylei:     "dracula",
		styleh:     "lovelace",
	}
}

// Hints for configuration values.
type Hints map[string]string

// Tip provides a brief help on the config file configurations.
func Tip() Hints {
	return Hints{
		editor:     "text editor to launch when using " + str.Example("config edit"),
		fontEmbed:  "encode and embed the font as Base64 binary-to-text within the CSS",
		fontFamily: "specifies the font to use with the HTML",
		layoutTmpl: "HTML template for the layout of CSS, JS and fonts",
		author:     "defines the name of the page authors",
		scheme:     "specifies one or more color schemes with which the page is compatible",
		desc:       "a short and accurate summary of the content of the page",
		genr:       fmt.Sprintf("include the %s version and page generation date?", meta.Name),
		keywords:   "words relevant to the page content",
		notlate:    "used to declare that the page should not be translated by Google Translate",
		referr:     "controls the Referer HTTP header attached to requests sent from the page",
		rtx:        "include a custom tag containing the meta information of the source textfile",
		bot:        "behavor that crawlers from Google, Bing and other engines should use with the page",
		theme:      "indicates a suggested color that user agents should use to customize the display of the page",
		title:      "page title that is shown in the browser tab",
		saveDir:    fmt.Sprintf("directory to store HTML assets created by %s", meta.Name),
		serve:      "serve files using an internal web server with this port",
		stylei:     "syntax highlighter for the config info output",
		styleh:     "syntax highlighter for html previews",
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
		if d := Reset()[key]; d != nil {
			sets[key] = viper.Get(key)
		}
	}
	return sets
}

// Keys list all the available configuration setting names sorted alphabetically.
func Keys() []string {
	keys := make([]string, len(Reset()))
	i := 0
	for key := range Reset() {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// KeySort list all the available configuration setting names sorted by hand.
func KeySort() []string {
	all := Keys()
	keys := []string{fontFamily, title, layoutTmpl, fontEmbed,
		saveDir, serve, editor, styleh, stylei}
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
	for key := range Reset() {
		if err := sc.marshals(key); err != nil {
			return sc, err
		}
	}
	return sc, nil
}

// marshals sets the default value for the key.
func (sc *Settings) marshals(key string) error { // nolint:gocyclo
	switch key {
	case editor:
		sc.Editor = getString(key)
	case fontEmbed:
		sc.HTML.Font.Embed = getBool(key)
	case fontFamily:
		sc.HTML.Font.Family = getString(key)
	case layoutTmpl:
		sc.HTML.Layout = getString(key)
	case author:
		sc.HTML.Meta.Author = getString(key)
	case scheme:
		sc.HTML.Meta.ColorScheme = getString(key)
	case desc:
		sc.HTML.Meta.Description = getString(key)
	case genr:
		sc.HTML.Meta.Generator = getBool(key)
	case keywords:
		sc.HTML.Meta.Keywords = getString(key)
	case notlate:
		sc.HTML.Meta.Notranslate = getBool(key)
	case referr:
		sc.HTML.Meta.Referrer = getString(key)
	case rtx:
		sc.HTML.Meta.RetroTxt = getBool(key)
	case bot:
		sc.HTML.Meta.Robots = getString(key)
	case theme:
		sc.HTML.Meta.ThemeColor = getString(key)
	case title:
		sc.HTML.Title = getString(key)
	case saveDir:
		sc.SaveDirectory = getString(key)
	case serve:
		sc.ServerPort = getUint(key)
	case stylei:
		sc.Style.Info = getString(key)
	case styleh:
		sc.Style.HTML = getString(key)
	default:
		return fmt.Errorf("marshals %q: %w", key, logs.ErrCfgName)
	}
	return nil
}

// getBool returns the default boolean value for the key.
func getBool(key string) bool {
	if v := viper.Get(key); v != nil {
		return v.(bool)
	}
	switch Reset()[key].(type) {
	case bool:
		return Reset()[key].(bool)
	default:
		logs.ProblemMarkFatal(key, ErrBool, logs.ErrCfgName)
	}
	return false
}

// getUint returns the default integer value for the key.
func getUint(key string) uint {
	if v := viper.GetUint(key); v != 0 {
		return v
	}
	switch Reset()[key].(type) {
	case uint:
		return Reset()[key].(uint)
	default:
		logs.ProblemMarkFatal(key, ErrUint, logs.ErrCfgName)
	}
	return 0
}

// getString returns the default string value for the key.
func getString(key string) string {
	if v := viper.GetString(key); v != "" {
		return v
	}
	switch Reset()[key].(type) {
	case string:
		return Reset()[key].(string)
	default:
		logs.ProblemMarkFatal(key, ErrString, logs.ErrCfgName)
	}
	return ""
}

// Missing returns the settings that are not found in the configuration file.
// This could be due to new features being added after the file was generated
// or because of manual file edits.
func Missing() (list []string) {
	if len(Reset()) == len(viper.AllSettings()) {
		return list
	}
	for key := range Reset() {
		if !viper.IsSet(key) {
			list = append(list, key)
		}
	}
	return list
}
