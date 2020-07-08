package config

import (
	"errors"
	"log"
	"os"
	"sort"
	"strings"

	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

// Defaults for configuration keys and values.
var Defaults = map[string]interface{}{
	"editor":                 "",
	"html.font.embed":        false,
	"html.font.family":       "vga",
	"html.layout":            "standard",
	"html.meta.author":       "",
	"html.meta.color-scheme": "",
	"html.meta.description":  "An example",
	"html.meta.generator":    true,
	"html.meta.keywords":     "",
	"html.meta.notranslate":  false,
	"html.meta.referrer":     "",
	"html.meta.retrotxt":     true,
	"html.meta.robots":       "index",
	"html.meta.theme-color":  "",
	"html.title":             "RetroTxt | example",
	"save-directory":         home(),
	"serve":                  uint(8080),
	"style.info":             "dracula",
	"style.html":             "lovelace",
}

// Hints provide brief help on the config file configurations.
var Hints = map[string]string{
	"editor":                 "text editor to launch when using " + str.Example("config edit"),
	"html.font.embed":        "encode and embed the font as Base64 binary-to-text within the CSS",
	"html.font.family":       "specifies the font to use with the HTML",
	"html.layout":            "HTML template for the layout of CSS, JS and fonts",
	"html.meta.author":       "defines the name of the page authors",
	"html.meta.color-scheme": "specifies one or more color schemes with which the page is compatible",
	"html.meta.description":  "a short and accurate summary of the content of the page",
	"html.meta.generator":    "include the RetroTxt version and page generation date?",
	"html.meta.keywords":     "words relevant to the page content",
	"html.meta.notranslate":  "used to declare that the page should not be translated by Google Translate",
	"html.meta.referrer":     "controls the Referer HTTP header attached to requests sent from the page",
	"html.meta.retrotxt":     "include a custom tag containing the meta information of the source textfile",
	"html.meta.robots":       "behaviour that crawlers from Google, Bing and other engines should use with the page",
	"html.meta.theme-color":  "indicates a suggested color that user agents should use to customize the display of the page",
	"html.title":             "page title that is shown in a browser title bar or tab",
	"save-directory":         "directory to store RetroTxt created HTML files",
	"serve":                  "serve HTML over an internal web server using this port",
	"style.info":             "syntax highlighter for the config info output",
	"style.html":             "syntax highlighter for html previews",
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

const (
	// PermF is posix permission bits for files
	PermF os.FileMode = 0660
	// PermD is posix permission bits for directories
	PermD os.FileMode = 0700

	cmdPath   = "retrotxt config"
	namedFile = "config.yaml"
)

var scope = gap.NewScope(gap.User, "retrotxt")

// Formats choices for flags
type Formats struct {
	Info    []string
	Shell   []string
	Version []string
}

// Format flag choices for info, shell and version commands.
var Format = Formats{
	Info:    []string{"color", "json", "json.min", "text", "xml"},
	Shell:   []string{"bash", "powershell", "zsh"},
	Version: []string{"color", "json", "json.min", "text"},
}

func (f Formats) String(field string) string {
	switch field {
	case "info":
		return strings.Join(f.Info, ", ")
	case "shell":
		return strings.Join(f.Shell, ", ")
	case "version":
		return strings.Join(f.Version, ", ")
	}
	return ""
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
		case "editor":
			sc.Editor = getString(key)
		case "html.font.embed":
			sc.HTML.Font.Embed = getBool(key)
		case "html.font.family":
			sc.HTML.Font.Family = getString(key)
		case "html.layout":
			sc.HTML.Layout = getString(key)
		case "html.meta.author":
			sc.HTML.Meta.Author = getString(key)
		case "html.meta.color-scheme":
			sc.HTML.Meta.ColorScheme = getString(key)
		case "html.meta.description":
			sc.HTML.Meta.Description = getString(key)
		case "html.meta.generator":
			sc.HTML.Meta.Generator = getBool(key)
		case "html.meta.keywords":
			sc.HTML.Meta.Keywords = getString(key)
		case "html.meta.notranslate":
			sc.HTML.Meta.Notranslate = getBool(key)
		case "html.meta.referrer":
			sc.HTML.Meta.Referrer = getString(key)
		case "html.meta.retrotxt":
			sc.HTML.Meta.RetroTxt = getBool(key)
		case "html.meta.robots":
			sc.HTML.Meta.Robots = getString(key)
		case "html.meta.theme-color":
			sc.HTML.Meta.ThemeColor = getString(key)
		case "html.title":
			sc.HTML.Title = getString(key)
		case "save-directory":
			sc.SaveDirectory = getString(key)
		case "serve":
			sc.ServerPort = getUint(key)
		case "style.info":
			sc.Style.Info = getString(key)
		case "style.html":
			sc.Style.HTML = getString(key)
		default:
			return sc, errors.New("default setting was missed, " + key)
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
		log.Fatal("config.getbool key does not exist or is not a bool value")
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
		log.Fatal("config.getuint key does not exist or is not a uint value")
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
		log.Fatal("config.getstring key does not exist or is not a string value")
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
		logs.LogCont(err)
		return ""
	}
	return dir
}
