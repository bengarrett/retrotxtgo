package config

import (
	"errors"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

// Defaults for configuration values.
var Defaults = map[string]interface{}{
	"create.layout":            "standard",
	"create.meta.author":       "",
	"create.meta.color-scheme": "",
	"create.meta.description":  "An example",
	"create.meta.generator":    true,
	"create.meta.keywords":     "",
	"create.meta.notranslate":  false,
	"create.meta.referrer":     "",
	"create.meta.robots":       "index",
	"create.meta.theme-color":  "",
	"create.save-directory":    home(),
	"create.server-port":       8080,
	"create.title":             "RetroTxt | example",
	"editor":                   "",
	"style.info":               "dracula",
	"style.html":               "lovelace",
}

// Hints provide brief help on the config file configurations.
var Hints = map[string]string{
	"create.layout":            "HTML output layout",
	"create.meta.author":       "defines the name of the page authors",
	"create.meta.color-scheme": "specifies one or more color schemes with which the page is compatible",
	"create.meta.description":  "a short and accurate summary of the content of the page",
	"create.meta.generator":    "include the RetroTxt version and page generation date?",
	"create.meta.keywords":     "words relevant to the page content",
	"create.meta.notranslate":  "used to declare that the page should not be translated by Google Translate",
	"create.meta.referrer":     "controls the Referer HTTP header attached to requests sent from the page",
	"create.meta.robots":       "behaviour that crawlers from Google, Bing and other engines should use with the page",
	"create.meta.theme-color":  "indicates a suggested color that user agents should use to customize the display of the page",
	"create.save-directory":    "directory to store RetroTxt created HTML files",
	"create.server":            "serve HTML over an internal web server",
	"create.server-port":       "port which the internal web server will use",
	"create.title":             "page title that is shown in a browser title bar or tab",
	"editor":                   "text editor to launch when using " + str.Example("config edit"),
	"style.info":               "syntax highlighter for the config info output",
	"style.html":               "syntax highlighter for html previews",
}

// Settings types and names to be saved in YAML.
type Settings struct {
	Create struct {
		Layout string `yaml:"layout"`
		Meta   struct {
			Author      string `yaml:"author"`
			ColorScheme string `yaml:"color-scheme"`
			Description string `yaml:"description"`
			Generator   bool   `yaml:"generator"`
			Keywords    string `yaml:"keywords"`
			Notranslate bool   `yaml:"notranslate"`
			Referrer    string `yaml:"referrer"`
			Robots      string `yaml:"robots"`
			ThemeColor  string `yaml:"theme-color"`
		}
		SaveDirectory string `yaml:"save-directory"`
		ServerPort    int    `yaml:"server-port"`
		Title         string `yaml:"title"`
	}
	Editor string
	Style  struct {
		Info string `yaml:"info"`
		HTML string `yaml:"html"`
	}
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

// Keys list all the available configuration setting names.
func Keys() []string {
	var keys = make([]string, len(Defaults))
	i := 0
	for key := range Defaults {
		keys[i] = key
		i++
	}
	return keys
}

// Marshal default values.
func Marshal() (interface{}, error) {
	var sc = Settings{}
	for key, def := range Defaults {
		switch key {
		case "create.layout":
			sc.Create.Layout = def.(string)
		case "create.meta.author":
			sc.Create.Meta.Author = def.(string)
		case "create.meta.color-scheme":
			sc.Create.Meta.ColorScheme = def.(string)
		case "create.meta.description":
			sc.Create.Meta.Description = def.(string)
		case "create.meta.generator":
			sc.Create.Meta.Generator = def.(bool)
		case "create.meta.keywords":
			sc.Create.Meta.Keywords = def.(string)
		case "create.meta.notranslate":
			sc.Create.Meta.Notranslate = def.(bool)
		case "create.meta.referrer":
			sc.Create.Meta.Referrer = def.(string)
		case "create.meta.robots":
			sc.Create.Meta.Robots = def.(string)
		case "create.meta.theme-color":
			sc.Create.Meta.ThemeColor = def.(string)
		case "create.save-directory":
			sc.Create.SaveDirectory = def.(string)
		case "create.server-port":
			sc.Create.ServerPort = def.(int)
		case "create.title":
			sc.Create.Title = def.(string)
		case "editor":
			sc.Editor = def.(string)
		case "style.info":
			sc.Style.Info = def.(string)
		case "style.html":
			sc.Style.HTML = def.(string)
		default:
			return sc, errors.New("default setting was missed, " + key)
		}
	}
	return sc, nil
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
