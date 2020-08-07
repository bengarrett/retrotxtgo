package config

import (
	"errors"
	"fmt"
	"os"
	"sort"

	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

const (
	missingKey      = "key does not exist or is not a bool value"
	httpPort   uint = 8080

	editor   = "editor"
	fontEmb  = "html.font.embed"
	fontFam  = "html.font.family"
	layout   = "html.layout"
	author   = "html.meta.author"
	scheme   = "html.meta.color-scheme"
	desc     = "html.meta.description"
	genr     = "html.meta.generator"
	keywords = "html.meta.keywords"
	notlate  = "html.meta.notranslate"
	referr   = "html.meta.referrer"
	rtx      = "html.meta.retrotxt"
	bot      = "html.meta.robots"
	theme    = "html.meta.theme-color"
	title    = "html.title"
	saveDir  = "save-directory"
	serve    = "serve"
	stylei   = "style.info"
	styleh   = "style.html"
)

// Defaults for configuration keys and values.
var Defaults = map[string]interface{}{
	editor:   "",
	fontEmb:  false,
	fontFam:  "vga",
	layout:   "standard",
	author:   "",
	scheme:   "",
	desc:     "An example",
	genr:     true,
	keywords: "",
	notlate:  false,
	referr:   "",
	rtx:      true,
	bot:      "index",
	theme:    "",
	title:    "RetroTxt | example",
	saveDir:  home(),
	serve:    httpPort,
	stylei:   "dracula",
	styleh:   "lovelace",
}

// Hints provide brief help on the config file configurations.
var Hints = map[string]string{
	editor:        "text editor to launch when using " + str.Example("config edit"),
	fontEmb:       "encode and embed the font as Base64 binary-to-text within the CSS",
	fontFam:       "specifies the font to use with the HTML",
	"html.layout": "HTML template for the layout of CSS, JS and fonts",
	author:        "defines the name of the page authors",
	scheme:        "specifies one or more color schemes with which the page is compatible",
	desc:          "a short and accurate summary of the content of the page",
	genr:          "include the RetroTxt version and page generation date?",
	keywords:      "words relevant to the page content",
	notlate:       "used to declare that the page should not be translated by Google Translate",
	referr:        "controls the Referer HTTP header attached to requests sent from the page",
	rtx:           "include a custom tag containing the meta information of the source textfile",
	bot:           "behaviour that crawlers from Google, Bing and other engines should use with the page",
	theme:         "indicates a suggested color that user agents should use to customize the display of the page",
	title:         "page title that is shown in a browser title bar or tab",
	saveDir:       "directory to store RetroTxt created HTML files",
	serve:         "serve HTML over an internal web server using this port",
	stylei:        "syntax highlighter for the config info output",
	styleh:        "syntax highlighter for html previews",
}

// Settings types and names to be saved in YAML.
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
			Generator   bool   `yaml:"generator"`
			Keywords    string `yaml:"keywords"`
			Notranslate bool   `yaml:"notranslate"`
			Referrer    string `yaml:"referrer"`
			RetroTxt    bool   `yaml:"retrotxt"`
			Robots      string `yaml:"robots"`
			ThemeColor  string `yaml:"theme-color"`
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

// filemode is posix permission bits for files.
const filemode os.FileMode = 0660

const (
	cmdPath   = "retrotxt config"
	namedFile = "config.yaml"
)

var scope = gap.NewScope(gap.User, "retrotxt")

// Formats choices for flags
type Formats struct {
	Info    [5]string
	Shell   [3]string
	Version [4]string
}

// Format flag choices for info, shell and version commands.
var Format = Formats{
	Info:    [5]string{"color", "json", "json.min", "text", "xml"},
	Shell:   [3]string{"bash", "powershell", "zsh"},
	Version: [4]string{"color", "json", "json.min", "text"},
}

// Enabled returns all the Viper keys holding a value that are used by RetroTxt.
// This will hide all unrecognised manual edits to the configuration file.
func Enabled() map[string]interface{} {
	var sets = make(map[string]interface{})
	for _, key := range viper.AllKeys() {
		if d := Defaults[key]; d != nil {
			sets[key] = viper.Get(key)
		}
	}
	return sets
}

// Keys list all the available configuration setting names sorted.
func Keys() []string {
	var keys = make([]string, len(Defaults))
	i := 0
	for key := range Defaults {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// Marshal default values for use in a YAML configuration file.
func Marshal() (interface{}, error) {
	var sc = Settings{}
	for key := range Defaults {
		switch key {
		case editor:
			sc.Editor = getString(key)
		case fontEmb:
			sc.HTML.Font.Embed = getBool(key)
		case fontFam:
			sc.HTML.Font.Family = getString(key)
		case layout:
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
			return sc, fmt.Errorf("unknown configuration name: %q", key)
		}
	}
	return sc, nil
}

func getBool(key string) bool {
	if v := viper.Get(key); v != nil {
		return v.(bool)
	}
	switch Defaults[key].(type) {
	case bool:
		return Defaults[key].(bool)
	default:
		logs.Fatal("getbool", key, errors.New(missingKey))
	}
	return false
}

func getUint(key string) uint {
	if v := viper.GetUint(key); v != 0 {
		return v
	}
	switch Defaults[key].(type) {
	case uint:
		return Defaults[key].(uint)
	default:
		logs.Fatal("getunit", key, errors.New(missingKey))
	}
	return 0
}

func getString(key string) string {
	if v := viper.GetString(key); v != "" {
		return v
	}
	switch Defaults[key].(type) {
	case string:
		return Defaults[key].(string)
	default:
		logs.Fatal("getstring", key, errors.New(missingKey))
	}
	return ""
}

// Missing lists the settings that are not found in the configuration file.
// This could be due to new features being added after the file was generated
// or because of manual edits.
func Missing() (list []string) {
	d, l := len(Defaults), len(viper.AllSettings())
	if d == l {
		return list
	}
	list = make([]string, l)
	for key := range Defaults {
		if !viper.IsSet(key) {
			list = append(list, key)
		}
	}
	return list
}

func home() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		logs.Log(err)
		return ""
	}
	return dir
}
